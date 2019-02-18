import sys
import os
#sys.path.append(os.path.dirname(os.path.dirname(os.path.abspath(__file__))))
sys.path.insert(1, os.path.join(sys.path[0], '../..'))
# import ../../config.py
import config

from MDMCommand import MDMCommand

class VPPAssociateCommand( MDMCommand ):
    def __init__(self, udid, serial_number, adam_id_str, pricing_oaram = 'STDQ'):
        """
        Parameters
        ------
        udid : str
            udid of the device (40 characters long) and lower cased
        serial_number: str
            Serial Number of the device (12 characters long) and upper cased
        adam_id_str: str
            id of the application, adam_id_str can be retrieved by the AppStore URL
            https://itunes.apple.com/ch/app/abaclock-2/id1301278477?mt=8
            For AbaClocK the adam_id_str would be "1301278477"
        pricing_oaram: str
            Either "PLUS" or "STDQ"
            https://developer.apple.com/business/documentation/MDM-Protocol-Reference.pdf
        """
        self.serial_number = serial_number
        self.adam_id_str = adam_id_str
        self.pricing_oaram = pricing_oaram
        super(VPPAssociateCommand, self).__init__(udid)

    def _request_url(self):
        # https://developer.apple.com/business/documentation/MDM-Protocol-Reference.pdf
        return "https://vpp.itunes.apple.com/WebObjects/MZFinance.woa/wa/manageVPPLicensesByAdamIdSrv"

    def _request_type(self):
        return ""

    def _request_data(self):
        return {
            "sToken": config.MDM.s_token,
            "adamIdStr": self.adam_id_str,
            "pricingParam": self.pricing_oaram,
            "associateSerialNumbers": [
                self.serial_number
            ]
        }

    def command_identifier(self):
        return "vpp_associate"

    def command_id(self):
        return 3

    def _serialize_response(self, json):
        print(json)
        print("VPPAssociate implements an own response serialization, as the VPPAsssociate deiffers from the others MDM reuqests, it directly communicates with Apple's VPP Server")
        return {
            'command_uuid': "",
            'device_udid': self.udid,
            'device_serial_number': self.serial_number,
            'command_identifier': self._command_identifier(),
            'command_id': self._command_id(),
        }