#!/bin/bash

# This script generates a private key and a certificate signing request.
# Later, a CertificateSigningRequest yaml object is sent to Kubernetes and
# then approved. Finally, a secret is created and deployed to it.
# For more information: https://kubernetes.io/docs/tasks/tls/managing-tls-in-a-cluster/
# Credits to ExpediaDotCom's kubernetes-sidecar-injector for the original script,
# although, later on, a different approach has been adopted instead.
# (https://github.com/ExpediaDotCom/kubernetes-sidecar-injector/)

set -euo pipefail

# Does CFSSL exist?
if [[ ! "$(command -v cfssl)" ]]; then
    echo "CFSSL not found"
    exit 1
fi

# Some variables 
SERVICE=polycube-sidecar-injector
SECRET=polycube-sidecar-injector-certs
NAMESPACE=default
CSRNAME=${SERVICE}.${NAMESPACE}

# Directories
CURRDIR=$(cd `dirname $0` && pwd)
TMPDIR=$(mktemp -d)

# Create the certificate signing request 
echo "--- Create certificate signing request"
cd $TMPDIR
cat <<EOF | cfssl genkey - | cfssljson -bare server
{
  "hosts": [
    "${SERVICE}",
    "${SERVICE}.${NAMESPACE}",
    "${SERVICE}.${NAMESPACE}.svc"
  ],
  "CN": "${SERVICE}.${NAMESPACE}.pod.cluster.local",
  "key": {
    "algo": "ecdsa",
    "size": 256
  }
}
EOF
cd $CURRDIR

# Delete pre-existing CSR, if any. Ignore errors.
echo "--- Delete pre-existing Certificate Signing Requests"
kubectl delete csr ${CSRNAME} 2>/dev/null || true

# Deploy the Certificate Signing Request
echo "--- Deploy the Certificate Signing Requests Object"
REQUEST=$(cat ${TMPDIR}/server.csr | base64 | tr -d '\n')
cat <<EOF | kubectl apply -f -
apiVersion: certificates.k8s.io/v1beta1
kind: CertificateSigningRequest
metadata:
  name: ${CSRNAME}
spec:
  request: ${REQUEST}
  usages:
  - digital signature
  - key encipherment
  - server auth
EOF

# Wait until it's there
while true; do
    kubectl get csr ${CSRNAME}
    if [[ "$?" -eq 0 ]]; then
        echo "Certificate Signing request found"
        break
    fi
    sleep 1
done

# Approve the certificate
echo "--- Approve the certificate"
kubectl certificate approve ${CSRNAME}

# Get it
for i in $(seq 5); do
    SC=$(kubectl get csr ${CSRNAME} -o jsonpath='{.status.certificate}')
    
    # There?
    if [[ "$SC" != '' ]]; then
        echo "Signed certificate found."
	kubectl get csr ${CSRNAME} -o jsonpath='{.status.certificate}' | base64 --decode > ${TMPDIR}/server-cert.crt
	break
    fi

    # Not there yet?
    printf "Signed certificate is not present yet. " 
	
    # Is this the fifth time already? If so give up, dude.
    if [[ i -eq 5 ]]; then
      	printf "Going to give up.\n"
  	exit 1  
    fi
 
    printf "Trying again in 3 seconds.\n"
    sleep 3
done


# Create the secret
echo "--- Creating the secret ---"
kubectl create secret generic ${SECRET} \
    --from-file=cert.crt=${TMPDIR}/server-cert.crt \
    --from-file=key.pem=${TMPDIR}/server-key.pem \
    --dry-run -o yaml | kubectl -n ${NAMESPACE} apply -f -