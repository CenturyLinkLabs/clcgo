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
		"status": "active",
		"details": {
			"powerState": "started",
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

	SuccessfulStatusResponse = `{ "status":"succeeded" }`

	ServerCreationSuccessfulResponse = `{
		"server":"web",
		"isQueued":true,
		"links":[
			{
				"rel":"status",
				"href":"/v2/operations/alias/status/test-status-id",
				"id":"test-status-id"
			},
			{
				"rel":"self",
				"href":"/v2/servers/alias/test-uuid?uuid=True",
				"id":"test-uuid",
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
