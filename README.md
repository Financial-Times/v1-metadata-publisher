V1-Metadata-Publisher
=====================

This is a one off application that publishes V1-Metadata for all the content in UPP. 

This service gets all the content UUIDs from document-store-api, saves them into a temporary file and then reads each UUID, gets the metadata by calling the binding-service and sends the result to the cms-metadata-notifier.

Usage
---------

`go get github.com/Financial-Times/v1-metadata-publisher`

`./v1-metadata-publisher`

__With Docker:__

`docker pull coco/v1-metadata-publisher`

`docker run -ti 
    --env DELIVERY_CLUSTER=<delivery_cluster_doc_store> 
    --env DELIVERY_CLUSTER_CREDENTIALS=<username:password> 
    --env PUBLISHING_CLUSTER=<publishing_cluster_cms_metadata_publisher>  
    --env PUBLISHING_CLUSTER_CREDENTIALS=<username:password>  
    --env CMR_ADDRESS=<binding_service_address> 
    --env CMR_CREDENTIALS=<username:address> 
    --env SOURCE=<content_source> 
    --env BATCH_SIZE=<batch_size>  
    coco/v1-metadata-publisher`

