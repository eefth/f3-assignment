package client

import (
	"bytes"
	"fmt"
	"net/http"
)

// Account ...
type Account struct {
	Cdata Cdata `json:"data"`
}

// Cattributes ...
type Cattributes struct {
	Country                 string   `json:"country"`
	BaseCurrency            string   `json:"base_currency"`
	BankID                  string   `json:"bank_id"`
	BankIDCode              string   `json:"bank_id_code"`
	Bic                     string   `json:"bic"`
	Name                    []string `json:"name"`
	AlternativeNames        []string `json:"alternative_names"`
	AccountClassification   string   `json:"account_classification"`
	JointAccount            bool     `json:"joint_account"`
	AccountMatchingOptOut   bool     `json:"account_matching_opt_out"`
	SecondaryIdentification string   `json:"secondary_identification"`
}

// Cdata ...
type Cdata struct {
	Type           string      `json:"type"`
	ID             string      `json:"id"`
	OrganisationID string      `json:"organisation_id"`
	Cattributes    Cattributes `json:"attributes"`
}

// CreateRequestBody creates a struct of type Account
func CreateRequestBody(accountID, organisationID string) (account *Account) {

	account = &Account{
		Cdata: Cdata{
			Type:           "accounts",
			ID:             accountID,
			OrganisationID: organisationID,
			Cattributes: Cattributes{
				Country:                 "GB",
				BaseCurrency:            "GBP",
				BankID:                  "400300",
				BankIDCode:              "GBDSC",
				Bic:                     "NWBKGB22",
				Name:                    []string{"Samantha Holder"},
				AlternativeNames:        []string{"Sam Holder"},
				AccountClassification:   "Personal",
				JointAccount:            false,
				AccountMatchingOptOut:   false,
				SecondaryIdentification: "A1B2C3D4",
			},
		},
	}
	return account
}

// CreateAccount calls the form3 api with the specified accountID and organizationID
func CreateAccount(host string, account *Account) (*http.Response, error) {

	fmt.Println("in CreateAccount.go with id", account.Cdata.ID)

	jsonBytes, err := Marshaller(account)
	if err != nil {
		return nil, err
	}

	uri := "/v1/organisation/accounts"

	request, err := RequestCreator(http.MethodPost, host+uri, bytes.NewReader(jsonBytes))
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/vnd.api+json")
	client := &http.Client{}
	return client.Do(request)
}

// UnmarshallCreateAccountResponse returns the  Account struct from the http.Response
func UnmarshallCreateAccountResponse(response *http.Response) (*Account, error) {

	byteArr, err := IOResponseBodyReader(response.Body)
	if err != nil {
		return nil, err
	}

	createdAccount := &Account{}
	err = Unmarshaller(byteArr, &createdAccount)
	if err != nil {
		return nil, err
	}

	fmt.Println("Created account", createdAccount)

	return createdAccount, nil
}
