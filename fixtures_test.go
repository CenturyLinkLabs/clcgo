package clcgo

const (
	serverResponse = `{
		"id": "test-id",
		"name": "Test Name",
		"groupId": "123il",
		"details": {
			"ipAddresses": [
				{
					"internal": "10.0.0.1"
				},
				{
					"public": "8.8.8.8",
					"internal": "10.0.0.2"
				}
			]
		}
	}`
	serverCreationSuccessfulResponse = `{
		"server":"web",
		"isQueued":true,
		"links":[
			{
				"rel":"status",
				"href":"/path/to/status",
				"id":"id-for-status"
			},
			{
				"rel":"self",
				"href":"/path/to/self",
				"id":"id-for-self",
				"verbs": [ "GET" ]
			}
		]
	}`

	serverCreationMissingStatusResponse = `{
		"server":"web",
		"isQueued":true,
		"links":[
			{
				"rel":"self",
				"href":"/path/to/self",
				"id":"id-for-self",
				"verbs": [ "GET" ]
			}
		]
	}`

	serverCreationInvalidResponse = `{
		"message": "The request is invalid.",
		"modelState": {
			"body.name": ["The name field is required."],
			"body.sourceServerId":["The sourceServerId field is required."]
		}
	}`

	addPublicIPAddressSuccessfulResponse = `{
		"rel":"status",
		"href":"/path/to/status",
		"id":"id-for-status"
	}`

	dataCenterCapabilitiesResponse = `{
		"templates": [
			{ "name": "Name", "description": "Description" }
		]
	}`
)
