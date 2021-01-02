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
	account := client.CreateRequestBody(accountID, organisationID)
	createAccountResponse, _ := client.CreateAccount(host, account)
	createdAccount, _ := client.UnmarshallCreateAccountResponse(createAccountResponse)
	fmt.Printf("Created Account with AccountId %s", createdAccount.Cdata.ID)

	// fetch the account
	getAccountResponse, _ := client.GetAccount(host, accountID)
	existingAccount, _ := client.UnmarshallGetAccountResponse(getAccountResponse)
	fmt.Printf("Get Existing Account with AccountId %s", existingAccount.Gdata.ID)

	// get all existing accounts in db and print them out
	accounts := client.GatherAccounts(host, pageSize)
	fmt.Printf("No of accouns in db:%d\n", len(accounts))
	for _, d := range accounts {
		fmt.Println(d.Type, d.ID, d.OrganisationID, d.Version, d.Attributes.Country, d.Attributes.BaseCurrency)
	}

	// delete an account
	deleteAccountResponse, _ := client.DeleteAccount(host, "b483e082-9b9e-4362-b2e1-69ddc0fc5b20", 0)
	fmt.Printf("Delete account response status code %d", deleteAccountResponse.StatusCode)

}
