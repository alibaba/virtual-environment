# 功能扩展

如需支持其他类型的Service Mesh控制面，只需要实现[IsolationRouter](https://github.com/alibaba/virtual-environment/blob/master/pkg/component/router/router_interface.go)接口的方法，自定义新的Router类型，并注册到[router_pool.go](https://github.com/alibaba/virtual-environment/blob/master/pkg/component/router/router_pool.go)的pool列表中，然后在创建虚拟环境实例时选择路由类型即可（**注意**：此配置当前尚未实现，如有需要请留言提出）。

## 源码构建VirtualEnvironment

请使用Golang 1.13或以上版本，并安装[Docker](https://docs.docker.com/)和[operator-sdk](https://github.com/operator-framework/operator-sdk)工具。

项目分为Operator和Admission Webhook两个独立二进制产物，可通过Makefile直接打包出镜像。

构建Operator镜像：

```bash
make build-operator
```

构建Webhook镜像：

```bash
make build-admission
```

可通过`operator-sdk`工具本地启动并连接到集群运行：

```bash
export OPERATOR_NAME=virtual-env-operator
operator-sdk run local --watch-namespace=<指定一个目标Namespace>
```
