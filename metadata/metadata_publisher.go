package metadata

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"sync"

	"github.com/gosuri/uilive"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("v1-metadata-publisher")

type PublishService interface {
	Publish() error
	SendMetadataJob(contents []Content, errorsCh chan error, doneCh chan bool)
}

type V1MetadataPublishService struct {
	cs         ContentService
	publishing *Cluster
	mr         ReadService
	source     string
	batchSize  int
	client     *http.Client
}

func NewV1MetadataPublishService(contentService ContentService, publishing *Cluster, mr ReadService, source string, batchSize int) *V1MetadataPublishService {
	return &V1MetadataPublishService{
		cs:         contentService,
		publishing: publishing,
		mr:         mr,
		source:     source,
		batchSize:  batchSize,
		client:     &http.Client{Transport: transport},
	}
}

func (mp *V1MetadataPublishService) Publish() error {
	contentErr := make(chan error)
	defer close(contentErr)
	publishErr := make(chan error)
	defer close(publishErr)
	done := make(chan bool)
	defer close(done)
	writer := uilive.New()
	writer.Start()

	contentCh := mp.cs.GetContent(mp.source, contentErr)
	batch := []Content{}
	counter := 0
	progress := 0

	for {
		select {
		case err := <-contentErr:
			if err != nil {
				return err
			}
		case content, ok := <-contentCh:
			if !ok {
				if len(batch) > 0 {
					go mp.SendMetadataJob(batch, publishErr, done)
					wait(publishErr, done)
				}
				fmt.Fprintf(writer, "\nFinished: %d contents published for source %s\n", progress, mp.source)
				writer.Stop()
				return nil
			}
			if counter < mp.batchSize {
				batch = append(batch, content)
				counter++
			} else {
				progress = progress + len(batch)
				fmt.Fprintf(writer, "%d contents published\n", progress)

				go mp.SendMetadataJob(batch, publishErr, done)
				wait(publishErr, done)
				counter = 0
				batch = []Content{}

				if progress%50000 == 0 {
					time.Sleep(5 * time.Minute)
				}
			}
		}
	}
}
func (mp *V1MetadataPublishService) SendMetadataJob(contents []Content, errorsCh chan error, doneCh chan bool) {
	var wg sync.WaitGroup
	wg.Add(len(contents))
	rate := time.Second / time.Duration(len(contents))
	throttle := time.Tick(rate)

	for i, content := range contents {
		<-throttle
		go func(content Content, i int) {
			defer wg.Done()
			value, err := mp.mr.ReadByUUID(content)
			if err != nil {
				errorsCh <- err
				return
			}
			if len(value) == 0 {
				log.Infof("No metadata found for content=[%s]", content)
				return
			}
			err = mp.publishMetadataForUUID(content, value)
			if err != nil {
				j, _ := json.Marshal(content)
				log.Errorf("Metadata publish for content=[%s] failed because: [%s]", j, err)
				errorsCh <- err
				return
			}
		}(content, i)
	}
	wg.Wait()
	doneCh <- true
}

func (mp *V1MetadataPublishService) publishMetadataForUUID(content Content, metadata []byte) error {
	body, err := getPayload(content.UUID, metadata)
	if err != nil {
		return err
	}

	req, err := getPublishRequest(body, mp.publishing.GetAddress(), mp.publishing.GetUsername(), mp.publishing.GetPassword())
	if err != nil {
		return err
	}

	resp, err := mp.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		j, _ := json.Marshal(content)
		return fmt.Errorf("Publishing of metadata for content=[%s] failed with status code %d", j, resp.StatusCode)
	}

	tid := resp.Header.Get("X-Request-Id")
	log.Infof("Metadata published for content=[%s] having tid=[%s]", content.UUID, tid)
	return nil
}

func getPayload(UUID string, metadata []byte) ([]byte, error) {
	message := map[string]interface{}{
		"uuid":         UUID,
		"lastModified": time.Now().String(),
		"value":        metadata,
	}
	return json.Marshal(message)
}

func getPublishRequest(body []byte, url string, username string, password string) (*http.Request, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(username, password)
	req.Header.Add("X-Origin-System-Id", "binding-service")
	req.Header.Add("Content-Type", "application/json")
	return req, nil
}

func wait(publishErr chan error, done chan bool) {
	for {
		select {
		case err := <-publishErr:
			checkError(err, "metadata publishing")
		case <-done:
			return
		}
	}
}

func checkError(err error, operation string) bool {
	if err != nil {
		log.Errorf("Error occured while %s: %s", operation, err)
		return true
	}
	return false
}
