# 限制条件

虚拟环境的路由能力源于Service Mesh控制面规则。由于当前开源版本仅实现了基于Istio的路由规则生成，因而受限于Istio的能力。

主要包括：

- 无法支持框架自带的非Kubernetes原始服务发现机制，如Dubbo和SpringCloud等需要[进行适配](zh-cn/ve/spring-cloud.md)
- 无法支持非HTTP协议的通信，目前普通四层TCP流量在Istio中不能进行精细路由管控

此外，由于Sidecar机制不会侵入应用程序内部逻辑，因而需要在应用程序中自行实现Header标签在请求之间的传递。若项目已经在使用OpenTracing SDK，可复用其baggage机制完成标签透传。也可[使用SDK](zh-cn/ve/use-sdk.md)或直接在代码中实现传递。
