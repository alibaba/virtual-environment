# 部署虚拟环境

**前提**：集群已经启用Istio，本地已安装并配置kubectl。请参考：

- Istio：https://istio.io/docs/setup/install
- kubectl：https://kubernetes.io/docs/tasks/tools/install-kubectl

## 部署到集群

使用`kubectl apply`命令部署Operator到Kubernetes

```bash
kubectl apply -f https://virtual-environment.oss-cn-zhangjiakou.aliyuncs.com/release/0.1/env.alibaba.com_virtualenvironments_crd.yaml
kubectl apply -f https://virtual-environment.oss-cn-zhangjiakou.aliyuncs.com/release/0.1/operator.yaml
```

如果集群开启了RBAC，还需要部署相应的Role和ServiceAccount

```bash
kubectl apply -f https://virtual-environment.oss-cn-zhangjiakou.aliyuncs.com/release/0.1/service_account.yaml
kubectl apply -f https://virtual-environment.oss-cn-zhangjiakou.aliyuncs.com/release/0.1/role.yaml
kubectl apply -f https://virtual-environment.oss-cn-zhangjiakou.aliyuncs.com/release/0.1/role_binding.yaml
```

现在，Kubernetes集群就已经具备使用虚拟环境能力了。

## 创建虚拟环境

创建类型为`VirtualEnvironment`的资源定义文件，使用`kubectl apply`命令添加到Kubernetes集群

```bash
kubectl apply -f path-to-virtual-environment-cr.yaml
```

实例创建后，会自动监听**所在Namespace中的**所有Service和Deployment对象并自动生成路由隔离规则，形成虚拟环境。

资源定义文件内容请参考[配置虚拟环境](configuration.md)，根据实际情况修改配置参数。

## 应用程序适配

根据虚拟环境配置，为Pod添加虚拟环境标签，并让服务在调用链上透传虚拟环境标签。

- 为Deployment定义的Pod模板增加标识虚拟环境名称的Label（默认约定的标签键为`virtual-env`）

- 为应用程序添加透传标签Header的功能（默认约定的请求头键为`X-Virtual-Env`）

完整示例请参考[快速开始](quickstart.md)。
