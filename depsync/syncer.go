package depsync

import (
	"errors"
	"fmt"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/micromdm/dep"
)

// Syncer is responsible for fetching devices from DEP, keeping them in sync, and storing cursor information about the
// current fetch operation.
type Syncer interface {
	Start(deviceChan chan<- dep.Device)
	Fetch(deviceChan chan<- dep.Device) (bool, error)
}

// Cursor represents information about the current fetch operation. It is required to fetch multiple "pages" of devices.
type Cursor struct {
	Value   string
	Created time.Time // >7 days cursor is invalid.
}

type syncer struct {
	logger               log.Logger
	client               dep.Client
	fetchInterval        time.Duration
	ticker               *time.Ticker
	done                 <-chan struct{}
	Cursor               *Cursor
	InitialFetchComplete bool
	errorCount           int
}

func NewSyncer(client dep.Client, logger log.Logger, fetchInterval time.Duration, done <-chan struct{}) Syncer {
	return &syncer{
		logger:        logger,
		client:        client,
		fetchInterval: fetchInterval,
		errorCount:    0,
		done:          done,
	}
}

// Fetch fetches a list of DEP devices in batches. The first return value indicates that there are more devices to fetch.
func (s *syncer) Fetch(device chan<- dep.Device) (moreResults bool, err error) {
	var deviceResponse *dep.DeviceResponse

	if s.Cursor != nil {
		deviceResponse, err = s.client.FetchDevices(dep.Limit(100), dep.Cursor(s.Cursor.Value))
	} else {
		deviceResponse, err = s.client.FetchDevices(dep.Limit(100))
	}

	if err != nil {
		return true, err
	}

	if s.Cursor == nil {
		created := time.Now()
		s.Cursor = &Cursor{
			Value:   deviceResponse.Cursor,
			Created: created,
		}
	}

	s.logger.Log("level", "debug", "msg", fmt.Sprintf("Fetching %d devices", len(deviceResponse.Devices)))
	for _, dev := range deviceResponse.Devices {
		device <- dev
	}

	return deviceResponse.MoreToFollow, nil
}

// Sync fetches devices since the given cursor value
func (s *syncer) Sync(device chan<- dep.Device) (bool, error) {
	var deviceResponse *dep.DeviceResponse
	var err error

	if s.Cursor != nil {
		deviceResponse, err = s.client.SyncDevices(s.Cursor.Value, dep.Limit(100))
	} else {
		return true, errors.New("cannot sync dep devices without a cursor value")
	}

	if err != nil {
		return true, err
	}

	if s.Cursor == nil {
		created := time.Now()
		s.Cursor = &Cursor{
			Value:   deviceResponse.Cursor,
			Created: created,
		}
	}

	s.logger.Log("level", "debug", "msg", fmt.Sprintf("Syncing %d devices", len(deviceResponse.Devices)))
	for _, dev := range deviceResponse.Devices {
		device <- dev
	}

	return deviceResponse.MoreToFollow, nil
}

// Start starts the synchronisation schedule which runs at a configured interval.
// The first run will always download a complete list of devices. Deltas will be fetched if the process is still running
// after the first sync interval
func (s *syncer) Start(device chan<- dep.Device) {
	s.logger.Log("level", "debug", "msg", "DEP sync routine started")
	s.ticker = time.NewTicker(s.fetchInterval)

	select {
	case <-s.ticker.C:
		if !s.InitialFetchComplete {
			more, err := s.Fetch(device)
			if err != nil {
				s.logger.Log("level", "error", "msg", fmt.Sprintf("Fetching initial snapshot of devices from DEP: %v", err))
				s.errorCount++
			} else {
				s.logger.Log("level", "debug", "msg", fmt.Sprintf("More devices after this batch: %t", more))
				if !more {
					s.InitialFetchComplete = true
				}
			}
		} else {
			s.logger.Log("level", "debug", "msg", "Synchronizing devices from DEP service")

			more, err := s.Sync(device)
			if err != nil {
				s.logger.Log("level", "warn", "msg", fmt.Sprintf("Unable to fetch devices: %s", err))
			} else {
				s.logger.Log("level", "debug", "msg", fmt.Sprintf("More devices after this sync: %t", more))
			}
		}
	case <-s.done:
		s.logger.Log("level", "info", "msg", "stopping DEP sync routine")
		s.ticker.Stop()
		return
	}
}
