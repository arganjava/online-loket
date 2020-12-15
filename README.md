# online-loket v1.0.0


## Installation and Setup

- Postgresql Installation
- Golang Version 1.14


- Provided Environment Variable in your system
  - TICKET_DB_HOST=localhost;
  - TICKET_DB_PORT=5432
  - TICKET_DB_PASSWORD=?
  - TICKET_DB_USERNAME=?
  - TICKET_DB_NAME=tiketing
  - TICKET_APP_PORT=:8282 or provided specific port that not used in your system
  
### Clone

- Clone this repo to your local machine using `https://github.com/arganjava/online-loket`

### How To run App
- go to root source code
- type go mod download
- enter src folder then type `go run main.go` enter
- if no error the apps will run on port 8282 or specific port you provided
---

### How To run Test
- go to `src/test` type `go test -v -cover .` enter
---

## Technical Steps

- add new db repository in `src/models` folder 
- add request and response in `src/dto` folder 
- add new interface business logic in `src/interface` folder
- add new implement business logic on `src/service` folder
- create unit test to cover business logic on `src/test` folder
- add new route `src/routes` folder
- register all route in `src/routes/server.go` file
- main configuration in `src/main.go` file 
- postman documentation in `Online Loket.postman_collection.json` file 

## Author

> Argan Megariansyah arganjava@gmail.com **[Linkedin](https://www.linkedin.com/in/argan-megariansyah-65751a89/)**