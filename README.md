# Overview

`openstack-cni` is a [CNI](https://github.com/containernetworking/cni) plugin designed to allow OpenStack Neutron ports to be used directory by Pod containers.


# Requirements

`multus-cni` is a pre-requisite. Instructions for installing `multus-cni` can be found [Here](https://github.com/k8snetworkplumbingwg/multus-cni/blob/master/README.md)


# Configuration

`openstack-cni-daemon` configuration and secrets are injected into the environment via a volume mounted ConfigMap and Secret.

### Example ConfigMap

```
apiVersion: v1
kind: ConfigMap
metadata:
  name: openstack-cni-config
  namespace: mynamespace
data:
  OS_AUTH_URL: https://keystone.mycloud.com:5000/v3
  OS_USERNAME: mycloud-user
  OS_PROJECT_NAME: mycloud-project
  OS_DOMAIN_NAME: default
  CNI_API_URL: http://127.0.0.1:4242
```

### Example Secret

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

For local testing, configuration and secrets can be loaded from `config.conf` or `secrets.conf`.

## CNI spec

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

In order to run the openstack related tests valid credentials must be provided.

The tests require an extra config file named `testing.conf`

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

Running `make test` will run all of the tests.

# HTTP Server End points

* `GET /health` - returns the health of the server including whether OpenStack authentication is working
* `GET /ping` - returns "PONG"
* `GET /state/{containerId}/{ifname}` - returns the state for an container/interface tuple
* `DELETE /state/{containerId}/{ifname}` - deletes the state for a container/interface tuple
* `POST /state` - sets the state for a container/interface tuple
* `POST /cni` - handles `ADD/DEL/CHECK` CNI commands
