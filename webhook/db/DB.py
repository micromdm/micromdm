import sys
import os

sys.path.insert(1, os.path.join(sys.path[0], '..'))
# import ../config.py
import mdm

class DB (object):
    class __DB:
        def __init__(self):
            print("init")
            self.val = None

        def __str__(self):
            return `self` + self.val

    instance = None
    def __new__(cls): # __new__ always a classmethod
        if not DB.instance:
            import MysqlDB
            DB.instance = MysqlDB.MysqlDB()

        return DB.instance


    def _log_command_request(self, mdm_command):
        """
        Parameters
        ------
        mdm_command: mdm.MDMCommand
            Command that is going to be requested to the mdm
        Raises
        ------
        NotImplementedError
            If implementation class does not implement this method
        """
        raise NotImplementedError("Should have implemented this: _request_url")

    @staticmethod
    def log_command_request(mdm_command):
        """
        Parameters
        ------
        mdm_command: mdm.MDMCommand
           Command that is going to be requested to the mdm
        """
        DB()
        DB.instance._log_command_request(mdm_command)

    def _log_command_response(self, udid, command_uuid, status):
        """
        Parameters
        ------
        udid: str
            Device UDID. Length 40 Characters
        command_uuid: str
            UUID of related command. Length 40 Characters
        status: str
            Response of the command
            see: https://developer.apple.com/business/documentation/MDM-Protocol-Reference.pdf
                Acknowledged
                Error
                CommandFormatError
                Idle
                NotNow
        Raises
        ------
        NotImplementedError
            If implementation class does not implement this method
        """
        raise NotImplementedError("Should have implemented this: _request_url")

    @staticmethod
    def log_command_response(udid, command_uuid, status):
        """
        Parameters
        ------
        udid: str
            Device UDID. Length 40 Characters
        command_uuid: str
            UUID of related command. Length 40 Characters
        status: str
            Response of the command
            see: https://developer.apple.com/business/documentation/MDM-Protocol-Reference.pdf
                Acknowledged
                Error
                CommandFormatError
                Idle
                NotNow
        """
        DB()
        DB.instance._log_command_response(udid, command_uuid, status)