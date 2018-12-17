from flask import Flask, request, abort
import base64
import json
from mdm import ResponseHandler

app = Flask(__name__)

@app.route('/webhook', methods=['POST'])
def webhook():
    print(request.json)
    ResponseHandler.handleRequest(request.json)
    return ''



if __name__ == '__main__':
    app.run()
    
#     jsonObj = json.loads('{"topic": "mdm.Authenticate", "event_id": "c6a2d42d-63f5-4cb6-9c60-208a61cea179", "created_at": "2018-12-17T07:18:02.339214011Z", "checkin_event": {"udid": "08ad0661b9ab4cc9aad2badb8b3b64f106457df8", "url_params": null, "raw_payload": "PD94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz0iVVRGLTgiPz4KPCFET0NUWVBFIHBsaXN0IFBVQkxJQyAiLS8vQXBwbGUvL0RURCBQTElTVCAxLjAvL0VOIiAiaHR0cDovL3d3dy5hcHBsZS5jb20vRFREcy9Qcm9wZXJ0eUxpc3QtMS4wLmR0ZCI+CjxwbGlzdCB2ZXJzaW9uPSIxLjAiPgo8ZGljdD4KCTxrZXk+QnVpbGRWZXJzaW9uPC9rZXk+Cgk8c3RyaW5nPjE2QzUwNDNiPC9zdHJpbmc+Cgk8a2V5Pk1lc3NhZ2VUeXBlPC9rZXk+Cgk8c3RyaW5nPkF1dGhlbnRpY2F0ZTwvc3RyaW5nPgoJPGtleT5PU1ZlcnNpb248L2tleT4KCTxzdHJpbmc+MTIuMS4xPC9zdHJpbmc+Cgk8a2V5PlByb2R1Y3ROYW1lPC9rZXk+Cgk8c3RyaW5nPmlQYWQ1LDE8L3N0cmluZz4KCTxrZXk+U2VyaWFsTnVtYmVyPC9rZXk+Cgk8c3RyaW5nPkY5RlY1MTQ0R0hLSjwvc3RyaW5nPgoJPGtleT5Ub3BpYzwva2V5PgoJPHN0cmluZz5jb20uYXBwbGUubWdtdC5FeHRlcm5hbC41MTQyM2IyNi0wNTM4LTQyOTMtYTQ1Ny01OTkwNmQ2OTYzYmM8L3N0cmluZz4KCTxrZXk+VURJRDwva2V5PgoJPHN0cmluZz4wOGFkMDY2MWI5YWI0Y2M5YWFkMmJhZGI4YjNiNjRmMTA2NDU3ZGY4PC9zdHJpbmc+CjwvZGljdD4KPC9wbGlzdD4K"}}')
#     ResponseHandler.handleRequest(jsonObj)
    
#     jsonObj = json.loads('{"topic": "mdm.TokenUpdate", "event_id": "46380b5f-3472-4fad-a54c-ec2b7d3c6d1b", "created_at": "2018-12-17T07:18:03.773086664Z", "checkin_event": {"udid": "08ad0661b9ab4cc9aad2badb8b3b64f106457df8", "url_params": null, "raw_payload": "PD94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz0iVVRGLTgiPz4KPCFET0NUWVBFIHBsaXN0IFBVQkxJQyAiLS8vQXBwbGUvL0RURCBQTElTVCAxLjAvL0VOIiAiaHR0cDovL3d3dy5hcHBsZS5jb20vRFREcy9Qcm9wZXJ0eUxpc3QtMS4wLmR0ZCI+CjxwbGlzdCB2ZXJzaW9uPSIxLjAiPgo8ZGljdD4KCTxrZXk+QXdhaXRpbmdDb25maWd1cmF0aW9uPC9rZXk+Cgk8ZmFsc2UvPgoJPGtleT5NZXNzYWdlVHlwZTwva2V5PgoJPHN0cmluZz5Ub2tlblVwZGF0ZTwvc3RyaW5nPgoJPGtleT5QdXNoTWFnaWM8L2tleT4KCTxzdHJpbmc+NkQ5Njg5QzUtNkI1Mi00OTgwLThEODQtNDYzQjA4NjJGRTcxPC9zdHJpbmc+Cgk8a2V5PlRva2VuPC9rZXk+Cgk8ZGF0YT4KCUZKck1xdGNCaUVXb2tDM1hDS0xNTTcyNXAzYlIrS0JvVjh4QTMvUHpsUDQ9Cgk8L2RhdGE+Cgk8a2V5PlRvcGljPC9rZXk+Cgk8c3RyaW5nPmNvbS5hcHBsZS5tZ210LkV4dGVybmFsLjUxNDIzYjI2LTA1MzgtNDI5My1hNDU3LTU5OTA2ZDY5NjNiYzwvc3RyaW5nPgoJPGtleT5VRElEPC9rZXk+Cgk8c3RyaW5nPjA4YWQwNjYxYjlhYjRjYzlhYWQyYmFkYjhiM2I2NGYxMDY0NTdkZjg8L3N0cmluZz4KCTxrZXk+VW5sb2NrVG9rZW48L2tleT4KCTxkYXRhPgoJUkVGVVFRQUFCT1JXUlZKVEFBQUFCQUFBQUFWVVdWQkZBQUFBQkFBQUFBSlZWVWxFQUFBQUVJSitjRmNQdjBYZGtDbFIKCWlwQ3BvMk5JVFVOTEFBQUFLTGVZblJyaURLeittS1NBQ3hqN1FGdjNUaW9PQXlUSDNqMldvL0I1YzV4Z2FLUjVwVGEyCgl3UjFYVWtGUUFBQUFCQUFBQUFGVFFVeFVBQUFBRk5KUHZkNjRuekhYazU0Q3FGb1h2YkxDU0U0WVNWUkZVZ0FBQUFRQQoJQU1OUVZWVkpSQUFBQUJBOGd0RThjRDVDVWFqMGdmRjdodWNKUTB4QlV3QUFBQVFBQUFBTFYxSkJVQUFBQUFRQUFBQUIKCVMxUlpVQUFBQUFRQUFBQUFWMUJMV1FBQUFDQTB2U09PcTBHU25DL0lneklYUEYyOFNVV2d3aXhVV2RUYUFYV0g2WHYrCglRMVZWU1VRQUFBQVE2MEVxbTVXWFRPV1NKYUJ3K3kwck9VTk1RVk1BQUFBRUFBQUFDbGRTUVZBQUFBQUVBQUFBQTB0VQoJV1ZBQUFBQUVBQUFBQUZkUVMxa0FBQUFveExyVGE0WTZ3UUxWYUpnV3MwbWUzeHBFTEsva2s3S21YYWdsc25mTVc5bG4KCW96VGRZbHVEb0ZWVlNVUUFBQUFRL1gwUVkveVNTSEdkeVVLQk5oWTdCME5NUVZNQUFBQUVBQUFBQ1ZkU1FWQUFBQUFFCglBQUFBQTB0VVdWQUFBQUFFQUFBQUFGZFFTMWtBQUFBb0ZraFRNaDNnZlNad2FNUEdDb29peUd0OUhId3pCUGZaNnFaRAoJRnpHbTQxV2pXcDB3THV5NFlsVlZTVVFBQUFBUU90ZUFGNVNSUUV1RHpFeitvQWlseFVOTVFWTUFBQUFFQUFBQUNGZFMKCVFWQUFBQUFFQUFBQUFVdFVXVkFBQUFBRUFBQUFBRmRRUzFrQUFBQWd2azFsM0cvM1RtTFg3QllFYWVhbVFBaW1vaE5DCgk4YzF0UG4yZWdZdjkwNjVWVlVsRUFBQUFFSXU3aUIxalMwMUxpT1VJcHB1Z3oyTkRURUZUQUFBQUJBQUFBQWRYVWtGUQoJQUFBQUJBQUFBQU5MVkZsUUFBQUFCQUFBQUFCWFVFdFpBQUFBS0xJSkJvZkMvK1pLbE5xTldmR3VBL0p4aGlhU0dUaTYKCU51QmxmNHcrTS9Ka3R6OVRLSjVralp0VlZVbEVBQUFBRU9ZYVJzeUVRRWlNbkRZdXU5cUJLeEpEVEVGVEFBQUFCQUFBCglBQVpYVWtGUUFBQUFCQUFBQUFOTFZGbFFBQUFBQkFBQUFBQlhVRXRaQUFBQUtPQ0tNdllRaXVGSzdRQ1JHbnlySUhpawoJanNLVDJ3OHZ2VmwzYmRYVXpacFVVRFRBMkZJUGFiRlZWVWxFQUFBQUVFWC9aQ3phdkVyL3BGNnVLbkVET3NaRFRFRlQKCUFBQUFCQUFBQUFWWFVrRlFBQUFBQkFBQUFBTkxWRmxRQUFBQUJBQUFBQUJYVUV0WkFBQUFLSTAzMExnY0xMZGFKeWRKCgloOCtKMWZMamZ0V1BxcTRIT3hDa0gzVzdwMXhwVm1DVVBKakNZWmRWVlVsRUFBQUFFS1YwS01rVkpVWGJ0RzZaVTN3TQoJNjRWRFRFRlRBQUFBQkFBQUFBTlhVa0ZRQUFBQUJBQUFBQU5MVkZsUUFBQUFCQUFBQUFCWFVFdFpBQUFBS0ZnSUZPR0sKCUNIQWpETEJQVWovWFdON0lTUFlHeS8xVk1KRFVSK0FJYWp3ZFN2N25lZWxUZ2c1VlZVbEVBQUFBRU93MTByVW1zVUY2Cgl1U3Z4VmFGT0lEdERURUZUQUFBQUJBQUFBQUpYVWtGUUFBQUFCQUFBQUFOTFZGbFFBQUFBQkFBQUFBRlhVRXRaQUFBQQoJS05LU3I5NmNya0xJUW1QTmtQQ0hUcEg5bzdsR0hFSytCLzdyN2lINXdjQnloUEwyQ3JZOS9ZeFFRa3RaQUFBQUlCUmIKCU10NnRYMlcrQnBYOW1XeldDWkZrQmZQL1RMZnZVckVyZ0cyRSt5dGpWVlZKUkFBQUFCQkU3ZXhCZnp4TFY3d1NIdWN3CglpZXpsUTB4QlV3QUFBQVFBQUFBQlYxSkJVQUFBQUFRQUFBQURTMVJaVUFBQUFBUUFBQUFBVjFCTFdRQUFBQ2puVEdydwoJMXFwbnViU3JqaGtoeEg2M0t1S3REVjZ1SGE5MDVQdXVZUDhLRTlyUWZWZ2dUZDh3VTBsSFRnQUFBQlJHdEtLNkJXK08KCXBjYmF2Q2tjMUVVZVZnVHg0Zz09Cgk8L2RhdGE+CjwvZGljdD4KPC9wbGlzdD4K"}}')
#     ResponseHandler.handleRequest(jsonObj)

#	jsonObj = json.loads('{"topic": "mdm.Connect", "event_id": "016c5677-3c15-46ae-b946-4ebaf68d8d6c", "created_at": "2018-12-17T15:47:17.252785689Z", "acknowledge_event": {"udid": "5852cd48347416928de781b6c3f756696e0dcb31", "status": "Idle", "url_params": null, "command_uuid": "8951a276-e7ad-47c5-9336-ae20a53759c5", "raw_payload": "PD94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz0iVVRGLTgiPz4KPCFET0NUWVBFIHBsaXN0IFBVQkxJQyAiLS8vQXBwbGUvL0RURCBQTElTVCAxLjAvL0VOIiAiaHR0cDovL3d3dy5hcHBsZS5jb20vRFREcy9Qcm9wZXJ0eUxpc3QtMS4wLmR0ZCI+CjxwbGlzdCB2ZXJzaW9uPSIxLjAiPgo8ZGljdD4KCTxrZXk+U3RhdHVzPC9rZXk+Cgk8c3RyaW5nPklkbGU8L3N0cmluZz4KCTxrZXk+VURJRDwva2V5PgoJPHN0cmluZz4zMGY2ZGE5YTg1ZmYxYzRiOWJhMzZmYzFmMzM5ZDliOWIyZDQxNzI2PC9zdHJpbmc+CjwvZGljdD4KPC9wbGlzdD4K"}}')
#	ResponseHandler.handleRequest(jsonObj)