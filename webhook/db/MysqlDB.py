import sys
import os
import mysql.connector

sys.path.insert(1, os.path.join(sys.path[0], '..'))
import config
from datetime import date, datetime

from DB import DB
class MysqlDB ( DB ):

    class __MysqlDB:
        def __init__(self):
            pass

        def _log_command_request(self, mdm_command):
            if mdm_command == None or mdm_command.command_uuid == None:
                print("Raise Error? Command and its uuid should not be none!")
                return

            cnx = mysql.connector.connect(**config.Mysql.config)
            cursor = cnx.cursor()

            now = datetime.now()

            add_log = ("INSERT INTO psm_mdm_commands "
                       "(command_uuid, command_id, command_identifier, device_udid, created_at, updated_at) "
                       "VALUES (%s, %s, %s, %s, %s, %s)")
            data_log = (
                mdm_command.command_uuid,
                mdm_command.command_id(),
                mdm_command.command_identifier(),
                mdm_command.udid,
                now,
                now,
            )

            # Insert new command
            cursor.execute(add_log, data_log)

            # Make sure data is committed to the database
            cnx.commit()
            cnx.close()

        def _log_command_response(self, udid, command_uuid, status):
            cnx = mysql.connector.connect(**config.Mysql.config)
            cursor = cnx.cursor()

            now = datetime.now()

            add_log = ("UPDATE psm_mdm_commands "
                       "SET status = %s "
                       "WHERE command_uuid = %s;")
            data_log = (
                status,
                command_uuid,
            )

            # Insert new command
            cursor.execute(add_log, data_log)

            print(cursor.statement)

            # Make sure data is committed to the database
            cnx.commit()
            cnx.close()







    instance = None
    def __new__(cls): # __new__ always a classmethod
        if not MysqlDB.instance:
            MysqlDB.instance = MysqlDB.__MysqlDB()
        return MysqlDB.instance