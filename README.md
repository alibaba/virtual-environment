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

例如，假设将标识Label配置为`virtualEnv`，环境级联符号配置为`/`，则：

- 所有包含Label为`virtualEnv: dev`和`virtualEnv: dev/proj007`的Pod分别归属名称为`dev`和`dev/proj007`的两个虚拟环境
- 任何情况下，来自`dev/proj007`虚拟环境的HTTP请求，如果目标服务在该虚拟环境中不存在，则自动由上一级（即`dev`）虚拟环境中的Pod响应（如示意图中ServiceA调用ServiceB）
- 在同一链路中，即使途径上级（即`dev`）虚拟环境，在后续链路中，依然应当优先返回到来源（即`dev/proj007`）虚拟环境的Pod响应（如示意图中ServiceB调用ServiceC）
- 如果环境名存在多个级联符号，例如`dev/proj007/feature001`，则在找不到合适的目标Pod时，将逐级向上查找，直到找到路由的Pod为止

相比集团的路由隔离和兜底能力，此方案主要做了两点增强：

1. 基于Pod Label自动生成虚拟环境和级联路由规则，而不是固定的`项目环境`到`日常环境`兜底约定
2. 理论上支持任意多层虚拟环境级联

限制条件：

- 由于本质是动态生成Istio规则，因此仅支持Istio可配置的通信协议，目前为HTTP
- 在应用程序中需要实现Header标签在请求之间的传递，可通过OpenTracing的baggage机制完成，也可在请求代码中直接传递

## 使用示例

1. 准备添加了自动透传Header标签能力的应用程序（假设约定Header为`X-Virtual-Env`）
2. 将改程序打包为镜像，并在部署到Kubernetes时，为Deployment的Pod模板增加一个Label项（假设为`virtualEnv`）
3. 创建一个`VirtualEnvironment`的YAML文件，然后使用kubectl apply命令将它添加到Kubernetes集群

