#!/bin/bash
set -euo pipefail

# This script automates the steps needed to remove the polycube
# sidecar injector. 

# Remove the deployment
kubectl delete deployment polycube-sidecar-injector 

# Remove the service 
kubectl delete service polycube-sidecar-injector 

# Remove the secret
kubectl delete secret polycube-sidecar-injector-certs

# Remove the web hook
kubectl delete mutatingwebhookconfiguration polycube-sidecar-injector-webhook

# Remove the configmap
kubectl delete configmap polycube-sidecar-configmap

# Remove the csr
kubectl delete csr polycube-sidecar-injector.default