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
        raise NotImplementedError("Should have implemented this: _log_command_request")

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
        raise NotImplementedError("Should have implemented this: _log_command_response")

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


    def _log_device(self, device):
        """
        Parameters
        ------
        device: mdm.MDMDevice
            MDMDevice of package mdm, holding all relevant information about
                Device udid
                Serial Number
                Build Version
                OS Version
                Product Name
        Raises
        ------
        NotImplementedError
            If implementation class does not implement this method
        """
        raise NotImplementedError("Should have implemented this: _log_device")

    @staticmethod
    def log_device(device):
        """
        Parameters
        ------
        device: mdm.MDMDevice
            MDMDevice of package mdm, holding all relevant information about
                Device udid
                Serial Number
                Build Version
                OS Version
                Product Name
        """
        DB()
        DB.instance._log_device(device)


    def _log_profile_payload(self, profile_payload):
        """
        Parameters
        ------
        device: mdm.MDMProfilePayload
            MDMProfilePayload of package mdm, profile content of a related profile
        Raises
        ------
        NotImplementedError
            If implementation class does not implement this method
        """
        raise NotImplementedError("Should have implemented this: _log_profile_payload")

    def _log_profile(self, profile):
        """
        Parameters
        ------
        device: mdm.MDMProfile
            MDMProfile of package mdm, holding all relevant information about a Profile and its payload objects
        Raises
        ------
        NotImplementedError
            If implementation class does not implement this method
        """
        raise NotImplementedError("Should have implemented this: _log_profile")

    @staticmethod
    def log_profile(profile):
        """
        Parameters
        ------
        device: mdm.MDMProfile
            MDMProfile of package mdm, holding all relevant information about a Profile and its payload objects
        """
        DB()
        DB.instance._log_profile(profile)





    def _log_os_update_status(self, os_update_status):
        """
        Parameters
        ------
        os_update_status: mdm.MDMOSUpdateStatus
            MDMOSUpdateStatus of package mdm, holding all relevant information about
                Device udid
                DownloadPercentComplete
                IsDownloaded
                ProductKey
                Status: see https://developer.apple.com/business/documentation/MDM-Protocol-Reference.pdf
                    Idle: No action is being taken on this software update.
                    Downloading: The software update is being downloaded.
                    DownloadFailed: The download has failed.
                    DownloadRequiresComputer: The device must be connected to a computer to download this update (iOS only).
                    DownloadInsufficientSpace: There is not enough space to download the update.
                    DownloadInsufficientPower: There is not enough power to download the update.
                    DownloadInsufficientNetwork: There is insufficient network capacity to download the update.
                    Installing: The software update is being installed.
                    InstallInsufficientSpace: There is not enough space to install the update.
                    InstallInsufficientPower: There is not enough power to install the update.
                    InstallPhoneCallInProgress: Installation has been rejected because a phone call is in progress.
                    InstallFailed: Installation has failed for an unspecified reason.
        Raises
        ------
        NotImplementedError
            If implementation class does not implement this method
        """
        raise NotImplementedError("Should have implemented this: _log_os_update_status")

    @staticmethod
    def log_os_update_status(os_update_status):
        """
        Parameters
        ------
        os_update_status: mdm.MDMOSUpdateStatus
            MDMOSUpdateStatus of package mdm, holding all relevant information about
                Device udid
                DownloadPercentComplete
                IsDownloaded
                ProductKey
                Status: see https://developer.apple.com/business/documentation/MDM-Protocol-Reference.pdf
                    Idle: No action is being taken on this software update.
                    Downloading: The software update is being downloaded.
                    DownloadFailed: The download has failed.
                    DownloadRequiresComputer: The device must be connected to a computer to download this update (iOS only).
                    DownloadInsufficientSpace: There is not enough space to download the update.
                    DownloadInsufficientPower: There is not enough power to download the update.
                    DownloadInsufficientNetwork: There is insufficient network capacity to download the update.
                    Installing: The software update is being installed.
                    InstallInsufficientSpace: There is not enough space to install the update.
                    InstallInsufficientPower: There is not enough power to install the update.
                    InstallPhoneCallInProgress: Installation has been rejected because a phone call is in progress.
                    InstallFailed: Installation has failed for an unspecified reason.
        Raises
        ------
        NotImplementedError
            If implementation class does not implement this method
        """
        DB()
        DB.instance._log_os_update_status(os_update_status)