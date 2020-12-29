package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Account ...
type Account struct {
	Cdata Cdata `json:"data"`
}

// Cattributes ...
type Cattributes struct {
	Country      string `json:"country"`
	BaseCurrency string `json:"base_currency"`
	BankID       string `json:"bank_id"`
	BankIDCode   string `json:"bank_id_code"`
	Bic          string `json:"bic"`
}

// Cdata ...
type Cdata struct {
	Type           string      `json:"type"`
	ID             string      `json:"id"`
	OrganisationID string      `json:"organisation_id"`
	Cattributes    Cattributes `json:"attributes"`
}

// CreateRequestBody Creates a struct of type Account
func CreateRequestBody(accountID, organisationID string) (account *Account) {

	account = &Account{
		Cdata: Cdata{
			Type:           "accounts",
			ID:             accountID,
			OrganisationID: organisationID,
			Cattributes: Cattributes{
				Country:      "GB",
				BaseCurrency: "GBP",
				BankID:       "400300",
				BankIDCode:   "GBDSC",
				Bic:          "NWBKGB22",
			},
		},
	}
	return account
}

// CreateAccount Calls the form3 api with the specified accountID and organizationID
func CreateAccount(host, accountID, organisationID string) (statusCode int, createdAccount *Account) {

	fmt.Println("in CreateAccount.go")

	account := CreateRequestBody(accountID, organisationID)

	fmt.Println("account to create with id", account.Cdata.ID)

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(account)

	uri := "/v1/organisation/accounts"

	res, err := http.Post(host+uri, "application/vnd.api+json", b)

	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
		return
	}

	byteArr, _ := ioutil.ReadAll(res.Body)

	createdAccount = &Account{}
	json.Unmarshal(byteArr, &createdAccount)
	fmt.Println("Created account", createdAccount)

	return res.StatusCode, createdAccount

}
