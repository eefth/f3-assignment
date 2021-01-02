package client

import (
	"fmt"
	"net/http"
)

// DeleteAccount calls the form3 api with the specified accountID and version
func DeleteAccount(host, accountID string, version int) (*http.Response, error) {

	fmt.Println("in DeleteAccount, id", accountID, "version", version)

	uri := "/v1/organisation/accounts/"

	// Create client
	client := &http.Client{}

	// Create request
	req, err := RequestCreator("DELETE", host+uri+accountID+"?version="+fmt.Sprint(version), nil)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	// Fetch Request
	return client.Do(req)
}
