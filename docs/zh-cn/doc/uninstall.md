# 移除KtEnv

KtEnv的卸载包括移除各个Namespace中的Operator对象以及移除全局的CRD和Webhook组件。

首先在使用了虚拟环境功能的Namespace里移除VirtualEnvironment资源以及Operator实例，以`default` Namespace为例：

```bash
kubectl label Namespace default environment-tag-injection-
kubectl delete -n default VirtualEnvironment --all
kubectl delete -n default Deployment virtual-env-operator
```

若部署过ServiceAccount资源，也应当相应移除：

```bash
kubectl delete -n default ServiceAccount virtual-env-operator
kubectl delete -n default Role virtual-env-operator
kubectl delete -n default RoleBinding virtual-env-operator
```

最后移除全局的KtEnv组件：

```bash
kubectl delete CustomResourceDefinition virtualenvironments.env.alibaba.com
kubectl delete Namespace kt-virtual-environment
```
