#!/bin/bash

collections=(users passwordHashes userPasswordHashMapping sessions userSessionMapping)

echo "########### Loading data to Mongo DB ###########"

for i in ${!collections[*]}; do
    collection=${collections[$i]}
    echo "$(($i+1)) Loading $collection"
    
    if [ -f /tmp/data/$collection.json ]; then
        mongoimport \
            --jsonArray \
            --db $MONGO_INITDB_DATABASE \
            --collection $collection \
            --file /tmp/data/$collection.json
    else
        echo "No data for $collection"
    fi
done