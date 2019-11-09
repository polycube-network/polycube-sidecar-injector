#!/bin/bash
set -euo pipefail

# This script automates the steps needed to deploy and configure the polycube
# sidecar injector. 
# Credits to ExpediaDotCom's kubernetes-sidecar-injector for the original script,
# although, later on, a different approach has been adopted instead.
# (https://github.com/ExpediaDotCom/kubernetes-sidecar-injector/)

# Directories
CURRDIR=$(cd `dirname $0` && pwd)
DEPDIR=$(cd `dirname $0` && cd ../deployment && pwd)
TMPDIR=$(mktemp -d)

# Create the server certificates
${CURRDIR}/create-server-cert.sh

# Deploy the configmap used by the sidecar injector
echo "--- Deploy the sidecar configMap"
kubectl apply -f ${DEPDIR}/sidecar-configmap.yaml

# Deploy the sidecar injector deployment
kubectl apply -f ${DEPDIR}/sidecar-injector-deployment.yaml

# Deploy the mutating webhook with the appropriate CABUNDLE
CABUNDLE=$(kubectl get configmap -n kube-system extension-apiserver-authentication -o=jsonpath='{.data.client-ca-file}' | base64 | tr -d '\n')
cat ${DEPDIR}/mutating-webhook-template.yaml | sed -e "s|\${CABUNDLE}|${CABUNDLE}|g" | kubectl apply -f -