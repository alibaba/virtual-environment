# 跨集群隔离

默认情况下，VirtualEnvironment实例产生的隔离规则仅对所在Namespace内的Pod有效。但从原理上说，这种隔离能力是可以跨越Namespace以及跨越集群使用的。

当请求从一个Pod发送到另一个Namespace甚至另一个Kubernetes集群的Pod中，如果目标Pod所在的集群部署有VirtualEnvironment CRD，且所在的Namespace具有相同配置的VirtualEnvironment实例，则该请求依然会在目标Pod所在的Namespace内遵循相同隔离规则进行路由。

![cross-cluster.jpg](https://img.alicdn.com/imgextra/i3/O1CN01DV8hTa1EdKpgmETa8_!!6000000000374-0-tps-2154-932.jpg)
