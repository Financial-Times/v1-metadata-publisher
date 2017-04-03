package metadata

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/pkg/errors"
)

type MockMetadataReadService struct {
	mockReadByUUID func(content Content) ([]byte, error)
}

func (mr *MockMetadataReadService) ReadByUUID(content Content) ([]byte, error) {
	return mr.mockReadByUUID(content)
}

type MockContentService struct {
	mockGetContent func(source string, errCh chan error) chan Content
}

func (cs *MockContentService) GetContent(source string, errCh chan error) chan Content {
	return cs.mockGetContent(source, errCh)

}

func TestPublishMetadataForUUIDSuccessfully(t *testing.T) {
	ps := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "binding-service", r.Header.Get("X-Origin-System-Id"), "Invalid X-Origin-System-Id header value")
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"), "Invalid Content-Type header value")
		w.Header().Add("X-Request-Id", "tid_testtid")
	}))
	defer ps.Close()

	mps := V1MetadataPublishService{
		publishing: &Cluster{
			address:  ps.URL + "/__cms-metadata-notifier/notify",
			username: "foo",
			password: "bar",
		},
		client: http.DefaultClient,
	}

	cm, err := getMetadata()
	assert.NoError(t, err, "Failed to read metadata")
	err = mps.publishMetadataForUUID(testContent, cm)
	assert.NoError(t, err, "Failed to publish metadata")
}

func TestPublishMetadataForUUIDUnsuccesfully(t *testing.T) {
	ps := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ps.Close()

	mps := V1MetadataPublishService{
		publishing: &Cluster{
			address:  ps.URL + "/__cms-metadata-notifier/notify",
			username: "foo",
			password: "bar",
		},
		client: http.DefaultClient,
	}

	mc, err := getMetadata()
	assert.NoError(t, err, "Failed to read metadata")
	err = mps.publishMetadataForUUID(testContent, mc)
	assert.Error(t, err, "Expected metadata publish will return an error")
}

func TestSendMetadataJobSuccessfully(t *testing.T) {
	publishesDone := 0
	ps := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "binding-service", r.Header.Get("X-Origin-System-Id"), "Invalid X-Origin-System-Id header value")
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"), "Invalid Content-Type header value")
		publishesDone++
		w.Header().Add("X-Request-Id", "tid_testtid")
	}))
	defer ps.Close()

	mps := V1MetadataPublishService{
		publishing: &Cluster{
			address:  ps.URL + "/__cms-metadata-notifier/notify",
			username: "foo",
			password: "bar",
		},
		mr: &MockMetadataReadService{
			mockReadByUUID: func(content Content) ([]byte, error) {
				return getMetadata()
			},
		},
		client: http.DefaultClient,
	}

	errorsCh := make(chan error)
	doneCh := make(chan bool)
	contents := []Content{
		{
			UUID:        "0cd42702-f789-11e6-9516-2d969e0d3b65",
			Identifiers: []Identifier{{Authority: "http://api.ft.com/system/FTCOM-METHODE"}},
		},
	}

	go mps.SendMetadataJob(contents, errorsCh, doneCh)
	go func(errorCh chan error) {
		for err := range errorsCh {
			assert.NoError(t, err, "Error occured while publising metadata")
		}
	}(errorsCh)
	<-doneCh
	assert.Equal(t, 1, publishesDone, "Actual number of publisher requested is different from expected value")
}

func TestSendMetadataJobMetadataServiceNotAvailable(t *testing.T) {
	ps := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "binding-service", r.Header.Get("X-Origin-System-Id"), "Invalid X-Origin-System-Id header value")
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"), "Invalid Content-Type header value")
		w.Header().Add("X-Request-Id", "tid_testtid")
	}))
	defer ps.Close()

	mps := V1MetadataPublishService{
		publishing: &Cluster{
			address:  ps.URL + "/__cms-metadata-notifier/notify",
			username: "foo",
			password: "bar",
		},
		mr: &MockMetadataReadService{
			mockReadByUUID: func(content Content) ([]byte, error) {
				return nil, fmt.Errorf("Cannot get metadata")
			},
		},
		client: http.DefaultClient,
	}

	errorsCh := make(chan error)
	doneCh := make(chan bool)
	contents := []Content{
		{
			UUID:        "0cd42702-f789-11e6-9516-2d969e0d3b65",
			Identifiers: []Identifier{{Authority: "http://api.ft.com/system/FTCOM-METHODE"}},
		},
	}

	go mps.SendMetadataJob(contents, errorsCh, doneCh)
	go func(errorCh chan error) {
		for err := range errorsCh {
			assert.Error(t, err, "Expecting error occured while publising metadata")
		}
	}(errorsCh)
	<-doneCh
}

func TestSendMetadataJobPublishingClusterNotAvailable(t *testing.T) {
	ps := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ps.Close()

	mps := V1MetadataPublishService{
		publishing: &Cluster{
			address:  ps.URL + "/__cms-metadata-notifier/notify",
			username: "foo",
			password: "bar",
		},
		mr: &MockMetadataReadService{
			mockReadByUUID: func(content Content) ([]byte, error) {
				return nil, fmt.Errorf("Cannot get metadata")
			},
		},
		client: http.DefaultClient,
	}

	errorsCh := make(chan error)
	doneCh := make(chan bool)
	contents := []Content{
		{
			UUID:        "0cd42702-f789-11e6-9516-2d969e0d3b65",
			Identifiers: []Identifier{{Authority: "http://api.ft.com/system/FTCOM-METHODE"}},
		},
	}

	go mps.SendMetadataJob(contents, errorsCh, doneCh)
	go func(errorCh chan error) {
		for err := range errorsCh {
			assert.Error(t, err, "Expecting error occured while publising metadata")
		}
	}(errorsCh)
	<-doneCh
}

func TestPublishSuccessful(t *testing.T) {
	ps := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "binding-service", r.Header.Get("X-Origin-System-Id"), "Invalid X-Origin-System-Id header value")
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"), "Invalid Content-Type header value")
		w.Header().Add("X-Request-Id", "tid_testtid")
	}))
	defer ps.Close()

	mps := V1MetadataPublishService{
		cs: &MockContentService{
			mockGetContent: func(source string, errCh chan error) chan Content {
				contentCh := make(chan Content)
				go func() {
					defer close(contentCh)
					contentCh <- testContent
				}()
				return contentCh
			},
		},
		publishing: &Cluster{
			address:  ps.URL + "/__cms-metadata-notifier/notify",
			username: "foo",
			password: "bar",
		},
		mr: &MockMetadataReadService{
			mockReadByUUID: func(content Content) ([]byte, error) {
				m, err := getMetadata()
				assert.NoError(t, err, "Error occured while getting metadata")
				return m, nil
			},
		},
		batchSize: 10,
		source:    "METHODE",
		client:    http.DefaultClient,
	}

	err := mps.Publish()
	assert.NoError(t, err, "Error while trying to publish metadata")
}

func TestPublishUnsuccessfulDeliveryClusterNotAvailble(t *testing.T) {
	ps := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "binding-service", r.Header.Get("X-Origin-System-Id"), "Invalid X-Origin-System-Id header value")
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"), "Invalid Content-Type header value")
		w.Header().Add("X-Request-Id", "tid_testtid")
	}))
	defer ps.Close()

	mps := V1MetadataPublishService{
		cs: &MockContentService{
			mockGetContent: func(source string, errCh chan error) chan Content {
				go func() {
					errCh <- errors.New("Error getting content")
				}()
				return nil
			},
		},
		batchSize: 10,
		source:    "METHODE",
		client:    http.DefaultClient,
	}

	err := mps.Publish()
	assert.Error(t, err, "Expecting error while trying to publish metadata")
}
