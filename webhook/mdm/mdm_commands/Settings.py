import sys
import os
#sys.path.append(os.path.dirname(os.path.dirname(os.path.abspath(__file__))))
sys.path.insert(1, os.path.join(sys.path[0], '../..'))
# import ../../config.py
import config

from MDMCommand import MDMCommand

class Settings( MDMCommand ):
    def __init__(self, udid, device_name = "New AbaClocK iPad"):
        """
        Parameters
        ------
        udid : str
            udid of the device (40 characters long) and lower cased
        device_name: str, optional
            New name of the device
        """
        self.device_name = device_name
        super(Settings, self).__init__(udid)

    def _request_url(self):
        return config.MDM.server_url + "/v1/commands"

    def _request_type(self):
        return "Settings"

    def _request_data(self):
        # raw:
        # '{"udid":"' + self.udid + '","request_type":"Settings","settings":[{"item":"Bluetooth","enabled":true},{"item":"DeviceName","device_name":"' + self.device_name + '"}]}'
        return {
            "udid": self.udid,
            "request_type": self._request_type(),
            "settings": [
                {"item": "Bluetooth", "enabled": True},
                {"item": "DeviceName", "device_name": self.device_name},
            ]
        }

    def command_identifier(self):
        return "settings"

    def command_id(self):
        return 1