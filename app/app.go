package main

import (
	"fmt"

	guuid "github.com/google/uuid"

	"github.com/eefth/f3-assignment/client"
)

const host = "http://localhost:8080"
const pageSize = 6

func main() {
	fmt.Println("Program is starting")

	var accountID, organisationID string
	accountID = guuid.New().String()
	organisationID = guuid.New().String()

	// create the account
	createAccountStatusCode, createAccountResponse := client.CreateAccount(host, accountID, organisationID)
	fmt.Printf("Create Account response status code %d\nCreated AccountId %s", createAccountStatusCode, createAccountResponse.Cdata.ID)

	// fetch the account
	getAccountStatusCode, getAccountResponse := client.GetAccount(host, accountID)
	fmt.Printf("Get Account response status code %d\nFetched AccountId %s", getAccountStatusCode, getAccountResponse.Gdata.ID)

	// get all existing accounts in db and print them out
	accounts := client.GatherAccounts(host, pageSize)
	fmt.Printf("No of accouns in db:%d\n", len(accounts))
	for _, d := range accounts {
		fmt.Println(d.Type, d.ID, d.OrganisationID, d.Version, d.Attributes.Country, d.Attributes.BaseCurrency)
	}

	// delete an account
	deleteAccountStatusCode := client.DeleteAccount(host, "b483e082-9b9e-4362-b2e1-69ddc0fc5b20", 0)
	fmt.Printf("Delete account response status code %d", deleteAccountStatusCode)

}
