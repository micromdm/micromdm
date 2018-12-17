import sys
import os
#sys.path.append(os.path.dirname(os.path.dirname(os.path.abspath(__file__))))
sys.path.insert(1, os.path.join(sys.path[0], '../..'))
# import ../../config.py
import config

from MDMCommand import MDMCommand

class ProfileList( MDMCommand ):
    def __init__(self, udid):
        super(ProfileList, self).__init__(udid)

    def _request_url(self):
        return config.MDM.server_url + "/v1/commands"

    def _request_type(self):
        return "ProfileList"

    def _request_data(self):
        # raw:
        # '{"udid":"' + self.udid + '","request_type":"Settings","settings":[{"item":"Bluetooth","enabled":true},{"item":"DeviceName","device_name":"' + self.device_name + '"}]}'
        return {
            "udid": self.udid,
            "request_type": self._request_type(),
        }

    def _command_identifier(self):
        return "profile_list"

    def _command_id(self):
        return 8