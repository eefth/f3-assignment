package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	guuid "github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// newTestServer creates a multiplex server to handle API endpoints
func newTestServer(path string, h func(w http.ResponseWriter, r *http.Request)) *httptest.Server {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	mux.HandleFunc(path, h)
	return server
}

func TestCreateAccount_success(t *testing.T) {

	accountID := guuid.New().String()
	organizationID := guuid.New().String()
	account := CreateRequestBody(accountID, organizationID)

	uri := "/v1/organisation/accounts/"

	server := newTestServer(uri, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
	})

	defer server.Close()
	response, _ := CreateAccount(server.URL, account)

	msg := fmt.Sprintf("TestCreateAccount failed. Status code expected to be %d but it was %d", http.StatusOK, response.StatusCode)

	if response.StatusCode != 201 {
		t.Errorf(msg)
	}
}

func TestCreateAccount_whenMarshallerFails_shouldReturnError(t *testing.T) {

	accountID := guuid.New().String()
	organizationID := guuid.New().String()
	account := CreateRequestBody(accountID, organizationID)

	uri := "/v1/organisation/accounts/"

	server := newTestServer(uri, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
	})

	Marshaller = func(v interface{}) ([]byte, error) {
		return nil, errors.New("Marshaller exception")
	}

	defer server.Close()
	_, err := CreateAccount(server.URL, account)

	assert.NotNil(t, err)
	assert.EqualValues(t, fmt.Sprint(err), "Marshaller exception")
}

func TestGetAccount_whenForm3ApiReturnesOK(t *testing.T) {

	accountID := guuid.New().String()
	uri := "/v1/organisation/accounts/"

	var body = GetAccountResponse{Gdata: Gdata{Type: "accounts", ID: "0673746b-8dd3-4bd2-b398-941bdf2865df"}}

	server := newTestServer(uri, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(body)
	})
	defer server.Close()

	response, error := GetAccount(server.URL, accountID)

	assert.Nil(t, error)
	assert.EqualValues(t, response.StatusCode, 200)
}

func TestGetAccount_whenForm3ApiReturnes500_shouldReturn500(t *testing.T) {

	accountID := guuid.New().String()
	uri := "/v1/organisation/accounts/"

	server := newTestServer(uri, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	defer server.Close()

	getAccountResponse, err := GetAccount(server.URL, accountID)

	assert.Nil(t, err)
	assert.EqualValues(t, getAccountResponse.StatusCode, 500)
}

func TestListAccounts_whenForm3ApiReturnsOK(t *testing.T) {

	pageNumber := 1
	pageSize := 30
	uri := "/v1/organisation/accounts"

	server := newTestServer(uri, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	defer server.Close()
	statusCode, _ := ListAccounts(server.URL, pageNumber, pageSize)

	msg := fmt.Sprintf("TestListAccounts failed. Status code expected to be %d but it was %d", http.StatusOK, statusCode)

	if statusCode != 200 {
		t.Errorf(msg)
	}
}

func TestGatherAccounts_WhenListAccountsReturns500(t *testing.T) {

	pageSize := 30
	uri := "/v1/organisation/accounts"

	server := newTestServer(uri, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	defer server.Close()

	accounts := GatherAccounts(server.URL, pageSize)

	if len(accounts) > 0 {
		t.Errorf("Expected found accounts to be %v but got %v", 0, len(accounts))
	}
}

func TestGatherAccounts_WhenListAccountsReturns200AndTwoResults(t *testing.T) {

	pageSize := 30
	uri := "/v1/organisation/accounts"

	var body = GetAccountsResponse{Data: []Data{{Type: "accounts", ID: "0673746b-8dd3-4bd2-b398-941bdf2865df"},
		{Type: "accounts", ID: "0673746b-8dd3-4bd2-b398-941bdf2865df"}}}

	server := newTestServer(uri, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(body)

	})
	defer server.Close()

	accounts := GatherAccounts(server.URL, pageSize)

	if len(accounts) != 2 {
		t.Errorf("Expected found accounts to be %v but got %v", 2, len(accounts))
	}
}

func TestDeleteAccount_whenForm3ApiReturns204(t *testing.T) {

	accountID := guuid.New().String()
	version := 0
	uri := "/v1/organisation/accounts/"

	server := newTestServer(uri, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
	defer server.Close()
	statusCode := DeleteAccount(server.URL, accountID, version)

	msg := fmt.Sprintf("TestDeleteAccount failed. Status code expected to be %d but it was %d", http.StatusNoContent, statusCode)

	if statusCode != http.StatusNoContent {
		t.Errorf(msg)
	}
}

func TestCreateRequestBody_WithAccountIdAndOrganisationId(t *testing.T) {

	var actualAccountID, actualOrganisationID string
	actualAccountID = guuid.New().String()
	actualOrganisationID = guuid.New().String()

	account := CreateRequestBody(actualAccountID, actualOrganisationID)

	if actualAccountID != account.Cdata.ID {
		t.Errorf("Expected %v but got %v", actualAccountID, account.Cdata.ID)
	}
	if actualOrganisationID != account.Cdata.OrganisationID {
		t.Errorf("Expected %v but got %v", actualOrganisationID, account.Cdata.OrganisationID)
	}
}

func TestUnmarshallCreateAccountResponse(t *testing.T) {

	actualAccountID := guuid.New().String()
	actualOrganisationID := guuid.New().String()

	account := CreateRequestBody(actualAccountID, actualOrganisationID)

	jsonBytes, _ := json.Marshal(account)
	body := ioutil.NopCloser(bytes.NewReader(jsonBytes))

	response := &http.Response{
		StatusCode: 201,
		Body:       body,
	}

	accountFromResponse := UnmarshallCreateAccountResponse(response)
	assert.EqualValues(t, actualAccountID, accountFromResponse.Cdata.ID)
	assert.EqualValues(t, actualOrganisationID, accountFromResponse.Cdata.OrganisationID)
}

func TestUnmarshallGetAccountResponse(t *testing.T) {

	var getAccountResponse = GetAccountResponse{Gdata: Gdata{Type: "accounts", ID: "0673746b-8dd3-4bd2-b398-941bdf2865df"}}
	jsonBytes, _ := json.Marshal(getAccountResponse)
	body := ioutil.NopCloser(bytes.NewReader(jsonBytes))

	response := &http.Response{
		StatusCode: 201,
		Body:       body,
	}

	accountFromResponse := UnmarshallGetAccountResponse(response)
	assert.EqualValues(t, "0673746b-8dd3-4bd2-b398-941bdf2865df", accountFromResponse.Gdata.ID)
}
