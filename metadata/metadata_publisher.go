package metadata

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"sync"

	log "github.com/Sirupsen/logrus"
)

const UUIDFile = "content.json"

type MetadataPublishService interface {
	Publish() error
}

type V1MetadataPublishService struct {
	cs         ContentService
	publishing *Cluster
	mr         MetadataReadService
	source     string
	batchSize  int
}

func NewV1MetadataPublishService(contentSerivce ContentService, publishing *Cluster, mr MetadataReadService, source string, batchSize int) *V1MetadataPublishService {
	return &V1MetadataPublishService{
		cs:         contentSerivce,
		publishing: publishing,
		mr:         mr,
		source:     source,
		batchSize:  batchSize,
	}
}

func (mp *V1MetadataPublishService) Publish() error {
	f, total, err := mp.cs.SaveContent(UUIDFile)
	if err != nil {
		return err
	}
	defer f.Close()
	defer os.Remove(UUIDFile)

	errorsCh := make(chan error)
	doneCh := make(chan bool)
	metadataBatch := []Content{}
	counter := 0
	scanner := bufio.NewScanner(f)
	progress := 0.0
	t := float64(total)
	for scanner.Scan() {
		var content Content
		json.Unmarshal(scanner.Bytes(), &content)
		cs, err := content.getSource()
		if err != nil || cs != mp.source {
			log.Warnf("Cannot publish metadata for content with UUID=[%d]. Skipping...", content.UUID)
			continue
		}

		if counter < mp.batchSize {
			metadataBatch = append(metadataBatch, content)
			counter++
		} else {
			log.Infof("Publishing metadata for content %s", metadataBatch)
			progress = progress + float64(len(metadataBatch))
			log.Infof("Progress is %.2f", (progress*100)/t)
			go mp.sendMetadataJob(metadataBatch, errorsCh, doneCh)
			<-doneCh
			counter = 0
			metadataBatch = []Content{}
			time.Sleep(1 * time.Second)
		}
	}

	if len(metadataBatch) > 0 {
		go mp.sendMetadataJob(metadataBatch, errorsCh, doneCh)
		<-doneCh
	}
	return nil
}

func (mp *V1MetadataPublishService) sendMetadataJob(contents []Content, errorsCh chan error, doneCh chan bool) {
	var wg sync.WaitGroup
	wg.Add(len(contents))
	rate := time.Second / time.Duration(len(contents))
	throttle := time.Tick(rate)

	for _, content := range contents {
		<-throttle
		go func(content Content) {
			defer wg.Done()
			metadata, err := mp.mr.ReadByUUID(content)
			if err != nil {
				errorsCh <- err
			}
			if len(metadata) == 0 {
				log.Infof("No metadata found for content with UUID=[%s]", content.UUID)
				return
			}
			err = mp.publishMetadataForUUID(content.UUID, metadata)
			if err != nil {
				errorsCh <- err
			}
		}(content)
	}
	go handleErrors(errorsCh)
	wg.Wait()
	doneCh <- true
}

func (mp *V1MetadataPublishService) publishMetadataForUUID(UUID string, metadata []byte) error {
	client := &http.Client{}
	body, err := getPayload(UUID, metadata)
	if err != nil {
		return err
	}

	req, err := getPublishRequest(body, mp.publishing.GetAddress(), mp.publishing.GetUsername(), mp.publishing.GetPassword())
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Publishing of content metadata with UUID=[%s] failed with status code %d", UUID, resp.StatusCode)
	}

	tid := resp.Header.Get("X-Request-Id")
	log.Infof("Metadata published for content with UUID=[%s] having tid=[%s]", UUID, tid)
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

func handleErrors(errors chan error) {
	for e := range errors {
		checkError(e, "metadata publishing")
	}
}

func checkError(err error, operation string) bool {
	if err != nil {
		log.Errorf("Error occured while %s: %s", operation, err)
		return true
	}
	return false
}
