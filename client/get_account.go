package client

import (
	"fmt"
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

// GetAccount calls the form3 api with the specified accountID
func GetAccount(host, accountID string) (*http.Response, error) {
	fmt.Println("in GetAccount, id", accountID)

	uri := "/v1/organisation/accounts/" + accountID

	request, err := RequestCreator(http.MethodGet, host+uri, nil)
	if err != nil {
		fmt.Print(err)
		return nil, err
	}

	client := &http.Client{}
	return client.Do(request)
}

// UnmarshallGetAccountResponse returns the  GetAccountResponse struct from the http.Response
func UnmarshallGetAccountResponse(response *http.Response) (account *GetAccountResponse, err error) {

	byteArr, err := IOResponseBodyReader(response.Body)
	if err != nil {
		fmt.Print(err)
		return nil, err
	}

	account = &GetAccountResponse{}
	err = Unmarshaller(byteArr, &account)
	if err != nil {
		fmt.Print(err)
		return nil, err
	}

	fmt.Println("Created account", account)

	return account, nil
}
