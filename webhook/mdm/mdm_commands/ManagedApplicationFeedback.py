import sys
import os
#sys.path.append(os.path.dirname(os.path.dirname(os.path.abspath(__file__))))
sys.path.insert(1, os.path.join(sys.path[0], '../..'))
# import ../../config.py
import config

from MDMCommand import MDMCommand

class ManagedApplicationFeedback( MDMCommand ):
    def __init__(self, udid, app_identifier = "ch.abacus.abaclock.client2"):
        """
        Parameters
        ------
        udid : str
            udid of the device (40 characters long) and lower cased
        app_identifier: str, optional
            Bundle ID of the app
        """
        self.app_identifier = app_identifier
        super(ManagedApplicationFeedback, self).__init__(udid)

    def _request_url(self):
        return config.MDM.server_url + "/v1/commands"

    def _request_type(self):
        return "ManagedApplicationFeedback"

    def _request_data(self):
        # raw:
        # '{"udid":"' + self.udid + '","request_type":"Settings","settings":[{"item":"Bluetooth","enabled":true},{"item":"DeviceName","device_name":"' + self.device_name + '"}]}'
        return {
            "udid": self.udid,
            "request_type": self._request_type(),
            "identifiers": [
                self.app_identifier
            ]
        }

    def command_identifier(self):
        return "managed_application_feedback_abaclock2"

    def command_id(self):
        return 9