import sys
import os
#sys.path.append(os.path.dirname(os.path.dirname(os.path.abspath(__file__))))
sys.path.insert(1, os.path.join(sys.path[0], '../..'))
# import ../../config.py
import config

from MDMCommand import MDMCommand

class InstallVPPAbaClocKUnmanaged( MDMCommand ):
    def __init__(self, udid, adam_id_str = 1301278477):
        """
        Parameters
        ------
        udid : str
            udid of the device (40 characters long) and lower cased
        adam_id_str: int64
            id of the application, adam_id_str can be retrieved by the AppStore URL
            https://itunes.apple.com/ch/app/abaclock-2/id1301278477?mt=8
            For AbaClocK the adam_id_str would be "1301278477"
        """
        self.adam_id_str = adam_id_str
        super(InstallVPPAbaClocKUnmanaged, self).__init__(udid)

    def _request_url(self):
        return config.MDM.server_url + "/v1/commands"

    def _request_type(self):
        return "InstallApplication"

    def _request_data(self):
        # raw:
        # '{"udid":"' + self.udid + '","request_type":"Settings","settings":[{"item":"Bluetooth","enabled":true},{"item":"DeviceName","device_name":"' + self.device_name + '"}]}'
        return {
            "udid": self.udid,
            "request_type": self._request_type(),
            "itunes_store_id": self.adam_id_str,
            "options": {
                "purchase_method": 1,
            }
        }

    def _command_identifier(self):
        return "install_vpp_abaclock_unmanaged"

    def _command_id(self):
        return 5