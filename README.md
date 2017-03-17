V1-Metadata-Publisher
=====================

This is a one off application that publishes V1-Metadata for all the content in UPP. 

This service gets all the content UUIDs from document-store-api, saves them into a temporary file and then reads each UUID, gets the metadata by calling the binding-service and sends the result to the cms-metadata-notifier.

Usage
---------

```
    go get github.com/Financial-Times/v1-metadata-publisher`

    cd $GOPATH/src/github.com/Financial-Times/v1-metadata-publisher

    go build

    ./v1-metadata-publisher
```

__With Docker:__
```bash
docker pull coco/v1-metadata-publisher:latest

##Set environment variables:

## Ger URL for content UUIDs
export DELIVERY_CLUSTER=

## Delivery cluster credentials as username:password
export DELIVERY_CLUSTER_CREDENTIALS=

## Metadata publishing URL
export PUBLISHING_CLUSTER=

## Publishing cluster credentials as username:password
export PUBLISHING_CLUSTER_CREDENTIALS=

## URL of binding-service (this must contain placeholder for {source} and {uuid})
export CMR_ADDRESS=

## Binding service credentials as username:password
export CMR_CREDENTIALS=

## the source of content to be published (valid values: METHODE or BLOGS)
export SOURCE=

## number of requests to be sent / second
export BATCH_SIZE=

docker run -ti \
    --env DELIVERY_CLUSTER=$DELIVERY_CLUSTER \
    --env DELIVERY_CLUSTER_CREDENTIALS=$DELIVERY_CLUSTER_CREDENTIALS \ 
    --env PUBLISHING_CLUSTER=$PUBLISHING_CLUSTER \ 
    --env PUBLISHING_CLUSTER_CREDENTIALS=$PUBLISHING_CLUSTER_CREDENTIALS \ 
    --env CMR_ADDRESS=$CMR_ADDRESS \
    --env CMR_CREDENTIALS=$CMR_CREDENTIALS \ 
    --env SOURCE=$SOURCE \
    --env BATCH_SIZE=$BATCH_SIZE  \
    coco/v1-metadata-publisher:latest
```

