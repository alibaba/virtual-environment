# 问题排查

- 说明：以下示例中，假设虚拟环境实例所在的Namespace已经存在环境变量`$NS`中

在排查各类问题前，请按照[检查部署结果](zh-cn/doc/deployment.md?id=检查部署结果)小节验证KtEnv的CRD和Webhook组件已正确安装到集群中。

## 路由规则不符合预期

首先检查目标Namespace中是否正确创建了VirtualEnvironment实例：

```bash
kubectl -n $NS get VirtualEnvironment
```

然后检查是否生成了预期的Istio资源。对于每一个选中包含路由标签Pod的Service实例，应该生成一个同名的VirtualService实例和一个同名的DestinationRule实例。虽然没有命令能够快速找出这些Service，但通过标签筛选可以找到所有符合要求的Pod：

```bash
# 假设路由标签名是virtual-env（这个名称是在创建VirtualEnvironment实例时候候配置的）
kubectl -n $NS get Pod -l virtual-env
```

通过这些Pod试着列举一下有哪些应该参与路由隔离的Service，然后与Namespace里的Istio资源进行比较：

```bash
# 这两种资源的数目应该相同，且与参与路由隔离的Service逐一同名对应
kubectl -n $NS get VirtualService
kubectl -n $NS get DestinationRule
```

若数目不正确，请检查目标Service对象的端口命名：

```bash
kubectl -n $NS get Service <要路由的目标服务名> -o jsonpath='{.spec.ports}'
```

端口名称必须依据[Istio文档](https://istio.io/latest/docs/ops/configuration/traffic-management/protocol-selection/)要求采用`<协议>[-<后缀>]`结构。由于当前Istio仅支持对`HTTP`协议的消息进行精细路由控制，因此KtEnv仅会处理名称以`http`开头的端口。

如果端口命名没有问题，但实例数目依然不正常，可检查VirtualEnvironment的实例配置和运行日志，通常是配置不正确或生成Istio资源时候出错了：

```bash
# 先看VirtualEnvironment实例的日志，留意其中的错误信息，Ctrl+C结束
kubectl logs -n $NS $(kubectl get pod -l name=virtual-env-operator -o jsonpath='{.items[0].metadata.name}' -n $NS) virtual-env-operator --tail 50 --follow
# 若没有在日志中发现可以信息，请认真检查VirtualEnvironment配置是否符合实际情况
kubectl -n $NS get VirtualEnvironment -o yaml
```

如果Istio资源数目正常，说明路由规则已生成（可以用`kubectl -n $NS get VirtualService <服务名> -o yaml`查看具体规则，这里通常不会有问题），接下来可检查Pod的Envoy Sidecar日志：

```bash
kubectl -n $NS logs <任意一个Pod名字> istio-proxy --tail 100
```

若路由配置正常，Sidecar运行也无任何错误，则需结合Istio本身功能进一步排查原因。

以下是几种比较常见的错误原因：

- 同一个Pod被多个Service选中。当前Istio不支持一个Pod同时属于多个Service的情况
- Istio规则生效有延迟（参考 [Istio文档](https://istio.io/latest/zh/docs/ops/common-problems/network-issues/#route-rules-don't-seem-to-affect-traffic-flow) ）

## 流量未自动加环境标

流量自动加标的过程分为两步，首先在Pod创建时通过全局Webhook组件将记录在Pod label上的环境名称写入其Sidecar容器的`VIRTUAL_ENVIRONMENT_TAG`环境变量，然后在流量出口处通过Envoy Sidecar读取上下文环境变量的内容，将环境标最终写到HTTP请求的Header里。

首先检查Webhook是否成功的将环境标写入Pod环境变量：

```bash
kubectl -n $NS get pod <任意一个Pod名字> -o yaml -o yaml | grep -A 1 'VIRTUAL_ENVIRONMENT_TAG'
```

如果没有输出任何内容，说明环境变量未注入，请检查Admission Webhook组件是否正确部署。

若有输出Pod所处的环境标名称，则问题出在Envoy Sidecar上。

接下来查看Envoy容器日志，若是注入脚本出错，这里会看到报错信息：

```bash
kubectl -n $NS logs <任意一个Pod名字> istio-proxy --tail 100
```

同时检查生成的EnvoyFilter对象：

```bash
kubectl -n $NS get EnvoyFilter <与VirtualEnvironment实例同名> -o yaml
```

若该对象存在且内容正常，可导出Envoy的配置，检查注入脚本是否正确添加：

```bash
kubectl -n $NS exec <任意一个Pod名字> -c istio-proxy curl http://localhost:15000/config_dump | less
```

搜索`virtual.environment.lua`文本，其上下文位置应该在`configs.dynamic_listeners.active_state`区块内，若是出现在`configs.dynamic_listeners.error_state`区块，请检查是否与其他Operator生成的路由规则存在冲突。

## Sidecar容器始终未就绪

查看Envoy容器日志，若发现如下内容：

```text
Envoy proxy is NOT ready: config not received from Pilot
```

可先登录到Envoy容器中：

```bash
kubectl -n $NS exec -it <任意一个Pod名字> -c istio-proxy /bin/bash
```

然后执行`nc -zvw2 istio-pilot.istio-system 15010`，
正常情况应当返回内容类似`istio-pilot.istio-system.svc.cluster.local [172.21.7.52] 15010 (?) open`。
异常情况可能返回“Unknown host”、“Connection refused”、“Connection timeout”等。
然后根据情况进一步排查问题原因。

参考文章[二分之一活的微服务](https://juejin.im/post/5ecdf080e51d457841190d22)
