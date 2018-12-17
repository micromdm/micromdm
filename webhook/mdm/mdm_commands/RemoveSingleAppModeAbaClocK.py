import sys
import os
#sys.path.append(os.path.dirname(os.path.dirname(os.path.abspath(__file__))))
sys.path.insert(1, os.path.join(sys.path[0], '../..'))
# import ../../config.py
import config

from MDMCommand import MDMCommand

class RemoveSingleAppModeAbaClocK( MDMCommand ):
    def __init__(self, udid, profile_identifier = 'ch.abacus.abaclock.client2.EF21A748-D41A-444F-9F71-F5B5CDFEC19C'):
        """
        Parameters
        ------
        udid : str
            udid of the device (40 characters long) and lower cased
        profile_identifier: str, optional
            Identifier of the Profile to be removed
        """
        self.profile_identifier = profile_identifier
        super(RemoveSingleAppModeAbaClocK, self).__init__(udid)

    def _request_url(self):
        return config.MDM.server_url + "/v1/commands"

    def _request_type(self):
        return "RemoveProfile"

    def _request_data(self):
        return {
            "udid": self.udid,
            "request_type": self._request_type(),
            "identifier": self.profile_identifier,
        }

    def _command_identifier(self):
        return "remove_profile_single_app_mode_abaclock"

    def _command_id(self):
        return 10