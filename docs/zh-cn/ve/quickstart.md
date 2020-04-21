# 快速开始

示例代码和部署结构见[Examples目录](https://github.com/alibaba/virtual-environment/tree/master/examples)说明。

## 获取实例代码

拉取Git仓库，进入示例代码目录

```bash
git clone https://github.com/alibaba/virtual-environment.git
cd virtual-environment/examples
```

## 验证虚拟环境隔离

在不使用KtConnect工具的情况下，隔离仅对集群内的Pod间调用有效。为了便于观察，可在集群随意创建一个临时的Pod作为发送测试请求的容器。

若当前Namespace中已经存在其他可用的Pod，这个步骤可以忽略。

```bash
# 创建一个在集群中的容器用于访问服务
kubectl create deployment sleep --image=virtualenvironment/sleep --dry-run -o yaml \
        | istioctl kube-inject -f - | kubectl apply -n default -f -
```

使用`app.sh`脚本快速创建示例所需的VirtualEnvironment、Service和Deployment资源。

```bash
# 启动演示的服务实例
deploy/app.sh apply default
```

依次使用`kubectl get virtualenvironment`、`kubectl get service`、`kubectl get deployment`命令查看各资源的创建情况，等待所有资源部署完成。

进入同Namespace的任意一个Pod，例如前面步骤创建的sleep容器。

```bash
# 进入集群中的容器
kubectl exec -n default -it $(kubectl get -n default pod -l app=sleep -o jsonpath='{.items[0].metadata.name}') /bin/sh
```

分别在请求头加上不同的虚拟环境名称，使用`curl`工具调用`app-js`服务。注意该示例创建的VirtualEnvironment实例配置使用`-`作为环境层级分隔符，同时配置了传递标签Header的键名为`ali-env-mark`。

已知各服务输出文本结构为`[项目名 @ 响应的Pod所属虚拟环境] <- 请求标签上的虚拟环境名称`。观察实际响应的服务实例情况：

```bash
# 使用dev.proj1标签
> curl -H 'ali-env-mark: dev.proj1' app-js:8080/demo
  [springboot @ dev.proj1] <-dev.proj1
  [go @ dev] <-dev.proj1
  [node @ dev.proj1] <-dev.proj1

# 使用dev.proj1.feature1标签
> curl -H 'ali-env-mark: dev.proj1.feature1' app-js:8080/demo
  [springboot @ dev.proj1.feature1] <-dev.proj1.feature1
  [go @ dev] <-dev.proj1.feature1
  [node @ dev.proj1] <-dev.proj1.feature1

# 使用dev.proj2标签
> curl -H 'ali-env-mark: dev.proj2' app-js:8080/demo
  [springboot @ dev] <-dev.proj2
  [go @ dev.proj2] <-dev.proj2
  [node @ dev] <-dev.proj2

# 不带任何标签访问
# 由于启用了AutoInject配置，经过node服务后，请求自动加上了Pod所在虚拟环境的标签
> curl app-js:8080/demo
  [springboot @ dev] <-dev
  [go @ dev] <-dev
  [node @ dev] <-empty
```

## 本地服务加入隔离域

通过[KtConnect](https://github.com/alibaba/kt-connect)工具可将本地网络与集群打通，并将本地服务直接加入任意指定隔离域。

- 用`ktctl connect`实现本地直接访问集群

```bash
# 注意label参数指定了加入的隔离域名称
sudo ktctl --label virtual-env=dev.proj2 --namespace default connect
```

现在无需进入集群中的容器，本地就可以直接访问`app-js`服务了

```bash
# 由于开启了envHeader.autoInject配置，本地发出的请求会被自动加上隔离域Header
$ curl app-js:8080/demo
  [springboot @ dev] <-dev.proj2
  [go @ dev.proj2] <-dev.proj2
  [node @ dev] <-dev.proj2
```

- 用`ktctl mesh`实现本地服务加入隔离域

```bash
# 本地启动一个app-java服务实例，监听8080端口
# 注意此处设置envMark环境变量只是为了让本地服务输出隔离域名称，与实际路由控制无关
cd examples/springboot
envMark=local mvn spring-boot:run

# label参数指定本地服务加入的隔离域
# app-java-dev是app-java服务在集群公共隔离域的Deployment实例名
sudo ktctl --label virtual-env=dev.proj2 --namespace default mesh app-java-dev --expose 8080
```

此时本地`app-java`服务在`dev.proj2`隔离域，因而新的访问途径应当为：


```
            +----------+
dev         |  app-js  |
            +----------+
                           +----------+   +-----------------+
dev.proj2                  |  app-go  |   | app-java(local) |
                           +----------+   +-----------------+
```

再次调用，本地发请求经过集群的app-js和app-go服务之后，路由到了本地的app-java实例

```bash
$ curl app-js:8080/demo
  [springboot @ local] <-dev.proj2
  [go @ dev.proj2] <-dev.proj2
  [node @ dev] <-dev.proj2
```

## 清理示例资源

```bash
# 删除演示使用的服务实例
deploy/app.sh delete default
# 删除用于发起访问请求的临时容器
kubectl delete -n default deployment sleep
```
