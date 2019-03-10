# Sample REST API

This is a sample REST API written in Go, using the Gin framework.

## Building

 1. Run `go get -d ./...` to get all project dependencies
 2. Run `go build` in the project root
 
## Running
 1. Run `go-restexample.exe` (windows) or `go-restexample`(unix/osx)
 2. The project will then run on `localhost:8080`

## Tests
TODO

## Consuming
Below is documentation for the various endpoints implemented within this project and how to consume them.

### POST /certificates
Endpoint used to create a new certificate

|Field|Value  |
|--|--|
|Headers|`Authorization: token`|
|Content-Type|`application/json`|
|Request Fields|`title`: Work title, string, required<br/>`year`: Year of work, integer, required<br/>`note`: Note regarding the work, string, optional|
|Response (success)|`status`: `success`<br/>`id`: Certificate ID|
|Response (error)|`status`: `failed`<br/>`error`: Error|

### PUT /certificates/:id
Endpoint used to update a certificate

|Field|Value  |
|--|--|
|Headers|`Authorization: token`|
|Content-Type|`application/json`|
|URL Parameters|`id`: Certificate ID, string, required|
|Request Fields|`title`: Work title, string, required<br/>`year`: Year of work, integer, required<br/>`note`: Note regarding the work, string, optional|
|Response (success)|`status`: `success`<br/>`id`: Certificate ID|
|Response (error)|`status`: `failed`<br/>`error`: Error|

### DELETE /certificates/:id
Endpoint used to delete a certificate

|Field|Value  |
|--|--|
|Headers|`Authorization: token`|
|Content-Type|`application/json`|
|URL Parameters|`id`: Certificate ID, string, required|
|Response (success)|`status`: `success`<br/>`id`: Certificate ID|
|Response (error)|`status`: `failed`<br/>`error`: Error|

### GET /users/:userId/certificates
Endpoint used to get the certificates for a user

|Field|Value  |
|--|--|
|Headers|`Authorization: token`|
|Content-Type|`application/json`|
|URL Parameters|`userId`: User ID, string, required|
|Response (success)|`status`: `success`<br/>`certificates`: Array of users certificates|
|Response (error)|`status`: `failed`<br/>`error`: Error|

### POST /certificates/:id/transfers
Endpoint used to create a new certificate transfer

|Field|Value  |
|--|--|
|Headers|`Authorization: token`|
|Content-Type|`application/json`|
|URL Parameters|`id`: Certificate ID, string, required|
|Request Fields|`to`: Email of recipient, string, required|
|Response (success)|`status`: `success`|
|Response (error)|`status`: `failed`<br/>`error`: Error|

### PUT /certificates/:id/transfers
Endpoint used to accept a certificate transfer

|Field|Value  |
|--|--|
|Headers|`Authorization: token`|
|Content-Type|`application/json`|
|URL Parameters|`id`: Certificate ID, string, required|
|Request Fields|`to`: Email of recipient, string, required<br/>`status`: Updated status of the transfer (`accepted` or `rejected`), string, required|
|Response (success)|`status`: `success`|
|Response (error)|`status`: `failed`<br/>`error`: Error|