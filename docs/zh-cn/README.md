# 介绍

KtVirtualEnvironment是Kt系列工具的一员，来自阿里巴巴云研发部。

项目实现了基于流量染色的虚拟环境隔离，适用于Kubernetes集群。可独立使用，或结合[KtConnect](https://alibaba.github.io/kt-connect/)工具实现本地到集群的流量路由控制，详见[典型场景](ve/typical-scenario.md)介绍。

## 起源

对于微服务的开发者而言，拥有一套干净、独占的完整测试环境无疑能够提高软件研发过程中的功能调试和异常排查效率。

然而在中大型团队里，为每位开发者维护一整套专用测试服务集群，从经济成本和管理成本上考虑都并不现实。为此阿里巴巴的研发团队采用了基于路由隔离的"虚拟环境"方法。

![isolation](https://virtual-environment.oss-cn-zhangjiakou.aliyuncs.com/image/diagram-zh-cn.jpg)

这种方法通过在请求上携带约定标签，可以将个别需要调试或测试的特定版本服务实例与其他公共服务实例组成临时虚拟集群，形成开发者视角的专属测试环境。

阅读[这里](https://yq.aliyun.com/articles/700766)了解更多故事。

## 特性

- 基于流量标记划分虚拟隔离域
- 隔离域之间复用大部分公共环境服务实例
- 用户随意在项目环境中部署、调试不会影响其他开发者
- 支持本地运行服务无需部署到集群，直接加入任意隔离域

## 联系我们

请加入`kt-dev`钉钉群：

<img src="https://virtual-environment.oss-cn-zhangjiakou.aliyuncs.com/image/dingtalk-group-zh-cn.jpg" width="40%"></img>
