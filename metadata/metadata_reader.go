package metadata

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/bobziuchkovski/digest"
	"github.com/pkg/errors"
	"encoding/json"
)

const (
	SourcePlaceholder = "{source}"
	UUIDPlaceholder   = "{uuid}"
)

type ReadService interface {
	ReadByUUID(content Content) ([]byte, error)
}

type V1MetadataReadService struct {
	client *http.Client
	url    string
}

func NewV1MetadataReadService(cmr *Cluster) (*V1MetadataReadService, error) {
	url := cmr.GetAddress()
	if !strings.Contains(url, SourcePlaceholder) || !strings.Contains(url, UUIDPlaceholder) {
		return nil, errors.New("Metadata URL is invalid")
	}

	t := digest.NewTransport(cmr.GetUsername(), cmr.GetPassword())
	t.Transport = transport
	c, err := t.Client()
	if err != nil {
		return nil, err
	}
	return &V1MetadataReadService{
		client: c,
		url:    cmr.GetAddress()}, nil
}

func (c *V1MetadataReadService) ReadByUUID(content Content) ([]byte, error) {
	var result []byte
	url, err := c.buildURL(content)
	if err != nil {
		log.Errorf("Error while building metadata URL: %s", err)
		return result, err
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return result, err
	}
	req.Header.Add("ClientUserPrincipal", "upp")

	resp, err := c.client.Do(req)
	if err != nil {
		j, _ := json.Marshal(content)
		log.Errorf("Getting metadata for content=[%s] failed: %s", j, err)
		return result, err
	}
	defer resp.Body.Close()

	//if status is 204 means that there is no metadata for this piece of content
	if resp.StatusCode == http.StatusNoContent {
		return result, nil
	}
	if resp.StatusCode != http.StatusOK {
		j, _ := json.Marshal(content)
		return result, fmt.Errorf("Received response with status code %d from binding service for content=[%s]", resp.StatusCode, j)
	}
	result, err = ioutil.ReadAll(resp.Body)
	return result, err
}

func (c *V1MetadataReadService) buildURL(content Content) (string, error) {
	source, ok := content.getSource()
	if !ok {
		j, _ := json.Marshal(content)
		return "", fmt.Errorf("Cannot get source of content=[%s]", j)
	}

	url := strings.Replace(c.url, SourcePlaceholder, source, -1)
	url = strings.Replace(url, UUIDPlaceholder, content.UUID, -1)
	return url, nil

}
