package clcgo

const (
	serverResponse = `{
		"id": "test-id",
		"name": "Test Name",
		"groupId": "123il"
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
