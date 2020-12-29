package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// GetAccountsResponse ...
type GetAccountsResponse struct {
	Data []Data `json:"data"`
}

// Attributes ...
type Attributes struct {
	Country               string `json:"country"`
	BaseCurrency          string `json:"base_currency"`
	AccountNumber         string `json:"account_number"`
	BankID                string `json:"bank_id"`
	BankIDCode            string `json:"bank_id_code"`
	Bic                   string `json:"bic"`
	Iban                  string `json:"iban"`
	AccountClassification string `json:"account_classification"`
	JointAccount          bool   `json:"joint_account"`
	Switched              bool   `json:"switched"`
	AccountMatchingOptOut bool   `json:"account_matching_opt_out"`
	Status                string `json:"status"`
}

// Data ...
type Data struct {
	Type           string     `json:"type"`
	ID             string     `json:"id"`
	OrganisationID string     `json:"organisation_id"`
	Version        int        `json:"version"`
	Attributes     Attributes `json:"attributes"`
}

// ListAccounts Calls the form3 api with the specified pageNumber and pageSize
func ListAccounts(host string, pageNumber, pageSize int) (statusCode int, accounts *GetAccountsResponse) {
	fmt.Println("in ListAccounts", "pageNumber", pageNumber, "pageSize", pageSize)

	uri := "/v1/organisation/accounts?"

	response, err := http.Get(host + uri + "page[number]=" + fmt.Sprint(pageNumber) + "&page[size]=" + fmt.Sprint(pageSize))

	if err != nil {
		fmt.Println(err)
		return
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	json.Unmarshal(responseData, &accounts)

	return response.StatusCode, accounts
}

// GatherAccounts Gets the list of all existing accounts in db by calling
// the 'ListAccounts'
func GatherAccounts(host string, pageSize int) (allAccs []Data) {

	allAccs = make([]Data, 0)
	listAccountsStatusCode := 200

	for pageNumber := 0; listAccountsStatusCode == 200; pageNumber = pageNumber + 1 {
		listAccountsStatusCodeCurrent, accounts := ListAccounts(host, pageNumber, pageSize)

		listAccountsStatusCode = listAccountsStatusCodeCurrent

		if listAccountsStatusCode == 200 && len(accounts.Data) > 0 {
			for _, d := range accounts.Data {
				allAccs = append(allAccs, d)
			}
			if len(accounts.Data) < pageSize {
				break
			}
		} else {
			break
		}
	}

	return allAccs
}
