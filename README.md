V1-Metadata-Publisher
=====================

This is a one off application that publishes V1-Metadata for all the content in UPP. 

This service gets all the content UUIDs from document-store, gets the metadata by calling the binding-service and sends the result to the cms-metadata-notifier.

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
docker pull coco/v1-metadata-publisher:{latest_version}

##Set environment variables:

## Address of document-store
export DELIVERY_CLUSTER=

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
```
```bash
#Run docker image
docker run \
    -e "DELIVERY_CLUSTER=$DELIVERY_CLUSTER" \
    -e "PUBLISHING_CLUSTER=$PUBLISHING_CLUSTER" \
    -e "PUBLISHING_CLUSTER_CREDENTIALS=$PUBLISHING_CLUSTER_CREDENTIALS" \
    -e "CMR_ADDRESS=$CMR_ADDRESS" \
    -e "CMR_CREDENTIALS=$CMR_CREDENTIALS" \
    -e "SOURCE=$SOURCE" \
    -e "BATCH_SIZE=$BATCH_SIZE"  \
    coco/v1-metadata-publisher:{latest_version}
```
__NB:__ This app supposes that there is an ssh tunnel between the host where Mongo runs the the local machine:
```bash
ssh -L {localhost_private_ip}:27020:localhost:27020 {username}@{mongo_instance_ip}
``` 
