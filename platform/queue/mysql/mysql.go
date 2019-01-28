// Package queue implements a boldDB backed queue for MDM Commands.
package mysql

import (
	"context"
	"fmt"
	"strings"
	"time"
	"database/sql"
	"github.com/jmoiron/sqlx"
	sq "gopkg.in/Masterminds/squirrel.v1"
	
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/groob/plist"
	"github.com/pkg/errors"
	_ "github.com/go-sql-driver/mysql"

	"github.com/micromdm/micromdm/mdm"
	"github.com/micromdm/micromdm/platform/command"
	"github.com/micromdm/micromdm/platform/queue"
	"github.com/micromdm/micromdm/platform/pubsub"
)

const (
	DeviceCommandBucket = "mdm.DeviceCommands"
	CommandQueuedTopic = "mdm.CommandQueued"
	
	DeviceCommandTable = "device_commands"
	CommandQueueTable = "command_queue"
)

type Store struct {
	db *sqlx.DB
	logger log.Logger
}

type MysqlCommand struct {
	UUID    string			`db:"uuid"`
	DeviceUDID string		`db:"device_udid"`
	Payload []byte			`db:"payload"`

	CreatedAt    time.Time 	`db:"created_at"`
	LastSentAt   time.Time 	`db:"last_sent_at"`
	Acknowledged time.Time 	`db:"acknowledged_at"`

	TimesSent int 			`db:"times_sent"`

	LastStatus     string 	`db:"last_status"`
	FailureMessage []byte 	`db:"failure_message"`
	Order			int		`db:"command_order"`
}

func command_columns() []string {
	return []string{
		"uuid",
		"device_udid",
		"payload",
		"created_at",
		"last_sent_at",
		"acknowledged_at",
		"times_sent",
		"last_status",
		"failure_message",
		"command_order",
	}
}

type Option func(*Store)

func WithLogger(logger log.Logger) Option {
	return func(s *Store) {
		s.logger = logger
	}
}

func (db *Store) Next(ctx context.Context, resp mdm.Response) ([]byte, error) {
	cmd, err := db.nextCommand(ctx, resp)
	if err != nil {
		return nil, err
	}
	if cmd == nil {
		return nil, nil
	}
	return cmd.Payload, nil
}

func (db *Store) nextCommand(ctx context.Context, resp mdm.Response) (*queue.Command, error) {
	udid := resp.UDID
	if resp.UserID != nil {
		// use the user id for user level commands
		udid = *resp.UserID
	}
	
	update_error := db.UpdateCommandStatus(ctx, resp)
	if update_error != nil {
		return nil, update_error
	}
	
	dc, err := db.DeviceCommand(ctx, udid)
	if err != nil {
		if isNotFound(err) {
			return nil, nil
		}
		return nil, errors.Wrapf(err, "get device command from queue, udid: %s", resp.UDID)
	}
	
	var cmd *queue.Command
	switch resp.Status {
	case "NotNow":
		// We will try this command later when the device is not
		// responding with NotNow
		x, a := cut(dc.Commands, resp.CommandUUID)
		dc.Commands = a
		if x == nil {
			break
		}
		dc.NotNow = append(dc.NotNow, *x)

	case "Acknowledged":
		// move to completed, send next
		x, a := cut(dc.Commands, resp.CommandUUID)
		dc.Commands = a
		if x == nil {
			break
		}
		dc.Completed = append(dc.Completed, *x)
	case "Error":
		// move to failed, send next
		x, a := cut(dc.Commands, resp.CommandUUID)
	
		dc.Commands = a
		if x == nil { // must've already bin ackd
			break
		}
		dc.Failed = append(dc.Failed, *x)

	case "CommandFormatError":
		// move to failed
		x, a := cut(dc.Commands, resp.CommandUUID)
		dc.Commands = a
		if x == nil {
			break
		}
		dc.Failed = append(dc.Failed, *x)

	case "Idle":
		// will send next command below

	default:
		return nil, fmt.Errorf("unknown response status: %s", resp.Status)
	}
	
	
	// pop the first command from the queue and add it to the end.
	// If the regular queue is empty, send a command that got
	// refused with NotNow before.
	cmd, dc.Commands = popFirst(dc.Commands)
	if cmd != nil {
		dc.Commands = append(dc.Commands, *cmd)
	} else if resp.Status != "NotNow" {
		cmd, dc.NotNow = popFirst(dc.NotNow)
		if cmd != nil {
			dc.Commands = append(dc.Commands, *cmd)
		}
	}

	if err := db.Save(ctx, dc); err != nil {
		return nil, err
	}
	return cmd, nil
}

func popFirst(all []queue.Command) (*queue.Command, []queue.Command) {
	if len(all) == 0 {
		return nil, all
	}
	first := all[0]
	all = append(all[:0], all[1:]...)
	return &first, all
}

func cut(all []queue.Command, uuid string) (*queue.Command, []queue.Command) {
	for i, cmd := range all {
		if cmd.UUID == uuid {
			all = append(all[:i], all[i+1:]...)
			return &cmd, all
		}
	}
	return nil, all
}

func SetupDB(db *sqlx.DB) error {
	// Required for TIMESTAMP DEFAULT 0
	_,err := db.Exec(`SET sql_mode = '';`)

	// "github.com/micromdm/micromdm/platform/queue/internal/devicecommandproto"
	_,err = db.Exec(`CREATE TABLE IF NOT EXISTS `+DeviceCommandTable+` (
	    uuid VARCHAR(40) PRIMARY KEY,
	    device_udid VARCHAR(40) NOT NULL,
	    payload BLOB DEFAULT NULL,
	    created_at TIMESTAMP DEFAULT 0,
	    last_sent_at TIMESTAMP DEFAULT 0,
	    acknowledged_at TIMESTAMP DEFAULT 0,
	    times_sent int(11) DEFAULT 0,
	    last_status VARCHAR(32) DEFAULT NULL,
	    failure_message BLOB DEFAULT NULL,
	    command_order int(11) DEFAULT 0
	);`)

	if err != nil {
	   return errors.Wrap(err, "creating "+DeviceCommandTable+" sql table failed")
	}
	
	_,err = db.Exec(`ALTER TABLE `+DeviceCommandTable+` MODIFY payload MEDIUMBLOB DEFAULT NULL;`)
	if err != nil {
	   return errors.Wrap(err, "altering "+DeviceCommandTable+" sql table failed")
	}

	return nil
}

func NewQueue(db *sqlx.DB, pubsub pubsub.PublishSubscriber, opts ...Option) (*Store, error) {
	SetupDB(db)
	
	datastore := &Store{db: db, logger: log.NewNopLogger()}
	for _, fn := range opts {
		fn(datastore)
	}

	ctx := context.Background()
	if err := datastore.pollCommands(ctx, pubsub); err != nil {
		return nil, err
	}

	return datastore, nil
}

func (db *Store) SaveCommand(ctx context.Context, cmd queue.Command, deviceUDID string, order int) error {
	// Make sure we take the time offset into account for "zero" dates	
	t := time.Now()
	_, offset := t.Zone()

	// Don't multiply by zero
	if (offset <= 0) {
		offset = 1
	}
	var min_timestamp_sec int64 = int64(offset) * 60 * 60 * 24
	
	if (cmd.CreatedAt.IsZero() || cmd.CreatedAt.Unix() < min_timestamp_sec) {
		cmd.CreatedAt = time.Unix(min_timestamp_sec, 0)
	}
	
	if (cmd.LastSentAt.IsZero() || cmd.LastSentAt.Unix() < min_timestamp_sec) {
		cmd.LastSentAt = time.Unix(min_timestamp_sec, 0)
	}
	
	if (cmd.Acknowledged.IsZero() || cmd.Acknowledged.Unix() < min_timestamp_sec) {
		cmd.Acknowledged = time.Unix(min_timestamp_sec, 0)
	}
	
	updateQuery, args_update, err := sq.StatementBuilder.
		PlaceholderFormat(sq.Question).
		Update(DeviceCommandTable).
		Prefix("ON DUPLICATE KEY").
		Set("uuid", cmd.UUID).
		Set("device_udid", deviceUDID).
		Set("payload", cmd.Payload).
		Set("created_at", cmd.CreatedAt).
		Set("last_sent_at", cmd.LastSentAt).
		Set("acknowledged_at", cmd.Acknowledged).
		Set("times_sent", cmd.TimesSent).
		Set("last_status", cmd.LastStatus).
		Set("failure_message", cmd.FailureMessage).
		Set("command_order", order).
		ToSql()
	if err != nil {
		return errors.Wrap(err, "building update query for command save")
	}
	
	// MySql Convention
	// Replace "ON DUPLICATE KEY UPDATE TABLE_NAME SET" to "ON DUPLICATE KEY UPDATE"
	updateQuery = strings.Replace(updateQuery, DeviceCommandTable+" SET ", "", -1)

	query, args, err := sq.StatementBuilder.
		PlaceholderFormat(sq.Question).
		Insert(DeviceCommandTable).
		Columns(command_columns()...).
		Values(
			cmd.UUID,
			deviceUDID,
			cmd.Payload,
			cmd.CreatedAt,
			cmd.LastSentAt,
			cmd.Acknowledged,
			cmd.TimesSent,
			cmd.LastStatus,
			cmd.FailureMessage,
			order,
		).
		Suffix(updateQuery).
		ToSql()
	
	var all_args = append(args, args_update...)
	
	if err != nil {
		return errors.Wrap(err, "building command save query")
	}
	
	_, err = db.db.ExecContext(ctx, query, all_args...)
	return errors.Wrap(err, "exec command save in mysql")
}

func (db *Store) Save(ctx context.Context, cmd *queue.DeviceCommand) error {
	SetupDB(db.db)

	var err error
	
	for i, _command := range cmd.Commands {
		err = db.SaveCommand(ctx, _command, cmd.DeviceUDID, i)
		if err != nil {
			return err
		}
	}
	return err
}

func (db *Store) DeviceCommand(ctx context.Context, udid string) (*queue.DeviceCommand, error) {
	query, args, err := sq.StatementBuilder.
		PlaceholderFormat(sq.Question).
		Select(command_columns()...).
		From(DeviceCommandTable).
		Where(sq.Eq{"device_udid": udid}).
		OrderBy("command_order").
		ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "building sql")
	}

	var list []MysqlCommand
	err = db.db.SelectContext(ctx, &list, query, args...)
	if errors.Cause(err) == sql.ErrNoRows {
		return nil, deviceCommandNotFoundErr{}
	}
	dev, _err := UnmarshalMysqlCommand(udid, list)
	if _err != nil {
		return nil, _err
	}
	return &dev, errors.Wrap(err, "finding device_commands by udid")
}

type notFound struct {
	ResourceType string
	Message      string
}

func (e *notFound) Error() string {
	return fmt.Sprintf("not found: %s %s", e.ResourceType, e.Message)
}

func (db *Store) pollCommands(ctx context.Context, pubsub pubsub.PublishSubscriber) error {
	commandEvents, err := pubsub.Subscribe(context.TODO(), "command-queue", command.CommandTopic)
	if err != nil {
		return errors.Wrapf(err,
			"subscribing push to %s topic", command.CommandTopic)
	}
	go func() {
		for {
			select {
			case event := <-commandEvents:
				var ev command.Event
				if err := command.UnmarshalEvent(event.Message, &ev); err != nil {
					level.Info(db.logger).Log("msg", "unmarshal command event in queue", "err", err)
					continue
				}

				cmd := new(queue.DeviceCommand)
				cmd.DeviceUDID = ev.DeviceUDID
				byUDID, err := db.DeviceCommand(ctx, ev.DeviceUDID)
				if err == nil && byUDID != nil {
					cmd = byUDID
				}
				newPayload, err := plist.Marshal(ev.Payload)
				if err != nil {
					level.Info(db.logger).Log("msg", "marshal event payload", "err", err)
					continue
				}
				newCmd := queue.Command{
					UUID:    ev.Payload.CommandUUID,
					Payload: newPayload,
				}
				cmd.Commands = append(cmd.Commands, newCmd)
				if err := db.Save(ctx, cmd); err != nil {
					level.Info(db.logger).Log("msg", "save command in db", "err", err)
					continue
				}
				level.Info(db.logger).Log(
					"msg", "queued event for device",
					"device_udid", ev.DeviceUDID,
					"command_uuid", ev.Payload.CommandUUID,
					"request_type", ev.Payload.Command.RequestType,
				)

				cq := new(queue.QueueCommandQueued)
				cq.DeviceUDID = ev.DeviceUDID
				cq.CommandUUID = ev.Payload.CommandUUID

				msgBytes, err := queue.MarshalQueuedCommand(cq)
				if err != nil {
					level.Info(db.logger).Log("msg", "marshal queued command", "err", err)
					continue
				}

				if err := pubsub.Publish(context.TODO(), CommandQueuedTopic, msgBytes); err != nil {
					level.Info(db.logger).Log("msg", "publish command to queued topic", "err", err)
				}
			}
		}
	}()
	return nil
}

func isNotFound(err error) bool {
	if _, ok := err.(*notFound); ok {
		return true
	}
	return false
}


type deviceCommandNotFoundErr struct{}

func (e deviceCommandNotFoundErr) Error() string {
	return "device not found"
}

func (e deviceCommandNotFoundErr) NotFound() bool {
	return true
}


/*
func (deviceCommand *queue.DeviceCommand) marshalMysqlCommand() []MysqlCommand {
	// queue.DeviceCommand --> List of MysqlCommand
    return nil
}
*/

func (db *Store) UpdateCommandStatus(ctx context.Context, resp mdm.Response) error {
	query, args, err := sq.StatementBuilder.
		PlaceholderFormat(sq.Question).
		Update(DeviceCommandTable).
		Set("last_status", resp.Status).
		Where(sq.Eq{"uuid": resp.CommandUUID}).
		ToSql()
	_, err = db.db.ExecContext(ctx, query, args...)
	if err != nil {
		return errors.Wrap(err, "building update query for command save")
	}
	
	return errors.Wrap(err, "exec command update in mysql")
}

func UnmarshalMysqlCommand(udid string, mysqlCommands []MysqlCommand) (queue.DeviceCommand, error) {
	
	var dev = queue.DeviceCommand {
		DeviceUDID: udid,
	}
	
	for _, command := range mysqlCommands {
		if command.DeviceUDID == udid {
			var cmd = queue.Command {
				UUID:         	command.UUID,
				Payload:      	command.Payload,
				CreatedAt:    	command.CreatedAt,
				LastSentAt:   	command.LastSentAt,
				Acknowledged: 	command.Acknowledged,
	
				TimesSent: 		command.TimesSent,
	
				LastStatus:     command.LastStatus,
				FailureMessage: command.FailureMessage,
			}
			
			switch cmd.LastStatus {
			case "NotNow":
				dev.NotNow = append(dev.NotNow, cmd)
		
			case "Acknowledged":
				dev.Completed = append(dev.Completed, cmd)
				
			case "Error":
				dev.Failed = append(dev.Failed, cmd)
						
			case "CommandFormatError":
				dev.Failed = append(dev.Failed, cmd)
		
			case "Idle":
				// will send next command below
				dev.Commands = append(dev.Commands, cmd)
		
			default:
				// Not yet classified
				dev.Commands = append(dev.Commands, cmd)
			}
		}
	}
	return dev, nil
}