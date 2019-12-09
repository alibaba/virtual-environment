# Virtual Environment Operator

阿里测试环境服务隔离和联调机制的Kubernetes版实现，基于Istio。

阅读[这里](https://yq.aliyun.com/articles/700766)了解更多故事。

```
            +----------+   +----------+   +----------+
DailyEnv    | ServiceA |   | ServiceB |   | ServiceC |
            +----------+   +-^------+-+   +----------+
                             |      |
            +----------+     |      |     +----------+
ProjectEnv -> ServiceA +-----+      +-----> ServiceC |
            +----------+                  +----------+
```

原理为根据Pod上的Label将服务划分为独立的虚拟环境，并根据Label值自动生成级联路由规则（即相应的Istio规则对象）。

例如，假设将标识Label配置为`virtual-env`，环境级联符号配置为`.`，则：

- 所有包含Label为`virtual-env: dev`和`virtual-env: dev.proj1`的Pod分别归属名称为`dev`和`dev.proj1`的两个虚拟环境
- 任何情况下，来自`dev.proj1`虚拟环境的HTTP请求，如果目标服务在该虚拟环境中不存在，则自动由上一级（即`dev`）虚拟环境中的Pod响应（如示意图中ServiceA调用ServiceB）
- 在同一链路中，即使途径上级（即`dev`）虚拟环境，在后续链路中，依然应当优先返回到来源（即`dev.proj1`）虚拟环境的Pod响应（如示意图中ServiceB调用ServiceC）
- 如果环境名存在多个级联符号，例如`dev.proj1.user1`，则在找不到合适的目标Pod时，将逐级向上查找，直到找到路由的Pod为止

相比集团的路由隔离和兜底能力，此方案主要做了两点增强：

1. 基于Pod Label自动生成虚拟环境和级联路由规则，而不是固定的`项目环境`到`日常环境`兜底约定
2. 理论上支持任意多层虚拟环境级联

限制条件：

- 由于本质是动态生成Istio规则，因此仅支持Istio可配置的通信协议，目前为HTTP
- 在应用程序中需要实现Header标签在请求之间的传递，可通过OpenTracing的baggage机制完成，也可在请求代码中直接传递

## 使用概述

**[0]** 前提：集群已经部署Istio

**[1]** 部署CRD和Operator
```bash
kubectl apply -f deploy/crds/env.alibaba.com_virtualenvironments_crd.yaml
kubectl apply -f deploy/operator.yaml
```
如果开启了RBAC，还需要部署相应的Role和ServiceAccount
```bash
kubectl apply -f deploy/service_account.yaml
kubectl apply -f deploy/role.yaml
kubectl apply -f deploy/role_binding.yaml
```

**[2]** 为应用程序添加透传标签Header的功能（默认约定Header为`X-Virtual-Env`）

**[3]** 将改程序打包为镜像，并在部署到Kubernetes时，为Deployment的Pod模板增加一个表示虚拟环境名称的Label（默认约定为`virtual-env`）

**[4]** 创建类型为`VirtualEnvironment`的资源定义文件（参考`deploy/crds/env.alibaba.com_v1alpha1_virtualenvironment_cr.yaml`），根据实际情况修改配置参数，使用`kubectl apply`命令添加到Kubernetes集群

## VirtualEnvironment配置

```yaml
apiVersion: env.alibaba.com/v1alpha1
kind: VirtualEnvironment
metadata:
  name: example-virtualenv
spec:
  envHeader:
    name: X-Virtual-Env
    autoInject: true
  envLabel:
    name: virtual-env
    splitter: .
    defaultSubset: dev
  instancePostfix: kt-env
```

| 配置参数                | 默认值         | 说明  |
| :--------              | :-----:       | :---- |
| envHeader.name         | X-Virtual-Env | 用于记录虚拟环境名的HTTP头名称（虽然有默认值，强烈建议显性设置） |
| envHeader.autoInject   | false         | 是否为没有虚拟环境HTTP头记录的请求自动注入HTTP头（建议开启） |
| envLabel.name          | virtual-env   | Pod上标记虚拟环境名用的标签名称（除非确实必要，建议保留默认值） |
| envLabel.splitter      | .             | 虚拟环境名中用于划分环境默认路由层级的字符（只能是单个字符） |
| envLabel.defaultSubset |               | 请求未匹配到任何存在的虚拟环境时，进行兜底虚拟环境名（默认为随机路由） |
| instancePostfix        |               | 自动创建的Istio对象命名的名字尾缀，默认与相应服务名相同（无尾缀） |
