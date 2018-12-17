import base64
import json
import sys
import plistlib
import os
import DeviceSetup
import VPPAssociate

#sys.path.append(os.path.dirname(os.path.dirname(os.path.abspath(__file__))))
sys.path.insert(1, os.path.join(sys.path[0], '..'))
# import ../config.py
import config

def handleRequest(json):
    if 'acknowledge_event' in json:
        raw_payload = json['acknowledge_event']['raw_payload']
        payload = base64.b64decode(raw_payload).decode('utf-8')
        print(payload)
    elif 'checkin_event' in json:
        print("New Device was registered")
        raw_payload = json['checkin_event']['raw_payload']
        payload = base64.b64decode(raw_payload).decode('utf-8')
        print(payload)
        parseResponseJSON(json)
        return ''


def parseResponseString(request):
    response_json = json.loads(request)
    parseResponseJSON(response_json)


def parseResponseJSON(response_json):
    # mdm.Authenticate is very first request of a new iPad -> Setup
    print(response_json)
    if 'topic' in response_json and response_json['topic'] == 'mdm.Authenticate':
    	print("Response Json Topic: ", response_json['topic'])
        raw_payload = response_json['checkin_event']['raw_payload']
        payload_xml = base64.b64decode(raw_payload).decode('utf-8')

        if sys.version_info >= (3, 0):
            pl = plistlib.readPlistFromBytes(payload_xml.encode());
        else:
            pl = plistlib.readPlistFromString(payload_xml);

        if 'UDID' in pl:
            udid = pl['UDID']
    
            if 'SerialNumber' in pl:
            	serial_number = pl['SerialNumber']
            	VPPAssociate.VPPAssociate(udid, serial_number, vpp_associate_completed)
            	DeviceSetup.DeviceSetup(udid, setup_completed)
#             elif 'UnlockToken' in pl:
# 	        	# Device was setup
			    

    else:
        print(json)


def vpp_associate_completed(udid, serial_number, error):
    if error is not None:
        print("ERROR OCCURED")
        print("setup_completed", udid, serial_number, error)
    else:
        print("setup_completed", udid, serial_number, error)


def setup_completed(udid, error):
    if error is not None:
        print("ERROR OCCURED")
        print("setup_completed", udid, error)
    else:
        print("setup_completed", udid, error)






if __name__ == '__main__':
    print("This module is not made for running. only use for debugging.")
    print(config.MDM.server_url)
    print(config.Mysql.service_host)
    
    #VPPAssociate.VPPAssociate("5852cd48347416928de781b6c3f756696e0dcb31", "F9FWF8SWGHKL", vpp_associate_completed)
    #DeviceSetup.DeviceSetup("5852cd48347416928de781b6c3f756696e0dcb31", setup_completed)
