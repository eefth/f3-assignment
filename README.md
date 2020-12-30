# f3-assignment
My name is Efthimios Efthimiou. I am new to go as my experience is one week old only. During that week, i saw how much exciting is to code with go. 
Comparing to java, a field that i work more than 10 years, i see that go is a more modern language, as less configuration is needed (recall maven, dependencies, spring e.t.c).
Also, go is faster and as people say, go is super in concurrency.

## project architecture
There are 2 folders containing the src: app and client. The app folder contains the main method, which calls the functions containing in the client folder.
The client folder contians the functions which call the form3 api. Theres is also a *test.go file there.

The tests exist in account_test.go. There the local form3 api has been mocked using the so called mux server.
In some cases json.Marshall and Unmarshall are mocked too. Current test-coverage is about 82.7%.

You can run the main method to create, delete, list, get account from the form3 api served by the docker-compose file. 

The dockerCompose folder is a copy of the form3tech-oss/interview-accountapi but with the docker-compose.yml file updated as requested.
I have included also the Dockerfile used to build the image that docker-compose.yml will use.

## how to run the docker-compose and see the go tests running
From the folder dockerCompose run the following: docker-compose up

## how to run the application
From the folder app run the following: go run app.go

## how to run the tests
From the folder client run the following: go test
