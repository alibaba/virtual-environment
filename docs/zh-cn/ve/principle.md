# 虚拟环境原理

虚拟环境Operator的本质是生成Service Mesh控制面规则，当前实现基于Istio开源版本。

Operator启动后会遍历所在Namespace里所有Deployment的Pod模板上的路由标签（Label）， 为每个Service计算生成每种染色流量（Header）的Subset访问规则，然后持续监听Service和Deployment对象事件，动态创建和调整VirtualService和DestinationRule资源。

![calculate-rule](https://virtual-environment.oss-cn-zhangjiakou.aliyuncs.com/image/calculate-rule-zh-cn.jpg)
