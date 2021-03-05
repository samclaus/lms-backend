# lms-backend
Learning Management System backend written in Golang

This branch has gotten too big lol.

# Changes implemented in this branch

## New request structure for websocket requests from the frontend
Go Implementation of generic WebSocketRequest struct. All Request data must have keys and values be strings.
```go
type WebSocketRequest struct {
    RequestType string `json:"type"`
    RequestData map[string]string `json:"data"`
}
```
Current RequestTypes: "login", "newuser"

Example json for login
```json
{
    "type": "login",
    "data": {
        "username": "username123",
        "password": "password123"
    }
}
```


## Separated functions for handling requests from server.go
Put request handler functions for each type of request into files based on their general category, eg. user.go for authentication and creating new users. Each function will have a name in the style of HandleLogin() and HandleNewUser(). Each function will accept the request object and a pointer to the Server as a parameter. HandleLogin currently returns a success bool, the user that attempted authentication, and a descriptive message string. It may be better to make a result type or use errors instead of a bool/message combination.

Structs will be defined in the file containing relevant handler functions. For example, Server and Request are in server.go, while Login and UserInfo are in user.go

## Created debug/cache.sqlite3
cache.sqlite3 currently has one table, `user_infos`. The fields are as follows:
```sql
CREATE TABLE "user_infos" (
	"id"	INTEGER NOT NULL UNIQUE,
	"username"	TEXT NOT NULL UNIQUE,
	"salt"	BLOB NOT NULL,
	"rounds"	INTEGER NOT NULL,
	"password_hash"	BLOB NOT NULL,
	PRIMARY KEY("id" AUTOINCREMENT)
)
```
There is currently one user in the database.
- id: 1
- username: username123
- salt: 11111
- rounds: 49283
- password_hash: [binary jumble here]

The unhashed password is password123

## Created debug/debug.html
debug.html is a simple webpage that sends a login request to the server and logs all ws messages to the console. Take a look at the source for more information. debug.html is served at the ["/debug" path](localhost:8080/debug).