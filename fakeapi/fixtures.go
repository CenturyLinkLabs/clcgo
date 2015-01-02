package fakeapi

// Fixture constants for the fakeapi's hardcoded responses, including for
// successful and unsuccessful requests.
const (
	AuthenticationSuccessfulResponse = `{
		"bearerToken": "1234ABCDEF",
		"accountAlias": "ACME"
	}`

	ServerResponse = `{
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

	ServerCreationSuccessfulResponse = `{
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

	ServerCreationMissingStatusResponse = `{
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

	ServerCreationInvalidResponse = `{
		"message": "The request is invalid.",
		"modelState": {
			"body.name": ["The name field is required."],
			"body.sourceServerId":["The sourceServerId field is required."]
		}
	}`

	AddPublicIPAddressSuccessfulResponse = `{
		"rel":"status",
		"href":"/path/to/status",
		"id":"id-for-status"
	}`

	DataCenterCapabilitiesResponse = `{
		"templates": [
			{ "name": "Name", "description": "Description" }
		]
	}`
)
