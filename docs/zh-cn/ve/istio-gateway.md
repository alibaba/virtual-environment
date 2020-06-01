# 使用Istio Gateway

在基于Istio的实现中，VirtualEnvironment路由规则最终会以VirtualService资源体现，若集群中的某些服务使用了Istio Gateway对外暴露访问入口，则应当为相应的Service资源添加名为`kt-virtual-environment/gateways`和`kt-virtual-environment/hosts`的Annotation，多个值以逗号“,”分隔。例如：

```yaml
apiVersion: v1
kind: Service
metadata:
  annotations:
    kt-virtual-environment/gateways: demo-gateway
    kt-virtual-environment/hosts: demo.com,test.com
```

这些配置将自动填入到生成的VirtualService资源的`spec.gateways`和`spec.hosts`属性中。
