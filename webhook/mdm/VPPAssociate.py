import base64
import sys
import requests
import os
from mdm_commands import *

sys.path.append(os.path.dirname(os.path.dirname(os.path.abspath(__file__))))
sys.path.insert(1, os.path.join(sys.path[0], '..'))
#import ../config.py
import config

class VPPAssociate:
    def __init__(self, udid, serial_number, vpp_associate_completed):
        udid = str(udid)
        if not udid.islower():
            print("UDID should be provided as lower cased")
            udid = udid.lower()

        serial_number = str(serial_number)
        if not serial_number .isupper():
            print("serial_number should be provided as lower cased")
            serial_number = serial_number.upper()
        
        try:
            VPPAssociateCommand(udid, serial_number, config.MDM.adam_id_str)
            vpp_associate_completed(udid, serial_number, None)
        except:
            vpp_associate_completed(udid, serial_number, sys.exc_info()[0])