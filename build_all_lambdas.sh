#!/bin/bash

cd ./services

lambdas=()

for d in *Service/ ; do
    cd ./$d

    if [ -e cmd/ ]; then
        lambdas+=($d)
    fi

    cd ..
done

lambdas_count=${#lambdas[@]}
echo "$lambdas_count lambdas found: ${lambdas[@]}"

for (( i=0; i<$lambdas_count; i++ )); do
    lambda=${lambdas[$i]}

    cd $lambda/cmd
    echo "Building $lambda lambda"

    lambda_name=${lambda/Service\//Lambda}.zip

    build_and_upload_lambda --output-file-name $lambda_name --output-dir ../../../terraform/aws/files --quiet
    echo ""
    
    cd ../..
done