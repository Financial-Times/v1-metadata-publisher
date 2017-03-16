package metadata

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"os"

	"github.com/stretchr/testify/assert"
)

const TestFile = "resources/content-test.json"

func TestSaveUUIDsSuccesfully(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(UUIDResponse)
	}))
	defer ts.Close()

	cs := UPPContentService{
		delivery: &Cluster{
			address:  ts.URL + "/__document-store-api/content/__ids",
			username: "foo",
			password: "bar",
		},
	}

	f, total, err := cs.SaveContent(TestFile)
	assert.NoError(t, err, "Failed to save content UUIDs")
	assert.NoError(t, f.Close(), "Error in closing test file")
	assert.Equal(t, 1, total, "Actual number of items if different from actual number of items")
	defer os.Remove(TestFile)

	result, err := ioutil.ReadFile(TestFile)
	assert.NoError(t, err, "Failed to save content UUIDs")
	assert.Equal(t, UUIDResponse, result, "Expected UUIDs are different from actual UUIDs")
}

func TestSaveUUIDsGettingUUIDsFails(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	cs := UPPContentService{
		delivery: &Cluster{
			address:  ts.URL + "/__document-store-api/content/__ids",
			username: "username",
			password: "password",
		},
	}

	f, _, err := cs.SaveContent(TestFile)
	defer f.Close()
	defer os.Remove(TestFile)

	assert.Error(t, err, "Expecting error while trying to get the UUIDs")
}
