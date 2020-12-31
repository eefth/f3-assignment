package client

import (
	"fmt"
	"net/http"
)

// DeleteAccount calls the form3 api with the specified accountID and version
func DeleteAccount(host, accountID string, version int) (statusCode int) {

	fmt.Println("in DeleteAccount, id", accountID, "version", version)

	uri := "/v1/organisation/accounts/"

	// Create client
	client := &http.Client{}

	// Create request
	req, err := http.NewRequest("DELETE", host+uri+accountID+"?version="+fmt.Sprint(version), nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Fetch Request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	// Display Response code
	fmt.Println("response Status : ", resp.StatusCode)

	return resp.StatusCode
}
