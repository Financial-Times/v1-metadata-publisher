package metadata

import (
	"testing"

	"net/http"
	"net/http/httptest"

	"io/ioutil"

	"github.com/stretchr/testify/assert"
)

const BindingServiceURL = "/metadata-services/binding/1.0/sources/{source}/references/{uuid}"

var testContent = Content{
	UUID:        "0cd42702-f789-11e6-9516-2d969e0d3b65",
	Identifiers: []Identifier{{Authority: "http://api.ft.com/system/FTCOM-METHODE"}},
}

func TestReadByUUIDSuccessfull(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		metadata, err := getMetadata()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.Write(metadata)

	}))
	defer ts.Close()

	expectedResponse, err := getMetadata()
	assert.NoError(t, err, "Test failed because resources could not be read")

	cmr := Cluster{
		address:  ts.URL + BindingServiceURL,
		username: "foo",
		password: "bar",
	}
	reader, err := NewV1MetadataReadService(&cmr)
	assert.NoError(t, err, "Failed to initialise metadata reader")

	result, err := reader.ReadByUUID(testContent)
	assert.NoError(t, err, "Failed to read metadata")
	assert.Equal(t, expectedResponse, result, "Actual metadata differs from expected metadata")

}

func TestReadByUUIDEmptyResponseWhenNoMetadataIsReturned(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte{})
		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	expectedResponse := []byte{}

	cmr := Cluster{
		address:  ts.URL + BindingServiceURL,
		username: "foo",
		password: "bar",
	}
	reader, err := NewV1MetadataReadService(&cmr)
	assert.NoError(t, err, "Failed to initialise metadata reader")
	result, err := reader.ReadByUUID(testContent)
	assert.NoError(t, err, "Failed to read metadata")
	assert.Equal(t, expectedResponse, result, "Actual metadata differs from expected metadata")
}

func TestReadByUUIDUnsuccessfull(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	cmr := Cluster{
		address:  ts.URL + BindingServiceURL,
		username: "foo",
		password: "bar",
	}
	reader, err := NewV1MetadataReadService(&cmr)
	assert.NoError(t, err, "Failed to initialise metadata reader")
	_, err = reader.ReadByUUID(testContent)
	assert.Error(t, err, "Getting metadata should return error")
}

func TestBuildURL(t *testing.T) {
	expectedURL := "http://localhost:8080/metadata-services/binding/1.0/sources/METHODE/references/0cd42702-f789-11e6-9516-2d969e0d3b65"
	cmr := V1MetadataReadService{
		url: "http://localhost:8080" + BindingServiceURL,
	}
	actual, err := cmr.buildURL(testContent)
	assert.NoError(t, err, "msgAndArgs ...interface{}")
	assert.Equal(t, expectedURL, actual, "Result URL is different from expected URL")
}

func getMetadata() ([]byte, error) {
	return ioutil.ReadFile("resources/metadata-response.xml")
}
