package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// Gattributes ...
type Gattributes struct {
	AlternativeBankAccountNames interface{} `json:"alternative_bank_account_names"`
	BankID                      string      `json:"bank_id"`
	BankIDCode                  string      `json:"bank_id_code"`
	BaseCurrency                string      `json:"base_currency"`
	Bic                         string      `json:"bic"`
	Country                     string      `json:"country"`
}

// Gdata ...
type Gdata struct {
	Gattributes    Gattributes `json:"attributes"`
	CreatedOn      time.Time   `json:"created_on"`
	ID             string      `json:"id"`
	ModifiedOn     time.Time   `json:"modified_on"`
	OrganisationID string      `json:"organisation_id"`
	Type           string      `json:"type"`
	Version        int         `json:"version"`
}

// Glinks ...
type Glinks struct {
	Self string `json:"self"`
}

// GetAccountResponse  ...
type GetAccountResponse struct {
	Gdata  Gdata  `json:"data"`
	Glinks Glinks `json:"links"`
}

// GetAccount Calls the form3 api with the specified accountID
func GetAccount(host, accountID string) (statusCode int, account *GetAccountResponse) {
	fmt.Println("in GetAccount, id", accountID)

	uri := "/v1/organisation/accounts/" + accountID

	response, err := http.Get(host + uri)

	if err != nil {
		fmt.Print(err.Error())
		return
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	json.Unmarshal(responseData, &account)

	fmt.Println("Got account", account)

	return response.StatusCode, account

}
