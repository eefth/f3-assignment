# f3-assignment
My name is Efthimios Efthimiou. I am new to go as my experience is 10 days old only. During that week, i saw how much exciting is to code with go. 
Comparing to java, a field that i work more than 10 years, i see that go is a more modern language, as less configuration is needed (recall maven, dependencies, spring e.t.c).
Also, go is faster and as people say, go is super in concurrency.

## Project architecture
Source code resides in 2 folders: app and client. The app folder contains the main method, which calls the functions existing in the client folder.

The client folder contains the functions that call the form3 api. Also, some json marshalling and unmarshalling functions reside there together with the structs and other variables. More specifically:
### Package client
#### create_account.go
This file contains the functions used to create a form3 Account resource.
### get_account.go
This file contains the functions used to get a form3 Account resource based on the accountId
#### delete_account.go
This file contains the functions used to delete a form3 Account resource.
#### list_accounts.go
This file contains the functions used to list form3 Account resources with paging support.
#### inits.go
This file contains some initialization variables that wrap build in go functions. These variables can be used to mock those functions.
#### accounts_test.go
This file contains the tests. Inside that file, the local form3 api has been mocked using the so called mux server.
In some cases json.Marshall, json.Unmarshall, http.NewRequest and ioutil.ReadAll are mocked too. Currently the test-coverage is about 100%, a value got from the VS Code go extension api.

### Package main
### app.go
This file contains the main method, that is used to call the functions of the client package that is described above. You can run that file after the api is served from 'docker-compose up' 

## Note
Regarding the list accounts, the main method calls a helper method which calls then another method, both in the client package. The latter method calls the form3 api. At the end, all the existing accounts in db are fetched and printed in the main.go. The pageSize request parameter has been put very low and equal to 6 in such a way to require a few iterations(pages) in order to gather all accounts from db.

## Note
The dockerCompose folder is a copy of the form3tech-oss/interview-accountapi but with the docker-compose.yml file updated as requested.
The Dockerfile used to build the image that docker-compose.yml will take into account, is also included. As base image an golang-alpine one is used which is light weight. My image name is eefth/my-go-app and is pushed in Docker hub.

## How to run the docker-compose and see the go tests running
From the folder dockerCompose run the following: docker-compose up

## How to run the application
From the folder app run the following: go run app.go

## How to run the tests
Apart from watching the tests running when you do docker-compose up, you can also run them with the following ways
- From the folder client run the following: go test
- In the folder where the Dockerfile exists do:
docker build -t eefth/my-go-app . \
docker run eefth/my-go-app
