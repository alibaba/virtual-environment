Virtual Environment Operator
===========

阿里测试环境服务隔离和联调机制的Kubernetes版实现，当前基于Istio。

![isolation](https://virtual-environment.oss-cn-zhangjiakou.aliyuncs.com/image/diagram-zh-cn.jpg)

详见☞[项目文档](https://alibaba.github.io/virtual-environment/#/zh-cn/)☜

[English Instruction](./README_EN.md)

## 概述

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

限制条件：

- 由于本质是动态生成Istio规则，因此仅支持Istio可配置的通信协议，目前为HTTP
- 由于Istio目前仅支持Kubernetes内置DNS的服务发现，因此对于自带服务发现的框架（如Bubbo、SpringCloud）暂时无法使用
- 在应用程序中需要实现Header标签在请求之间的传递，可通过OpenTracing的baggage机制完成，也可在请求代码中直接传递

## 联系我们

请加入`kt-dev`钉钉群：

<img src="https://virtual-environment.oss-cn-zhangjiakou.aliyuncs.com/image/dingtalk-group-zh-cn.jpg" width="50%"></img>
