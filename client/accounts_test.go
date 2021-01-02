package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
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
	// prepare
	accountID := guuid.New().String()
	organizationID := guuid.New().String()
	account := CreateRequestBody(accountID, organizationID)

	uri := "/v1/organisation/accounts/"

	var body = GetAccountResponse{Gdata: Gdata{Type: "accounts", ID: "0673746b-8dd3-4bd2-b398-941bdf2865df", OrganisationID: "9864746b-8dd3-4bd2-b398-941bdf2865df"}}

	server := newTestServer(uri, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(body)
	})

	defer server.Close()
	response, _ := CreateAccount(server.URL, account)

	// test & validate
	createdAccount, err := UnmarshallCreateAccountResponse(response)

	msg := fmt.Sprintf("TestCreateAccount failed. Status code expected to be %d but it was %d", http.StatusOK, response.StatusCode)

	if response.StatusCode != 201 {
		t.Errorf(msg)
	}
	assert.Nil(t, err)
	assert.EqualValues(t, createdAccount.Cdata.ID, "0673746b-8dd3-4bd2-b398-941bdf2865df")
	assert.EqualValues(t, createdAccount.Cdata.OrganisationID, "9864746b-8dd3-4bd2-b398-941bdf2865df")
}

func TestCreateAccount_whenForm3ApiReturns500_shouldReturn500(t *testing.T) {
	// prepare
	accountID := guuid.New().String()
	organizationID := guuid.New().String()
	account := CreateRequestBody(accountID, organizationID)

	uri := "/v1/organisation/accounts/"

	server := newTestServer(uri, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})

	defer server.Close()

	// test & validate
	response, err := CreateAccount(server.URL, account)

	msg := fmt.Sprintf("TestCreateAccount failed. Status code expected to be %d but it was %d", http.StatusInternalServerError, response.StatusCode)

	if response.StatusCode != 500 {
		t.Errorf(msg)
	}

	assert.Nil(t, err)
	assert.NotNil(t, response)
}

func TestCreateAccount_whenMarshallerFails_shouldReturnError(t *testing.T) {
	// prepare
	accountID := guuid.New().String()
	organizationID := guuid.New().String()
	account := CreateRequestBody(accountID, organizationID)

	uri := "/v1/organisation/accounts/"

	var body = GetAccountResponse{Gdata: Gdata{Type: "accounts", ID: "0673746b-8dd3-4bd2-b398-941bdf2865df", OrganisationID: "9864746b-8dd3-4bd2-b398-941bdf2865df"}}

	server := newTestServer(uri, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(body)
	})
	defer server.Close()

	Marshaller = func(v interface{}) ([]byte, error) {
		return nil, errors.New("Marshaller faillure")
	}

	// test & validate
	response, err := CreateAccount(server.URL, account)

	assert.Nil(t, response)
	assert.NotNil(t, err)
	assert.EqualValues(t, fmt.Sprint(err), "Marshaller faillure")
}

func TestCreateAccount_whenRequestCreatorFails_shouldReturnError(t *testing.T) {
	// prepare
	accountID := guuid.New().String()
	organizationID := guuid.New().String()
	account := CreateRequestBody(accountID, organizationID)

	uri := "/v1/organisation/accounts/"

	var body = GetAccountResponse{Gdata: Gdata{Type: "accounts", ID: "0673746b-8dd3-4bd2-b398-941bdf2865df", OrganisationID: "9864746b-8dd3-4bd2-b398-941bdf2865df"}}

	server := newTestServer(uri, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(body)
	})
	defer server.Close()

	Marshaller = json.Marshal

	RequestCreator = func(method, url string, body io.Reader) (*http.Request, error) {
		return nil, errors.New("RequestCreator faillure")
	}

	// test & validate
	response, err := CreateAccount(server.URL, account)

	assert.Nil(t, response)
	assert.NotNil(t, err)
	assert.EqualValues(t, fmt.Sprint(err), "RequestCreator faillure")
}

func TestGetAccount_success(t *testing.T) {
	// prepare
	accountID := guuid.New().String()
	uri := "/v1/organisation/accounts/"

	var body = Account{Cdata: Cdata{Type: "accounts", ID: "0673746b-8dd3-4bd2-b398-941bdf2865df"}}

	server := newTestServer(uri, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(body)
	})
	defer server.Close()

	RequestCreator = http.NewRequest

	// test & validate
	response, error := GetAccount(server.URL, accountID)

	assert.Nil(t, error)
	assert.EqualValues(t, response.StatusCode, 200)
}

func TestGetAccount_whenForm3ApiReturnes500_shouldReturn500(t *testing.T) {
	// prepare
	accountID := guuid.New().String()
	uri := "/v1/organisation/accounts/"

	server := newTestServer(uri, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	defer server.Close()

	RequestCreator = http.NewRequest
	// test & validate
	getAccountResponse, err := GetAccount(server.URL, accountID)

	assert.Nil(t, err)
	assert.EqualValues(t, getAccountResponse.StatusCode, 500)
}

func TestGetAccount_whenRequestCreatorFails_shouldReturnError(t *testing.T) {
	// prepare
	accountID := guuid.New().String()
	uri := "/v1/organisation/accounts/"

	server := newTestServer(uri, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	defer server.Close()

	RequestCreator = func(method, url string, body io.Reader) (*http.Request, error) {
		return nil, errors.New("RequestCreator faillure")
	}
	// test & validate
	getAccountResponse, err := GetAccount(server.URL, accountID)

	assert.Nil(t, getAccountResponse)
	assert.NotNil(t, err)
	assert.EqualValues(t, fmt.Sprint(err), "RequestCreator faillure")
}

func TestListAccounts_success(t *testing.T) {
	// prepare
	pageNumber := 1
	pageSize := 30
	uri := "/v1/organisation/accounts"

	var body = GetAccountsResponse{Data: []Data{{Type: "accounts", ID: "0673746b-8dd3-4bd2-b398-941bdf2865df"},
		{Type: "accounts", ID: "9673746b-8dd3-4bd2-b398-941bdf2865df"}}}

	server := newTestServer(uri, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(body)
	})
	defer server.Close()

	RequestCreator = http.NewRequest

	// test
	getAccountsResponse, err := ListAccounts(server.URL, pageNumber, pageSize)

	// validate
	msg := fmt.Sprintf("TestListAccounts failed. Status code expected to be %d but it was %d", http.StatusOK, getAccountsResponse.StatusCode)

	if getAccountsResponse.StatusCode != http.StatusOK {
		t.Errorf(msg)
	}

	accounts := UnmarshallGetAccountsResponse(getAccountsResponse)

	assert.Nil(t, err)
	assert.EqualValues(t, len(accounts.Data), 2)
	assert.EqualValues(t, accounts.Data[0].ID, "0673746b-8dd3-4bd2-b398-941bdf2865df")
	assert.EqualValues(t, accounts.Data[1].ID, "9673746b-8dd3-4bd2-b398-941bdf2865df")
}

func TestListAccounts_whenForm3Returns500_shouldReturn500(t *testing.T) {
	// prepare
	pageNumber := 1
	pageSize := 30
	uri := "/v1/organisation/accounts"

	server := newTestServer(uri, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	defer server.Close()

	RequestCreator = http.NewRequest

	// test & validate
	getAccountsResponse, _ := ListAccounts(server.URL, pageNumber, pageSize)

	msg := fmt.Sprintf("TestListAccounts failed. Status code expected to be %d but it was %d", http.StatusInternalServerError, getAccountsResponse.StatusCode)

	if getAccountsResponse.StatusCode != http.StatusInternalServerError {
		t.Errorf(msg)
	}
}

func TestListAccounts_whenRequestCreatorFails_shouldReturnError(t *testing.T) {
	// prepare
	pageNumber := 1
	pageSize := 30
	uri := "/v1/organisation/accounts"

	server := newTestServer(uri, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	defer server.Close()

	RequestCreator = func(method, url string, body io.Reader) (*http.Request, error) {
		return nil, errors.New("RequestCreator faillure")
	}

	// test & validate
	response, err := ListAccounts(server.URL, pageNumber, pageSize)

	assert.Nil(t, response)
	assert.NotNil(t, err)
	assert.EqualValues(t, fmt.Sprint(err), "RequestCreator faillure")
}

func TestGatherAccounts_WhenListAccountsReturns500(t *testing.T) {
	// prepare
	pageSize := 30
	uri := "/v1/organisation/accounts"

	server := newTestServer(uri, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	defer server.Close()

	RequestCreator = http.NewRequest

	// test & validate
	accounts := GatherAccounts(server.URL, pageSize)

	if len(accounts) > 0 {
		t.Errorf("Expected found accounts to be %v but got %v", 0, len(accounts))
	}
}

func TestGatherAccounts_WhenListAccountsReturns200AndTwoResults(t *testing.T) {
	// prepare
	pageSize := 30
	uri := "/v1/organisation/accounts"

	var body = GetAccountsResponse{Data: []Data{{Type: "accounts", ID: "0673746b-8dd3-4bd2-b398-941bdf2865df"},
		{Type: "accounts", ID: "0673746b-8dd3-4bd2-b398-941bdf2865df"}}}

	server := newTestServer(uri, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(body)
	})
	defer server.Close()

	RequestCreator = http.NewRequest

	// test & validate
	accounts := GatherAccounts(server.URL, pageSize)

	if len(accounts) != 2 {
		t.Errorf("Expected found accounts to be %v but got %v", 2, len(accounts))
	}
}

func TestDeleteAccount_success(t *testing.T) {
	// prepare
	accountID := guuid.New().String()
	version := 0
	uri := "/v1/organisation/accounts/"

	server := newTestServer(uri, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
	defer server.Close()

	RequestCreator = http.NewRequest

	// test & validate
	deleteAccountResponse, _ := DeleteAccount(server.URL, accountID, version)

	msg := fmt.Sprintf("TestDeleteAccount failed. Status code expected to be %d but it was %d", http.StatusNoContent, deleteAccountResponse.StatusCode)

	if deleteAccountResponse.StatusCode != http.StatusNoContent {
		t.Errorf(msg)
	}
}

func TestDeleteAccount_whenForm3ApiReturns404_shouldReturn404(t *testing.T) {
	// prepare
	accountID := guuid.New().String()
	version := 0
	uri := "/v1/organisation/accounts/"

	server := newTestServer(uri, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})
	defer server.Close()

	RequestCreator = http.NewRequest

	// test & validate
	deleteAccountResponse, _ := DeleteAccount(server.URL, accountID, version)

	msg := fmt.Sprintf("TestDeleteAccount failed. Status code expected to be %d but it was %d", http.StatusNotFound, deleteAccountResponse.StatusCode)

	if deleteAccountResponse.StatusCode != http.StatusNotFound {
		t.Errorf(msg)
	}
}

func TestDeleteAccount_whenForm3ApiReturns409_shouldReturn409(t *testing.T) {
	// prepare
	accountID := guuid.New().String()
	version := 0
	uri := "/v1/organisation/accounts/"

	server := newTestServer(uri, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusConflict)
	})
	defer server.Close()

	RequestCreator = http.NewRequest

	// test & validate
	deleteAccountResponse, _ := DeleteAccount(server.URL, accountID, version)

	msg := fmt.Sprintf("TestDeleteAccount failed. Status code expected to be %d but it was %d", http.StatusNotFound, deleteAccountResponse.StatusCode)

	if deleteAccountResponse.StatusCode != http.StatusConflict {
		t.Errorf(msg)
	}
}

func TestDeleteAccount_whenRequestCreatorFails_shouldReturnError(t *testing.T) {
	// prepare
	accountID := guuid.New().String()
	version := 0
	uri := "/v1/organisation/accounts/"

	server := newTestServer(uri, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	defer server.Close()

	RequestCreator = func(method, url string, body io.Reader) (*http.Request, error) {
		return nil, errors.New("RequestCreator faillure")
	}

	// test & validate
	deleteAccountResponse, err := DeleteAccount(server.URL, accountID, version)

	assert.Nil(t, deleteAccountResponse)
	assert.NotNil(t, err)
	assert.EqualValues(t, fmt.Sprint(err), "RequestCreator faillure")
}

func TestCreateRequestBody_WithAccountIdAndOrganisationId(t *testing.T) {
	// prepare
	var actualAccountID, actualOrganisationID string
	actualAccountID = guuid.New().String()
	actualOrganisationID = guuid.New().String()

	// test & validate
	account := CreateRequestBody(actualAccountID, actualOrganisationID)

	if actualAccountID != account.Cdata.ID {
		t.Errorf("Expected %v but got %v", actualAccountID, account.Cdata.ID)
	}
	if actualOrganisationID != account.Cdata.OrganisationID {
		t.Errorf("Expected %v but got %v", actualOrganisationID, account.Cdata.OrganisationID)
	}
}

func TestUnmarshallCreateAccountResponse_success(t *testing.T) {
	// prepare
	actualAccountID := guuid.New().String()
	actualOrganisationID := guuid.New().String()

	account := CreateRequestBody(actualAccountID, actualOrganisationID)

	jsonBytes, _ := json.Marshal(account)
	body := ioutil.NopCloser(bytes.NewReader(jsonBytes))

	response := &http.Response{
		StatusCode: 201,
		Body:       body,
	}

	// test
	accountFromResponse, err := UnmarshallCreateAccountResponse(response)

	// validate
	assert.Nil(t, err)
	assert.EqualValues(t, actualAccountID, accountFromResponse.Cdata.ID)
	assert.EqualValues(t, actualOrganisationID, accountFromResponse.Cdata.OrganisationID)
}

func TestUnmarshallCreateAccountResponse_whenUnmarshallerFails_returnsError(t *testing.T) {
	// prepare
	actualAccountID := guuid.New().String()
	actualOrganisationID := guuid.New().String()

	account := CreateRequestBody(actualAccountID, actualOrganisationID)

	jsonBytes, _ := json.Marshal(account)
	body := ioutil.NopCloser(bytes.NewReader(jsonBytes))

	response := &http.Response{
		StatusCode: 201,
		Body:       body,
	}

	Unmarshaller = func(data []byte, v interface{}) error {
		return errors.New("Unmarshaller faillure")
	}

	// test
	accountFromResponse, err := UnmarshallCreateAccountResponse(response)

	// validate
	assert.Nil(t, accountFromResponse)
	assert.NotNil(t, err)
	assert.EqualValues(t, fmt.Sprint(err), "Unmarshaller faillure")
}

func TestUnmarshallGetAccountResponse_success(t *testing.T) {
	// prepare
	var getAccountResponse = GetAccountResponse{Gdata: Gdata{Type: "accounts", ID: "0673746b-8dd3-4bd2-b398-941bdf2865df"}}
	jsonBytes, _ := json.Marshal(getAccountResponse)
	body := ioutil.NopCloser(bytes.NewReader(jsonBytes))

	response := &http.Response{
		StatusCode: 201,
		Body:       body,
	}

	Unmarshaller = json.Unmarshal

	// test & validate
	accountFromResponse, err := UnmarshallGetAccountResponse(response)
	assert.EqualValues(t, "0673746b-8dd3-4bd2-b398-941bdf2865df", accountFromResponse.Gdata.ID)
	assert.Nil(t, err)
}

func TestUnmarshallGetAccountResponse_whenUnmarshallerFails_shouldReturnError(t *testing.T) {
	// prepare
	var getAccountResponse = GetAccountResponse{Gdata: Gdata{Type: "accounts", ID: "0673746b-8dd3-4bd2-b398-941bdf2865df"}}
	jsonBytes, _ := json.Marshal(getAccountResponse)
	body := ioutil.NopCloser(bytes.NewReader(jsonBytes))

	response := &http.Response{
		StatusCode: 201,
		Body:       body,
	}

	Unmarshaller = func(data []byte, v interface{}) error {
		return errors.New("Unmarshaller faillure")
	}

	// test & validate
	accountFromResponse, err := UnmarshallGetAccountResponse(response)
	assert.Nil(t, accountFromResponse)
	assert.NotNil(t, err)
}
