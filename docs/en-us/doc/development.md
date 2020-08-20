## Build from source code

Please use Golang version 1.13 or above, [Docker](https://docs.docker.com/) and [operator-sdk](https://github.com/operator-framework/operator-sdk) are also required.

This project contains two parts: the `Operator` and the `Admission Webhook`, both can be build using Makefile lay on the root folder of the repository.

Build `Operator` docker image:

```bash
make build-operator
```

Build `Admission Webhook` docker image:

```bash
make build-admission
```

You could use `operator-sdk` tool to directly run the operator in a kubernetes cluster from local:

```bash
export OPERATOR_NAME=virtual-env-operator
operator-sdk run local --watch-namespace=<指定一个目标Namespace>
```

# Document preview

Use `docsify` tool to preview edited document locally:

```bash
docsify serve docs
```

# Function extension

In order to support other Service Mesh controller, you should implement [IsolationRouter](https://github.com/alibaba/virtual-environment/blob/master/pkg/component/router/router_interface.go) interface and define a new Router class, and register it to the pool in [router_pool.go](https://github.com/alibaba/virtual-environment/blob/master/pkg/component/router/router_pool.go). Then select route type when configuring the virtual environment instance (**Notice**: this configure item currently not exist yet, please raise a github issue when required)
