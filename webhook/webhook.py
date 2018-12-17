from flask import Flask, request, abort
import base64

app = Flask(__name__)

@app.route('/webhook', methods=['POST'])
def webhook():
	print(request.json)
	if 'acknowledge_event' in request.json:
		raw_payload = request.json['acknowledge_event']['raw_payload']
		payload = base64.b64decode(raw_payload).decode('utf-8')
		print(payload)
	elif 'checkin_event' in request.json:
		print("New Device was registered")
		raw_payload = request.json['checkin_event']['raw_payload']
		payload = base64.b64decode(raw_payload).decode('utf-8')
		print(payload)
	return ''
	

if __name__ == '__main__':
    app.run()
    #DeviceSetup.DeviceSetup("5852cd48347416928de781b6c3f756696e0dcb31", "F9FWF8SWGHKL", setup_completed)