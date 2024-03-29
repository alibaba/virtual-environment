# 适配Spring Cloud开发框架

由于KtEnv的网络隔离能力依赖于Kubernetes的服务路由机制，而Spring Cloud框架采用的客户端服务发现机制会绕过Kubernetes的路由过程，流量直接定向到达目标服务的Pod IP，导致虚拟环境功能失效。

为了兼容Kubernetes的路由机制，一种方法是在向Spring Cloud服务注册中心注册实例时，使用服务相应的Kubernetes Service资源名称作为目标地址，具体做法为在程序的`application.properties`或`application.yaml`文件中，根据所用的服务注册中心类型，将相应的注册地址改为Service资源名字。

这样一来，实际上是使用Kubernetes的服务发现机制替代了Spring Cloud原有的相应能力，可能会导致某些依赖该机制的功能，譬如注册中心会将相同服务的所有实例合并显示为一个，但对于测试环境而言不会带来实际影响。

假设Service资源在Kubernetes里的`spec.name`值为`app-js`，若使用Eureka作为服务注册中心，则配置如下：

```properties
eureka.instance.hostname = app-js
```

若使用Consul作为服务注册中心，则配置为：

```properties
spring.cloud.consul.host = app-js
```

若使用Nacos作为服务注册中心，则配置为：

```properties
spring.cloud.nacos.discovery.ip = app-js
```

与此同时，由于Spring Cloud的`RibbonClient`组件在发送HTTP请求时默认不会携带Istio路由所需的`Host`请求头<sup>[1]</sup>，开发者还需要自行添加HTTP Client的拦截器来在请求发出前添加内容为`Host: <当前服务名>`的Header，使得请求能被Istio Sidecar处理。

为测试环境中的所有服务采用如上修改后，即可正常使用KtEnv的环境隔离功能了。

特别感谢：
- @[宋锦丰](https://github.com/SJFCS) 对SpringCloud使用虚拟环境时缺失`Host`请求头问题的补充（友情链接：[TA的博客](https://songjinfeng.com/)）

参考引用：
- [1] https://github.com/IBM/spring-cloud-kubernetes-with-istio/blob/master/README.md#discovery-client-and-istio
