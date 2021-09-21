package httpclient

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	// "log"
)

var (
	reqBody = `{
		"test": "test",
	}`
	dummyHandler = func(w http.ResponseWriter, r *http.Request) {
		var bodyData struct {
			Sleep int `json:"sleep"`
		}
		switch r.Method {
		case "GET":
			w.WriteHeader(http.StatusOK)
			out, err := w.Write([]byte("Hello"))
			if err != nil {
				http.Error(w, "can't write response ("+strconv.Itoa(out)+") "+err.Error(), http.StatusBadRequest)
				return
			}
		case "POST":
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "can't read body", http.StatusBadRequest)
				return
			}

			if err = json.Unmarshal(body, &bodyData); err != nil {
				http.Error(w, "Error Unmarshal", http.StatusBadRequest)
				return
			}

			w.WriteHeader(http.StatusOK)
			out, err := w.Write(body)
			if err != nil {
				http.Error(w, "can't write response ("+strconv.Itoa(out)+") "+err.Error(), http.StatusBadRequest)
				return
			}
		case "PUT":
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "can't read body", http.StatusBadRequest)
				return
			}

			if err = json.Unmarshal(body, &bodyData); err != nil {
				http.Error(w, "Error Unmarshal", http.StatusBadRequest)
				return
			}

			w.WriteHeader(http.StatusOK)
			out, err := w.Write(body)
			if err != nil {
				http.Error(w, "can't write response ("+strconv.Itoa(out)+") "+err.Error(), http.StatusBadRequest)
				return
			}
		case "PATCH":
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "can't read body", http.StatusBadRequest)
				return
			}

			if err = json.Unmarshal(body, &bodyData); err != nil {
				http.Error(w, "error unmarshal", http.StatusBadRequest)
				return
			}
			w.WriteHeader(http.StatusOK)
			out, err := w.Write(body)
			if err != nil {
				http.Error(w, "can't write response ("+strconv.Itoa(out)+") "+err.Error(), http.StatusBadRequest)
				return
			}
		case "DELETE":
			w.WriteHeader(http.StatusOK)
			out, err := w.Write([]byte("Deleted"))
			if err != nil {
				http.Error(w, "can't write response ("+strconv.Itoa(out)+") "+err.Error(), http.StatusBadRequest)
				return
			}
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
	}
)

func TestHttpDoGet(t *testing.T) {
	var httpDoer HTTPDoer
	var httpParam HTTPParam
	_ = json.Unmarshal([]byte(reqBody), &httpParam)

	server := httptest.NewServer(http.HandlerFunc(dummyHandler))
	defer server.Close()

	httpParam.URL = server.URL
	httpParam.Method = "get"
	httpDoer = &httpParam
	response, err := httpDoer.HTTPDo()
	require.NoError(t, err, "should not have failed to make a Get request")
	body, _ := ioutil.ReadAll(response.Body)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, "Hello", string(body))
}

func TestHttpDoPost(t *testing.T) {
	var httpDoer HTTPDoer
	var httpParam HTTPParam
	_ = json.Unmarshal([]byte(reqBody), &httpParam)

	server := httptest.NewServer(http.HandlerFunc(dummyHandler))
	defer server.Close()

	httpParam.URL = server.URL
	httpParam.Method = "post"
	httpDoer = &httpParam
	response, err := httpDoer.HTTPDo()
	require.NoError(t, err, "should not have failed to make a POST request")
	body, _ := ioutil.ReadAll(response.Body)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, httpParam.Body, string(body))
}

func TestHttpDoPut(t *testing.T) {
	var httpDoer HTTPDoer
	var httpParam HTTPParam
	_ = json.Unmarshal([]byte(reqBody), &httpParam)

	server := httptest.NewServer(http.HandlerFunc(dummyHandler))
	defer server.Close()

	httpParam.URL = server.URL
	httpParam.Method = "put"
	httpDoer = &httpParam
	response, err := httpDoer.HTTPDo()
	require.NoError(t, err, "should not have failed to make a POST request")
	body, _ := ioutil.ReadAll(response.Body)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, httpParam.Body, string(body))
}

func TestHttpDoPatch(t *testing.T) {
	var httpDoer HTTPDoer
	var httpParam HTTPParam
	_ = json.Unmarshal([]byte(reqBody), &httpParam)

	server := httptest.NewServer(http.HandlerFunc(dummyHandler))
	defer server.Close()

	httpParam.URL = server.URL
	httpParam.Method = "patch"
	httpDoer = &httpParam
	response, err := httpDoer.HTTPDo()
	require.NoError(t, err, "should not have failed to make a POST request")
	body, _ := ioutil.ReadAll(response.Body)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, httpParam.Body, string(body))
}

func TestHttpDoDelete(t *testing.T) {
	var httpDoer HTTPDoer
	var httpParam HTTPParam
	_ = json.Unmarshal([]byte(reqBody), &httpParam)

	server := httptest.NewServer(http.HandlerFunc(dummyHandler))
	defer server.Close()

	httpParam.URL = server.URL
	httpParam.Method = "delete"
	httpDoer = &httpParam
	response, err := httpDoer.HTTPDo()
	require.NoError(t, err, "should not have failed to make a Get request")
	body, _ := ioutil.ReadAll(response.Body)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, "Deleted", string(body))
}

func TestHttpDoDefault(t *testing.T) {
	var httpDoer HTTPDoer
	var httpParam HTTPParam
	_ = json.Unmarshal([]byte(reqBody), &httpParam)

	server := httptest.NewServer(http.HandlerFunc(dummyHandler))
	defer server.Close()

	httpParam.URL = server.URL
	httpParam.Method = ""
	httpDoer = &httpParam
	response, err := httpDoer.HTTPDo()
	require.NoError(t, err, "should not have failed to make a POST request")
	body, _ := ioutil.ReadAll(response.Body)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, httpParam.Body, string(body))
}

func TestMakeHeader(t *testing.T) {
	var httpParam HTTPParam

	httpParam.Header = make(map[string]string)
	httpParam.Header["Content-Type"] = "application/json"

	header := makeHeader(httpParam.Header)
	assert.Equal(t, "application/json", header.Get("Content-Type"))
}
