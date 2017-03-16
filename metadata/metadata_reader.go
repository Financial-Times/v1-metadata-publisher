package metadata

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/bobziuchkovski/digest"
)

const SourcePlaceholder = "{source}"
const UUIDPlaceholder = "{uuid}"

type MetadataReadService interface {
	ReadByUUID(content Content) ([]byte, error)
}

type V1MetadataReadService struct {
	transport *digest.Transport
	url       string
}

func NewV1MetadataReadService(cmr *Cluster) (*V1MetadataReadService, error) {
	url := cmr.GetAddress()
	if !strings.Contains(url, SourcePlaceholder) || !strings.Contains(url, UUIDPlaceholder) {
		return nil, fmt.Errorf("Metadata URL is invalid")
	}
	return &V1MetadataReadService{
		transport: digest.NewTransport(cmr.GetUsername(), cmr.GetPassword()),
		url:       cmr.GetAddress()}, nil
}

func (c *V1MetadataReadService) ReadByUUID(content Content) ([]byte, error) {
	var metadata []byte
	url, err := c.buildURL(content)
	if err != nil {
		log.Errorf("Error while building metadata URL: %s", err)
		return metadata, err
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return metadata, err
	}
	req.Header.Add("ClientUserPrincipal", "upp")

	resp, err := c.transport.RoundTrip(req)
	if err != nil {
		return metadata, err
	}

	//if status is 204 means that there is no metadata for this piece of content
	if resp.StatusCode == http.StatusNoContent {
		return metadata, nil
	}
	if resp.StatusCode != http.StatusOK {
		return metadata, fmt.Errorf("Received response with status code [%d] from binding service for UUID=[%s]", resp.StatusCode, content.UUID)
	}

	metadata, err = ioutil.ReadAll(resp.Body)
	return metadata, err
}

func (c *V1MetadataReadService) buildURL(content Content) (string, error) {
	source, err := content.getSource()
	if err != nil {
		return "", err
	}

	url := strings.Replace(c.url, SourcePlaceholder, source, -1)
	url = strings.Replace(url, UUIDPlaceholder, content.UUID, -1)
	return url, nil

}
