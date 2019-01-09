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
                       "(command_uuid, command_id, command_identifier, udid, created_at, updated_at) "
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
            #print(cursor.statement)

            # Make sure data is committed to the database
            cnx.commit()
            cnx.close()

        def _log_device(self, device):
            cnx = mysql.connector.connect(**config.Mysql.config)
            cursor = cnx.cursor()

            now = datetime.now()

            add_log = ("INSERT INTO psm_mdm_devices "
                       "(build_version, os_version, product_name, serial_number, udid) "
                       "VALUES (%s, %s, %s, %s, %s) "
                       "ON DUPLICATE KEY UPDATE "
                       "build_version=VALUES(build_version), "
                       "os_version=VALUES(os_version), "
                       "product_name=VALUES(product_name), "
                       "serial_number=VALUES(serial_number), "
                       "udid=VALUES(udid)")
            data_log = (
                device.build_version,
                device.os_version,
                device.product_name,
                device.serial_number,
                device.udid,
            )

            # Insert new command
            cursor.execute(add_log, data_log)
            print(cursor.statement)

            # Make sure data is committed to the database
            cnx.commit()
            cnx.close()

        def _log_profile_payload(self, profile_payload):
            cnx = mysql.connector.connect(**config.Mysql.config)
            cursor = cnx.cursor()

            now = datetime.now()

            add_log = ("INSERT INTO psm_mdm_profiles_payload "
                       "(device_udid, profile_uuid, payload_version, payload_description, payload_type,"
                       "payload_identifier, payload_display_name) "
                       "VALUES (%s, %s, %s, %s, %s, %s, %s) "
                       "ON DUPLICATE KEY UPDATE "
                       "payload_version=VALUES(payload_version), "
                       "payload_description=VALUES(payload_description), "
                       "payload_type=VALUES(payload_type), "
                       "payload_identifier=VALUES(payload_identifier), "
                       "payload_display_name=VALUES(payload_display_name)")
            data_log = (
                profile_payload.udid,
                profile_payload.profile_uuid,
                profile_payload.payload_version,
                profile_payload.payload_description,
                profile_payload.payload_type,
                profile_payload.payload_identifier,
                profile_payload.payload_display_name,
            )

            # Insert new command
            cursor.execute(add_log, data_log)
            print(cursor.statement)

            # Make sure data is committed to the database
            cnx.commit()
            cnx.close()

        def _log_profile(self, profile):
            cnx = mysql.connector.connect(**config.Mysql.config)
            cursor = cnx.cursor()

            now = datetime.now()

            add_log = ("INSERT INTO psm_mdm_profiles "
                       "(udid, payload_uuid, payload_description, has_removal_passcode, payload_identifier,"
                       "payload_display_name, is_managed, is_encrypted, payload_version, payload_organization,"
                       "payload_removal_disallowed) "
                       "VALUES (%s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s) "
                       "ON DUPLICATE KEY UPDATE "
                       "udid=VALUES(udid), "
                       "payload_uuid=VALUES(payload_uuid), "
                       "payload_description=VALUES(payload_description), "
                       "has_removal_passcode=VALUES(has_removal_passcode), "
                       "payload_identifier=VALUES(payload_identifier), "
                       "payload_display_name=VALUES(payload_display_name), "
                       "is_managed=VALUES(is_managed), "
                       "is_encrypted=VALUES(is_encrypted), "
                       "payload_version=VALUES(payload_version), "
                       "payload_organization=VALUES(payload_organization), "
                       "payload_removal_disallowed=VALUES(payload_removal_disallowed)")
            data_log = (
                profile.udid,
                profile.payload_uuid,
                profile.payload_description,
                profile.has_removal_passcode,
                profile.payload_identifier,
                profile.payload_display_name,
                profile.is_managed,
                profile.is_encrypted,
                profile.payload_version,
                profile.payload_organization,
                profile.payload_removal_disallowed,
            )

            # Insert new command
            cursor.execute(add_log, data_log)
            print(cursor.statement)

            # Make sure data is committed to the database
            cnx.commit()
            cnx.close()

            for payload in profile.profile_payloads:
                self._log_profile_payload(payload)

        def _log_os_update_status(self, os_update_status):
            cnx = mysql.connector.connect(**config.Mysql.config)
            cursor = cnx.cursor()

            now = datetime.now()

            add_log = ("INSERT INTO psm_mdm_os_update_status "
                       "(download_percent_complete, is_downloaded, product_key, status, udid) "
                       "VALUES (cast(%s as decimal(6,5)), %s, %s, %s, %s) "
                       "ON DUPLICATE KEY UPDATE "
                       "download_percent_complete=VALUES(download_percent_complete), "
                       "is_downloaded=VALUES(is_downloaded), "
                       "product_key=VALUES(product_key), "
                       "status=VALUES(status)")

            data_log = (
                os_update_status.download_percent_complete,
                os_update_status.is_downloaded,
                os_update_status.product_key,
                os_update_status.status,
                os_update_status.udid,
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