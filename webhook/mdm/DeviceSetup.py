import base64
import sys
import requests
import os
from mdm_commands import *

sys.path.append(os.path.dirname(os.path.dirname(os.path.abspath(__file__))))
sys.path.insert(1, os.path.join(sys.path[0], '..'))
#import ../config.py
import config

class DeviceSetup:
    def __init__(self, udid, setup_completed):
        udid = str(udid)
        if not udid.islower():
            print("UDID should be provided as lower cased")
            udid = udid.lower()

        try:
            Settings(udid, "New AbaClocK iPad")
            RestrictApps(udid, serial_number)
            HomescreenLayout(udid, serial_number)
            InstallVPPAbaClocKUnmanaged(udid, int(config.MDM.adam_id_str))
            ManagedApplicationConfiguration(udid, "ch.abacus.abaclock.client2")
            SingleAppModeAbaClocK(udid, serial_number)
            ProfileList(udid, serial_number)
            ManagedApplicationFeedback(udid, "ch.abacus.abaclock.client2")
            RemoveSingleAppModeAbaClocK(udid, serial_number)
            setup_completed(udid, None)
        except:
            setup_completed(udid, sys.exc_info()[0])