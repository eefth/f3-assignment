package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	guuid "github.com/google/uuid"
)

// newTestServer creates a multiplex server to handle API endpoints
func newTestServer(path string, h func(w http.ResponseWriter, r *http.Request)) *httptest.Server {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	mux.HandleFunc(path, h)
	return server
}

func TestCreateAccount(t *testing.T) {

	accountID := guuid.New().String()
	organizationID := guuid.New().String()
	uri := "/v1/organisation/accounts/"

	server := newTestServer(uri, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
	})
	defer server.Close()
	statusCode, _ := CreateAccount(server.URL, accountID, organizationID)

	msg := fmt.Sprintf("TestCreateAccount failed. Status code expected to be %d but it was %d", http.StatusOK, statusCode)

	if statusCode != 201 {
		t.Errorf(msg)
	}
}

func TestGetAccount(t *testing.T) {

	accountID := guuid.New().String()
	uri := "/v1/organisation/accounts/"

	server := newTestServer(uri, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	defer server.Close()
	statusCode, _ := GetAccount(server.URL, accountID)

	msg := fmt.Sprintf("TestGetAccount failed. Status code expected to be %d but it was %d", http.StatusOK, statusCode)

	if statusCode != 200 {
		t.Errorf(msg)
	}
}

func TestListAccounts(t *testing.T) {

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

func TestDeleteAccount(t *testing.T) {

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
