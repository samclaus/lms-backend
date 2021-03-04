# lms-backend
Learning Management System backend written in Golang

# Changes to be implemented in this fork

## New request structure for the frontend
Go Implementation of generic Request struct
```go
type Request struct {
    RequestType string `json:"type"`
    RequestData map[string]interface{} `json:"data"`
}
```
Current RequestTypes: "login", "newuser"

Example json for login
```json
{
    "type": "login",
    "data": {
        "username": "user123",
        "password": "password123"
    }
}
```


## Separate functions for handling requests from server.go
Put request handler functions for each type of request into files based on their general category, eg. user.go for authentication and creating new users. Each function will have a name in the style of HandleLogin() and HandleNewUser(). Each function will accept the request object as a parameter and return a marshalled json string that the main server function will send back to the client via ws.

Structs will be defined in the file containing relevant handler functions. For example, Server and Request are in server.go, while Login and UserInfo are in user.go