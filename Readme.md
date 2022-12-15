# Smartio test assignment

MQTT on por 1883

Server on port 8000

Routes:

* GET /devices - gets the list of all devices
* GET /device/{id} - gets the specific device
* POST /device - creates the device, using _input 1_
* POST /devices - creates the list of devices, using the _input 2_
* DELETE /device/{id}, deletes the device by id
* DELETE /devices, deletes devices by ids using _input 2_

Ref:

* _input 1_:
```json
{
	"id": "newdevice10000"
}
```

* _input 2_:
```json
[
	{
		"id": "newdevice10000"
	},
	{
		"id": "newdevice10001"
	}
]
```
