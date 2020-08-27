# 部署虚拟环境

**前提**：集群已经启用Istio，本地已安装并配置kubectl。请参考：

- Istio：https://istio.io/docs/setup/install
- kubectl：https://kubernetes.io/docs/tasks/tools/install-kubectl

## 部署KtEnv组件

KtEnv系统包含Operator CRD和Admission Webhook两个组件。Webhook组件用于将Pod的虚拟环境标写入到其Sidecar容器的运行时环境变量内；CRD组件用于创建监听集群服务变化并动态生成路由规则的VirtualEnvironment资源实例。

从 [发布页面](https://github.com/alibaba/virtual-environment/releases) 下载最新的部署文件包，并解压。

```bash
wget https://github.com/alibaba/virtual-environment/releases/download/v0.3.2/kt-virtual-environment-v0.3.2.zip
unzip kt-virtual-environment-v0.3.2.zip
cd v0.3.2/
```

将目录中的CRD和Webhook组件添加到Kubernetes（其中Webhook组件携带了默认的自签名秘钥，可参考[Webhook配置文档](zh-cn/doc/webhook.md)替换）。

```bash
kubectl apply -f global/ktenv_crd.yaml
kubectl apply -f global/ktenv_webhook.yaml
```

## 检查部署结果

Webhook组件默认被部署到名为`kt-virtual-environment`的Namespace中，包含一个Service和一个Deployment对象，以及它们创建的子资源对象，可用以下命令查看：

```bash
kubectl -n kt-virtual-environment get all
```

若输出类似以下信息，则表明KtEnv的Webhook组件已经部署且正常运行。

```
NAME                                  READY   STATUS    RESTARTS   AGE
pod/webhook-server-5dd55c79b5-rf6dl   1/1     Running   0          86s

NAME                     TYPE        CLUSTER-IP     EXTERNAL-IP   PORT(S)   AGE
service/webhook-server   ClusterIP   172.21.0.254   <none>        443/TCP   109s

NAME                             READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/webhook-server   1/1     1            1           109s

NAME                                        DESIRED   CURRENT   READY   AGE
replicaset.apps/webhook-server-5dd55c79b5   1         1         1       86s
```

检查上述输出中各资源对象的`AGE`属性（从创建到现在已经过的时间），可确定该对象是否为刚刚新部署的Webhook组件所创建。

CRD组件会在Kubernetes集群内新增一种名为`VirtualEnvironment`的资源类型，在下一步我们将会用到它。可通过以下命令验证其安装状态：

```bash
kubectl get crd virtualenvironments.env.alibaba.com
```

若输出类似以下信息，则表明KtEnv的CRD组件已经正确部署。

```
NAME                                  CREATED AT
virtualenvironments.env.alibaba.com   2020-04-21T13:20:35Z
```

检查输出中`CREATED AT`属性（资源创建时间），可确定该对象是否为刚刚新部署的CRD组件。

## 部署KtEnv Operator

Operator是由CRD组件定义的虚拟环境管理器实例，需要在**每个**使用虚拟环境的Namespace里单独部署。同时为了让Webhook组件对目标Namespace起作用，还应该为其添加值为`enabled`的`environment-tag-injection`标签。

以使用`default` Namespace为例，通过以下命令完成部署。

```bash
kubectl apply -n default -f ktenv_operator.yaml
kubectl label namespace default environment-tag-injection=enabled
```

如果集群开启了RBAC，还需要部署相应的Role和ServiceAccount。

```bash
kubectl apply -n default -f ktenv_service_account.yaml
```

现在，Kubernetes集群就已经具备使用虚拟环境能力了。

## 创建虚拟环境

创建类型为`VirtualEnvironment`的资源定义文件，使用`kubectl apply`命令添加到Kubernetes集群的目标Namespace中

```bash
kubectl apply -n default -f path-to-virtual-environment-cr.yaml
```

实例创建后，会自动监听**所在Namespace中的**所有Service、Deployment和StatefulSet对象并自动生成路由隔离规则，形成虚拟环境。

资源定义文件内容可参考[virtualenv.yaml](https://github.com/alibaba/virtual-environment/blob/master/examples/deploy/virtualenv.yaml)，在[配置虚拟环境](zh-cn/doc/configuration.md)文档中列举了所有可用的配置参数，请根据实际情况进行修改。

## 应用程序适配

根据虚拟环境配置，为Pod添加虚拟环境标签，并让服务在调用链上透传虚拟环境标签。

- 为Deployment定义的Pod模板增加标识虚拟环境名称的Label（默认约定的标签键为`virtual-env`）

- 为应用程序添加透传标签Header的功能（默认约定的请求头键为`X-Virtual-Env`）

完整示例请参考[快速开始](zh-cn/doc/quickstart.md)。
