import requests
import sys
import os
import json
import base64

#sys.path.append(os.path.dirname(os.path.dirname(os.path.abspath(__file__))))
#sys.path.insert(1, os.path.join(sys.path[0], '../..'))
sys.path.insert(1, os.path.join(sys.path[0], '..'))
# import ../../config.py
import db

class MDMCommand (object):

    """Command is an abstract class that passes serialization and deserialization
    responsibility to its implementation classes, but takes care of storing the
    data to Mysql database."""
    def __init__(self, udid):
        """
        Parameters
        ------
        udid : str
            udid of the device (40 characters long) and lower cased
        """
        self.udid = udid
        self.command_uuid = None
        self.request()

    def _request_url(self):
        """
        Raises
        ------
        NotImplementedError
            If implementation class does not implement this method

        Returns
        -------
        str
            a string representing the api command url subpath of our MDM
            The prefix (server_url) goes WITHOUT trailing slash
            Only the suffix / path of the command is returned
            e.g.: config.MDM.server_url + "/v1/commands"
        """
        raise NotImplementedError("Should have implemented this: _request_url")

    def _request_type(self):
        """
        Raises
        ------
        NotImplementedError
            If implementation class does not implement this method

        Returns
        -------
        str
            a string representing the request_type of the command
            request_types are predefined by Apple.
            See: Mobile Device Management Protocol Reference
            https://developer.apple.com/business/documentation/MDM-Protocol-Reference.pdf
        """
        raise NotImplementedError("Should have implemented this: _request_type")

    def _request_data(self):
        """
        Raises
        ------
        NotImplementedError
            If implementation class does not implement this method

        Returns
        -------
        dict
            a dict, which will be transferred to json string representing the command request payload
        """
        raise NotImplementedError("Should have implemented this: _request_data")

    def command_identifier(self):
        """
        Raises
        ------
        NotImplementedError
            If implementation class does not implement this method

        Returns
        -------
        str
            a unique identifier string used for storing the command in the database
        """
        raise NotImplementedError("Should have implemented this: _command_identifier")

    def command_id(self):
        """
        Raises
        ------
        NotImplementedError
            If implementation class does not implement this method

        Returns
        -------
        int
            a unique id int used for storing the command in the database
        """
        raise NotImplementedError("Should have implemented this: _command_id")

    def _serialize_response(self, json):
        """
        Parameters
        ------
        json : str
            json being returned as response of the implemented post request

        Returns
        -------
        dict
            Dictionary of Key,Values, which are used for insertion into Mysql Database
            https://dev.mysql.com/doc/connector-python/en/connector-python-example-cursor-transaction.html
        """
        k_payload = 'payload'
        k_command_uuid = 'command_uuid'
        if k_payload in json and k_command_uuid in json[k_payload]:
            payload = json[k_payload]
            self.command_uuid = payload[k_command_uuid].encode("utf-8")
            # Used for Mysql insertion
            return {
                'command_uuid': self.command_uuid,
                'device_udid': self.udid,
                #'device_serial_number': self.serial_number,
                'command_identifier': self.command_identifier(),
                'command_id': self.command_id(),
            }
        else:
            raise ValueError("Given JSON not in expected format"+json)

    def _profile_payload(self, profile_file_path):
        profilePayload = open(profile_file_path, 'r').read()
        profilePayload = base64.b64encode(profilePayload)
        return profilePayload

    def request(self):
        data = json.dumps(self._request_data())
        url = self._request_url()
        c_identifier = self.command_identifier()
        c_id = self.command_id()
        headers = {
            'Content-Type': 'application/json',
        }
        response = requests.post(url, data=data, headers=headers, auth=('micromdm', 'secret'))
        if 'command_uuid' in response.json():
            self.command_uuid = response.json()['command_uuid']
        print(c_id, c_identifier, self.udid, response)
        response_serialized = self._serialize_response(response.json())

        if 'command_uuid' in response_serialized:
            self.command_uuid = response_serialized['command_uuid']

        db.DB.DB.log_command_request(self)
        print(response_serialized)