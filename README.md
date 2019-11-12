# Polycube Sidecar Injector

Kubernetes Mutating Webhook that automatically injects Polycube as a sidecar in pods matching some criteria.

## Polycube
``Polycube`` is an **open source** software framework that provides **fast** and **lightweight** **network functions** such as bridges, routers, firewalls, and others.

Polycube services, called `cubes`, can be composed to build arbitrary **service chains** and provide custom network connectivity to **namespaces**, **containers**, **virtual machines**, and **physical hosts**.

For more information, jump to the project [Documentation](https://polycube-network.readthedocs.io/en/latest/).

## Polycube as a sidecar

From monitoring to security purposes, Polycube running as a sidecar in your pods may bring several benefits. In case your CNI does not provide firewall capabilities or, for some reason, you don't want to use that one, you may leverage on Polycube's API to create a firewall inside the pods you want to protect; or instantiate a DDOS mitigator to reduce the impact of DDOS attacks. 
The aforementioned situations are just two simple examples, refer to the documentation to know more about all the features and network functions that Polycube provides. 

### CNI requirements

Running ``pcn-k8s`` (Polycube's own [CNI](https://polycube-network.readthedocs.io/en/latest/components/k8s/pcn-kubernetes.html)) as your CNI of choice is recommended, as it can be made aware of the presence of the sidecar injector and, thus, make the proper adjustments to help it be more efficient. Nonetheless, the sidecar injector is CNI-agnostic and has no requirements about the CNI installed.

### Injection requirements

Polycube will be injected as a sidecar only in pods that match some particular criteria. Once the sidecar injector is installed, it will work only on pods that have following annotation: ``polycube.network/sidecar`` with value ``enabled``. Additionally, such pods must run on namespaces that have the mentioned key/pair as label.

### Example

In this example, we will deploy a pod that will be injected with the Polycube sidecar.

Supposing that the namespace where you want to deploy such pod is called ``enabled-ns``, you need to first label it with the neabled label:

``kubectl label ns enabled-ns polycube.network/sidecar=enabled``

Deploy the pod: 

```bash
cat <<EOF | kubectl apply -f -
apiVersion: v1
kind: Pod
metadata:
  namespace: enabled-ns
  name: myapp-pod
  annotations:
    polycube.network/sidecar: enabled
  labels:
    app: myapp
spec:
  containers:
  - name: myapp-container
    image: busybox
    command: ['sh', '-c', 'echo Hello Kubernetes! && sleep 3600']
EOF
```

After some time, you will see that the pod has 2 containers running inside it:

```bash
kubectl get pods -n enabled-ns
NAME        READY   STATUS    RESTARTS   AGE
myapp-pod   2/2     Running   0          90s
```

### Interact with the polycube sidecar

You can interact with polycube's API by contacting the pod's IP on port 9000. Once again, refer to the [Documentation](https://polycube-network.readthedocs.io/en/latest/) to know more.

## Installation

In order to launch the sidecar injector, run the ``deploy.sh`` script inside the ``scripts`` folder.

Please make sure you have [CFSSL](https://github.com/cloudflare/cfssl) installed before running the script:

``sudo apt install golang-cfssl``

### Remove

Run the ``remove.sh`` script inside the ``scripts`` folder to remove every resource deployed by the sidecar injector.

### Configuration

The sidecar injector is set to inject the latest polycube docker image. But in case this does not suit your needs, i.e. if you have compiled and uploaded a version of polycube with only the firewall component present, you may edit the ``polycubeImage`` field in the ``sidecar-configmap.yaml`` file inside the ``deployment`` folder:

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: polycube-sidecar-configmap
data:
  sidecarconfig.yaml: |
    polycubeImage: user/image:tag
```