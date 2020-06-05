# 部署虚拟环境

**前提**：集群已经启用Istio，本地已安装并配置kubectl。请参考：

- Istio：https://istio.io/docs/setup/install
- kubectl：https://kubernetes.io/docs/tasks/tools/install-kubectl

## 部署到集群

从 [发布页面](https://github.com/alibaba/virtual-environment/releases) 下载最新的CRD文件，并解压。

```bash
wget https://github.com/alibaba/virtual-environment/releases/download/v0.3/kt-virtual-environment-v0.3.zip
unzip kt-virtual-environment-v0.3.zip
cd v0.3/
```

使用`kubectl apply`命令将解压后目录中的CRD和Webhook配置应用到Kubernetes，其中Webhook携带了默认的自签名秘钥，可参考[Webhook配置文档](zh-cn/ve/webhook.md)替换。

```bash
kubectl apply -f crds/env.alibaba.com_virtualenvironments_crd.yaml
kubectl apply -f webhooks/virtualenvironment_tag_injector_webhook.yaml
```

将Operator部署到每个需要使用虚拟环境的目标Namespace里，比如`default`。

```bash
kubectl apply -n default -f operator.yaml
```

如果集群开启了RBAC，还需要部署相应的Role和ServiceAccount。

```bash
kubectl apply -n default -f service_account.yaml
kubectl apply -n default -f role.yaml
kubectl apply -n default -f role_binding.yaml
```

现在，Kubernetes集群就已经具备使用虚拟环境能力了。

## 创建虚拟环境

创建类型为`VirtualEnvironment`的资源定义文件，使用`kubectl apply`命令添加到Kubernetes集群的目标Namespace中

```bash
kubectl apply -n default -f path-to-virtual-environment-cr.yaml
```

实例创建后，会自动监听**所在Namespace中的**所有Service和Deployment对象并自动生成路由隔离规则，形成虚拟环境。

资源定义文件内容请参考[配置虚拟环境](zh-cn/ve/configuration.md)，根据实际情况修改配置参数。

## 应用程序适配

根据虚拟环境配置，为Pod添加虚拟环境标签，并让服务在调用链上透传虚拟环境标签。

- 为Deployment定义的Pod模板增加标识虚拟环境名称的Label（默认约定的标签键为`virtual-env`）

- 为应用程序添加透传标签Header的功能（默认约定的请求头键为`X-Virtual-Env`）

完整示例请参考[快速开始](zh-cn/ve/quickstart.md)。
