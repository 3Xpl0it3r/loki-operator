#!/usr/bin/env bash

CURRENT_DIR=$(echo "$(pwd)/$line")
REPO_DIR="$CURRENT_DIR"
IMAGE_NAME="kubernetes-codegen:latest"

echo "Building codgen Docker image ...."
docker build -f "${CURRENT_DIR}/hack/docker/codegen.dockerfile" \
             -t "${IMAGE_NAME}" \
             "${REPO_DIR}" 
            

cmd="go mod tidy && /go/src/k8s.io/code-generator/generate-groups.sh  all  \
        "github.com/l0calh0st/loki-operator/pkg/client" \
        "github.com/l0calh0st/loki-operator/pkg/apis" \
        lokioperator.l0calh0st.cn:v1alpha1 -h /go/src/k8s.io/code-generator/hack/boilerplate.go.txt"
    
echo "Generating client codes ...."

docker run --rm -v "${REPO_DIR}:/go/src/github.com/l0calh0st/loki-operator" \
        "${IMAGE_NAME}" /bin/bash -c "${cmd}"
