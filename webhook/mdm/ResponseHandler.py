import base64
import json
import sys
import plistlib
import os
import DeviceSetup
import VPPAssociate
import db.MysqlDB

sys.path.insert(1, os.path.join(sys.path[0], '..'))
import config


def handleRequest(json):
    if 'checkin_event' in json:
        print("New Device was registered")

    parseResponseJSON(json)
    return ''


def parseResponseString(request):
    response_json = json.loads(request)
    parseResponseJSON(response_json)


def pl_payload(response_json, key):
    raw_payload = response_json[key]['raw_payload']
    payload_xml = base64.b64decode(raw_payload).decode('utf-8')

    if sys.version_info >= (3, 0):
        pl = plistlib.readPlistFromBytes(payload_xml.encode());
    else:
        pl = plistlib.readPlistFromString(payload_xml);

    return pl

def parseResponseJSON(response_json):
    # mdm.Authenticate is very first request of a new iPad -> Setup
    print(response_json)
    if 'topic' in response_json and response_json['topic'] == 'mdm.Authenticate':
        print("Response Json Topic: ", response_json['topic'])
        pl = pl_payload(response_json, 'checkin_event')
        if 'UDID' in pl:
            udid = pl['UDID']

            if 'SerialNumber' in pl:
                serial_number = pl['SerialNumber']
                VPPAssociate.VPPAssociate(udid, serial_number, vpp_associate_completed)
                DeviceSetup.DeviceSetup(udid, setup_completed)

    elif 'topic' in response_json and response_json['topic'] == 'mdm.Connect'\
            and 'acknowledge_event' in response_json:
        udid = response_json['acknowledge_event']['udid'].decode('utf-8')
        command__uuid = response_json['acknowledge_event']['command_uuid'].decode('utf-8')
        status = response_json['acknowledge_event']['status'].decode('utf-8')
        from db import DB
        DB.DB.log_command_response(udid, command__uuid, status)
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
    #dbHandler = db.MysqlDB.MysqlDB.log_command_request(None)
    #print("DB HANDLER", dbHandler)


    # VPPAssociate.VPPAssociate("5852cd48347416928de781b6c3f756696e0dcb31", "F9FWF8SWGHKL", vpp_associate_completed)
    #DeviceSetup.DeviceSetup("5852cd48347416928de781b6c3f756696e0dcb31", setup_completed)