# openGemini-operator

## Description
Cloud Native Installation and Deployment of [opengemini](https://github.com/openGemini/openGemini) database using K8S.

## Getting Started
You’ll need a Kubernetes cluster to run against. You can use [KIND](https://sigs.k8s.io/kind) to get a local cluster for testing, or run against a remote cluster.
**Note:** Your controller will automatically use the current context in your kubeconfig file (i.e. whatever cluster `kubectl cluster-info` shows).

### Running on the cluster

1. Install helm:

```
curl https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash
```

2. Add chart repo:

```
helm repo add opengemini https://opengemini.github.io/openGemini-operator
```

3. Install chart(both CRDs and operator) in specified namespace:

```
kubectl create ns opengemini-system
helm install opengemini-operator opengemini/operator -n opengemini-system
```

4. Deploy cluster in specified namespace, you can copy [example configuration file](config/samples/_v1_opengeminicluster.yaml) to local path and modify it:

```
kubectl create ns opengemini
kubectl -n opengemini create -f config/samples/_v1_opengeminicluster.yaml
```

5. Uninstall chart:

```
helm uninstall opengemini-operator -n opengemini-system
```

## Contributing

Please read the [Contribution guide](https://github.com/openGemini/openGemini/blob/main/CONTRIBUTION.md) of opengemini repo first.

### How it works
This project aims to follow the Kubernetes [Operator pattern](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/).

It uses [Controllers](https://kubernetes.io/docs/concepts/architecture/controller/)
which provides a reconcile function responsible for synchronizing resources until the desired state is reached on the cluster.

### Test It Out
1. Install the CRDs into the cluster:

```sh
make install
```

2. Run your controller (this will run in the foreground, so switch to a new terminal if you want to leave it running):

```sh
make run
```

**NOTE:** You can also run this in one step by running: `make install run`

### Modifying the API definitions
If you are editing the API definitions, generate the manifests such as CRs or CRDs using:

```sh
make manifests
```

### Uninstall CRDs
To delete the CRDs from the cluster:

```sh
make uninstall
```

### Undeploy controller
UnDeploy the controller to the cluster:

```sh
make undeploy
```

**NOTE:** Run `make --help` for more information on all potential `make` targets

More information can be found via the [Kubebuilder Documentation](https://book.kubebuilder.io/introduction.html)

## License

Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
