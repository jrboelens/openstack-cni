# Overview

`openstack-cni` is a [CNI](https://github.com/containernetworking/cni) plugin designed to allow OpenStack Neutron ports to be used directory by Pod containers.


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

* `CNI_API_URL` - optionally the url `openstack-cni` will used to contact `openstack-cni-daemon`.  Also overrides `openstack-cni-daemon`'s listen address (`http://127.0.0.1:4242`)
* `CNI_CONFIG_FILE` - optionally override the configuration `openstack-cni` reads (`/etc/cni/net.d/openstack-cni.conf`)
* `CNI_REQUEST_TIMEOUT` - optionally `openstack-cni`'s request timeout in seconds (`60`)
* `OS_REGION_NAME` - optionally override the region (`RegionOne`)
* `CNI_READ_TIMEOUT` - optionally override the http server read timeout (`10s`)
* `CNI_WRITE_TIMEOUT` - optionally override the http server write timeout (`10s`)
* `CNI_REAP_INTERVAL` - optionally override the port cleanup interval (`300s`)
* `CNI_CACHE_TTL` - optionally override the cache ttl  (`300s`)

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
