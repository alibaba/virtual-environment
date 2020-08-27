# 移除KtEnv

KtEnv的卸载包括移除各个Namespace中的Operator对象以及移除全局的CRD和Webhook组件。

首先，进入KtEnv的部署文件包目录（见[部署文档](zh-cn/doc/deployment.md?id=部署KtEnv组件)）。

在使用了虚拟环境功能的Namespace里移除Operator，以`default` Namespace为例：

```bash
kubectl label namespace default environment-tag-injection-
kubectl delete -n default -f ktenv_operator.yaml
kubectl delete -n default -f ktenv_service_account.yaml
```

最后移除全局的KtEnv组件：

```bash
kubectl delete -f global/ktenv_crd.yaml
kubectl delete -f global/ktenv_webhook.yaml
```
