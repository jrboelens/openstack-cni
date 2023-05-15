# Overview

`openstack-cni` is a [CNI](https://github.com/containernetworking/cni) plugin that provides the ability to dynamically configure Openstack Neutron Ports and attach them to Pod containers. 

Configuration is achieved by created `NetworkAttachmentDefinitions` and referencing them with Pod Annotations.

Such configuration allows the implementer to configure a port for any OpenStack network and attach it to any Pod.


These [sequence diagrams](docs/diagrams.md) show the workflow and relationship between `kubelet`, `multus`, `openstack-cni` and `openstack-cni-daemon`.

# Installation

* Ensure that [multus-cni](https://github.com/k8snetworkplumbingwg/multus-cni/blob/master/README.md) is installed
* Add the required secret.

For example:
```
---
apiVersion: v1
kind: Secret
metadata:
  name: openstack-cni-secret
  namespace: mynamespace
type: Opaque
stringData:
  OS_PASSWORD: SECRETPASSWORD
```
* Create a helm values file ([example](helm/example-values.yaml))
* Run helm (`helm upgrade openstack-cni helm/ --install`)
* Create a pod with the proper annotations.

For example:
```
apiVersion: v1
kind: Pod
metadata:
  annotations:
    k8s.v1.cni.cncf.io/networks: '[{"name": "mycloud-network1", "interface": "ens37"},{"name": "mycloud-network2", "interface": "ens42"}]'
  name: dummypod
  namespace: mynamespace
...
```

# CNI spec

The `spec.config` portion of the `NetworkAttachmentDefinition` should contain the following configuration:

### Fields
* `cniVersion` is required
* `type` is required and must be `openstack-cni`
* `network` is required
* `project_name` is optional, but required if `subnet_name` is specified
* `subnet_name` is optional
* `security_groups` is optional

### Example
```
spec:
  config: '{
        "cniVersion": "0.3.1",
        "type": "openstack-cni",
        "name": "service-ingress",
        "network": "my-openstack-network",
        "project_name": "my-openstack-project-name",
        "subnet_name": "my-openstack-subnet",
        "security_groups": ["project_default", "default"],
        }'
```

# Testing

In order to run the full test suite valid OpenStack credentials must be present in the environment.

`testing.conf` will be sourced if present.

The following enviroment variables an be used to control the tests.

```
OS_TESTS="1" ## 0 = skip openstack tests 1 = execute openstack tests
CNI_CONFIG_FILE="../../config.conf" # path to the main config file
OS_VM_NAME="mytestvm"
OS_NETWORK_NAME="myproject-network"
OS_PORT_NAME="mytestport"
OS_PROJECT_NAME="myproject"
OS_SUBNET_NAME="myproject-subnet"
OS_SECURITY_GROUPS="default;project_default"
```

# HTTP Server End points

* `GET /health` - returns the health of the server including whether OpenStack authentication is working
* `GET /ping` - returns "PONG"
* `POST /cni` - handles `ADD/DEL/CHECK` CNI commands

# Environment Variables
### Runtime:
* `OS_PROJECT_NAME` - required

* `CNI_API_URL` - url `openstack-cni` will used to contact `openstack-cni-daemon`.  Also overrides `openstack-cni-daemon`'s listen address (`http://127.0.0.1:4242`)
* `CNI_CACHE_TTL` - cache ttl (`300s`)
* `CNI_CONFIG_FILE` - configuration file `openstack-cni` reads (`/etc/cni/net.d/openstack-cni.conf`)
* `CNI_MIN_PORT_AGE` - minimum age of ports to be cleaned up (`300s`)
* `CNI_READ_TIMEOUT` - http server read timeout (`10s`)
* `CNI_REAP_INTERVAL` - the port cleanup interval (`300s`)
* `CNI_REQUEST_TIMEOUT` - `openstack-cni`'s request timeout in seconds (`60`)
* `CNI_WRITE_TIMEOUT` - http server write timeout (`10s`)
* `OS_REGION_NAME` - OpenStack region (`RegionOne`)

### Testing:
The following vars control the test that interact directly with the OpenStack APIs
* `OS_TESTS` - `0` = skip OpenStack tests `1` = perform OpenStack tests
* `OS_PROJECT_NAME` - required when `OS_TESTS=1`
* `OS_NETWORK_NAME` - required when `OS_TESTS=1`
* `OS_SECURITY_GROUPS` - required when `OS_TESTS=1`
* `OS_SUBNET_NAME` - required when `OS_TESTS=1`
* `OS_PORT_NAME` - optionally override the port name
* `OS_VM_NAME` - optionally tell the OpenStack tests to use a hostname other than `os.Hostname()`


For local testing, configuration and secrets can be loaded from `config.conf` or `secrets.conf`.
