# 扩展路由规则

## 使用Istio Gateway

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

## 自定义VirtualService规则

对于更复杂的情况，譬如修改指定服务的URI重写规则，可用通过`kt-virtual-environment/rule`注解来定义。

注意：为了确保VirtualEnvironment能力的通用性，该配置的值为字符串类型。在基于Istio的实现上，请参考[HTTPRoute](https://github.com/istio/api/blob/1.7.0/networking/v1alpha3/virtual_service.pb.go#L689)结构以Json格式定义配置。

例如：

```yaml
apiVersion: v1
kind: Service
metadata:
  annotations:
    kt-virtual-environment/rule: '{"Rewrite":{"Uri":"/prefix"},"Match":[{"Uri":{"Prefix":"/prefix"}}]}'
```

则生成的VirtualService对象会包含以下额外规则：

```yaml
  - match:
    - headers:
        ...
      uri:  
        prefix: /prefix/
    rewrite:
      uri: /
    route:
    - destination:
        ...
```
