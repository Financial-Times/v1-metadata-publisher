package metadata

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"

	log "github.com/Sirupsen/logrus"
)

type ContentService interface {
	SaveContent(file string) (*os.File, int, error)
}

type UPPContentService struct {
	delivery *Cluster
}

func InitContentService(delivery *Cluster) *UPPContentService {
	return &UPPContentService{delivery: delivery}
}

func (c *UPPContentService) SaveContent(file string) (*os.File, int, error) {
	f, err := os.Create(file)
	if err != nil {
		return nil, 0, err
	}

	reader, err := c.readContent()
	if err != nil {
		return f, 0, err
	}
	total, err := c.writeContent(f, reader)
	if err != nil {
		return f, 0, err
	}

	f.Seek(0, 0)
	return f, total, nil
}

func (c *UPPContentService) readContent() (*bufio.Reader, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", c.delivery.GetAddress(), nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(c.delivery.GetUsername(), c.delivery.GetPassword())
	q := req.URL.Query()
	q.Add("includeSource", "true")
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Received response with status code %d", resp.StatusCode)
	}

	return bufio.NewReader(resp.Body), nil
}

func (c *UPPContentService) writeContent(f *os.File, r *bufio.Reader) (int, error) {
	counter := 0
	for {
		line, err := r.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				log.Infof("Received %d UUIDs", counter)
				return counter, nil
			}
			return 0, err
		}
		f.Write(line)
		counter++
	}
}
