# Simple User Authentication
> Written in Golang
## Installing / Getting started

- [Go](https://golang.org/)
- [Git](https://git-scm.com/)
```
1.git clone 
2.cd ./simple-userAuth
3.go mod download
4.go run main.go
```

## Configuration
 - Go to file main.go
 - Find line 41 and edit it.
 
 Example:
 ```
 41. db, err = sqlx.Open("mysql", "root:123456@/user")
 ```
 
 ## Features
 - Register and Login Systems
 - JWT (JSON WEBTOKEN)
 
 ## Links
 - [Golang Documentation](https://golang.org/doc/)
 - [MySQL Documentation](https://dev.mysql.com/doc/)
 - [JWT](https://dev.mysql.com/doc/)
 - API ?
   - [WHAT IS API ?](https://en.wikipedia.org/wiki/API)
   - [REST API](https://en.wikipedia.org/wiki/Representational_state_transfer)
 
## Licensing
"Tho project is licensed under MIT license."
