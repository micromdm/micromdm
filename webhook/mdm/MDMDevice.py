class MDMDevice ( object ):
    def __init__(self, udid, plist):
        self.udid = udid
        self.serial_number = plist['SerialNumber']
        self.build_version = plist['BuildVersion']
        self.os_version = plist['OSVersion']
        self.product_name = plist['ProductName']
