class MDMProfile ( object ):
    def __init__(self, udid, plist):
        self.udid = udid
        self.payload_uuid = plist['PayloadUUID']

        if 'PayloadDescription' in plist:
            self.payload_description = plist['PayloadDescription']
        else:
            self.payload_description = None

        self.has_removal_passcode = plist['HasRemovalPasscode']
        self.payload_identifier = plist['PayloadIdentifier']
        self.payload_display_name = plist['PayloadDisplayName']
        self.is_managed = plist['IsManaged']
        self.is_encrypted = plist['IsEncrypted']
        self.payload_version = plist['PayloadVersion']

        if 'PayloadOrganization' in plist:
            self.payload_organization = plist['PayloadOrganization']
        else:
            self.payload_organization = None

        self.payload_removal_disallowed = plist['PayloadRemovalDisallowed']

        self.profile_payloads = []
        if 'PayloadContent' in plist:
            for payload_content in plist['PayloadContent']:
                profile_payload = MDMProfilePayload(self.udid, self.payload_uuid, payload_content)
                self.profile_payloads.append(profile_payload)

        print(self.profile_payloads)


class MDMProfilePayload ( object ):
    def __init__(self, udid, profile_uuid, plist):
        self.udid = udid
        self.profile_uuid = profile_uuid

        self.payload_version = plist['PayloadVersion']

        if 'PayloadDescription' in plist:
            self.payload_description = plist['PayloadDescription']
        else:
            self.payload_description = None

        if 'PayloadType' in plist:
            self.payload_type = plist['PayloadType']
        else:
            self.payload_type = None

        self.payload_identifier = plist['PayloadIdentifier']

        if 'PayloadDisplayName' in plist:
            self.payload_display_name = plist['PayloadDisplayName']
        else:
            self.payload_display_name = None
