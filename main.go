package main

import (
	"os"

	"github.com/Financial-Times/v1-metadata-publisher/metadata"
	"github.com/op/go-logging"
	"github.com/jawher/mow.cli"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

var log = logging.MustGetLogger("v1-metadata-publisher.log")
var format = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} %{shortfunc}: %{level:.4s} %{message}`,
)

func main() {
	app := cli.App("v1-metadata-publisher", "Application for publishing V1 metadata")

	deliveryCluster := app.String(cli.StringOpt{
		Name:   "deliveryCluster",
		Value:  "http://localhost:8080",
		Desc:   "Address of the delivery cluster",
		EnvVar: "DELIVERY_CLUSTER",
	})

	publishingCluster := app.String(cli.StringOpt{
		Name:   "publishingCluster",
		Value:  "http://localhost:8080",
		Desc:   "Address of the publishing cluster cluster",
		EnvVar: "PUBLISHING_CLUSTER",
	})

	publishingClusterCredentials := app.String(cli.StringOpt{
		Name:   "publishingClusterCredentials",
		Desc:   "Address of the publishing cluster cluster",
		EnvVar: "PUBLISHING_CLUSTER_CREDENTIALS",
	})

	cmrAddress := app.String(cli.StringOpt{
		Name:   "cmrAddress",
		Value:  "http://localhost:8080",
		Desc:   "Address of the Central Metadata Repository",
		EnvVar: "CMR_ADDRESS",
	})

	cmrCredentials := app.String(cli.StringOpt{
		Name:   "cmrCredentials",
		Desc:   "Credentials for Central Metadata Repository",
		EnvVar: "CMR_CREDENTIALS",
	})

	source := app.String(cli.StringOpt{
		Name:   "source",
		Value:  "METHODE",
		Desc:   "Souce of the content",
		EnvVar: "SOURCE",
	})

	batchSize := app.Int(cli.IntOpt{
		Name:   "batchSize",
		Value:  10,
		Desc:   "Number of requests to be sent at a time",
		EnvVar: "BATCH_SIZE",
	})

	initLogging()

	app.Action = func() {
		delivery := metadata.GetCluster(*deliveryCluster, "")
		publishing := metadata.GetCluster(*publishingCluster, *publishingClusterCredentials)
		cmr := metadata.GetCluster(*cmrAddress, *cmrCredentials)

		cmrReader, err := metadata.NewV1MetadataReadService(cmr)
		if err != nil {
			log.Errorf("Cannot start application: %s", err)
			return
		}

		contentService, err := metadata.InitContentService(delivery)
		if err != nil {
			log.Errorf("Cannot start application: %s", err)
			return
		}
		mp := metadata.NewV1MetadataPublishService(contentService, publishing, cmrReader, *source, *batchSize)
		mp.Publish()

		httpHandler := metadata.NewHttpHandler(mp)
		listen(httpHandler, 8080)
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Errorf("Cannot start application: %s", err)
	}
}

func initLogging() {
	errorLog, err := os.OpenFile("v1-metadata-publisher-error.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	infoLog, err := os.OpenFile("v1-metadata-publisher.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	fileBackend := logging.NewLogBackend(errorLog, "", 0)
	stdBackend := logging.NewLogBackend(infoLog, "", 0)

	fileBackendFormatter := logging.NewBackendFormatter(fileBackend, format)
	stdBackendFormatter := logging.NewBackendFormatter(stdBackend, format)

	fileBackendLeveled := logging.AddModuleLevel(fileBackendFormatter)
	fileBackendLeveled.SetLevel(logging.ERROR, "")

	logging.SetBackend(fileBackendLeveled, stdBackendFormatter)
}

func listen(h *metadata.HttpHandler, port int) {
	r := mux.NewRouter()
	r.HandleFunc("/metadata/publish", h.Publish).Methods("POST")

	err := http.ListenAndServe(":"+strconv.Itoa(port), r)
	if err != nil {
		log.Error(err)
	}
}
