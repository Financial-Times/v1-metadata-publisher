package main

import (
	"os"

	"github.com/Financial-Times/v1-metadata-publisher/metadata"
	log "github.com/Sirupsen/logrus"
	"github.com/jawher/mow.cli"
)

func main() {
	app := cli.App("v1-metadata-publisher", "Application for publishing V1 metadata")

	deliveryCluster := app.String(cli.StringOpt{
		Name:   "deliveryCluster",
		Value:  "http://localhost:8080",
		Desc:   "Address of the delivery cluster",
		EnvVar: "DELIVERY_CLUSTER",
	})

	deliveryClusterCredentials := app.String(cli.StringOpt{
		Name:   "deliveryClusterCredentials",
		Desc:   "Credentials of the delivery cluster",
		EnvVar: "DELIVERY_CLUSTER_CREDENTIALS",
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

	app.Action = func() {
		delivery := metadata.GetCluster(*deliveryCluster, *deliveryClusterCredentials)
		publishing := metadata.GetCluster(*publishingCluster, *publishingClusterCredentials)
		cmr := metadata.GetCluster(*cmrAddress, *cmrCredentials)

		cmrReader, err := metadata.NewV1MetadataReadService(cmr)
		if err != nil {
			log.Errorf("Cannot start application: %s", err)
			return
		}
		
		contentService := metadata.InitContentService(delivery)
		mp := metadata.NewV1MetadataPublishService(contentService, publishing, cmrReader, *source, *batchSize)
		mp.Publish()
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Errorf("Cannot start application: %s", err)
	}
}
