class MDMOSUpdateStatus ( object ):
    def __init__(self, udid, plist):
        self.udid = udid
        self.download_percent_complete = plist['DownloadPercentComplete']
        self.is_downloaded = plist['IsDownloaded']
        self.product_key = plist['ProductKey']
        self.status = plist['Status']
