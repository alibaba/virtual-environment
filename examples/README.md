# 演示

### 环境拓扑结构

```
                    +-----------+   +-----------+   +-----------+
dev                 |  demo-js  |   |  demo-go  |   | demo-java |
                    +-----------+   +-----------+   +-----------+

                    +-----------+                   +-----------+
dev-proj1           |  demo-js  |                   | demo-java |
                    +-----------+                   +-----------+

                                                    +-----------+
dev-proj1-feature1                                  | demo-java |
                                                    +-----------+

                                    +-----------+
dev-proj2                           |  demo-go  |
                                    +-----------+
```

### 预期结果

测试`demo-js`->`demo-go`->`demo-java`调用路径。

- 来源请求包含`dev-proj1`头标签。
由于`dev-proj1`虚拟环境中只存在`demo-js`和`demo-java`服务，因此访问途径为：

```
                                    +-----------+                
dev                                 |  demo-go  |                
                                    +-----------+                
                    +-----------+                   +-----------+
dev-proj1           |  demo-js  |                   | demo-java |
                    +-----------+                   +-----------+
```

- 来源请求包含`dev-proj1-feature1`头标签。
由于`dev-proj1-feature1`虚拟环境中只存在`demo-java`服务，因此访问途径为：

```
                                    +-----------+                
dev                                 |  demo-go  |                
                                    +-----------+                
                    +-----------+                                
dev-proj1           |  demo-js  |                                
                    +-----------+                                
                                                    +-----------+
dev-proj1-feature1                                  | demo-java |
                                                    +-----------+
```

- 来源请求包含`dev-proj2`头标签。
由于`dev-proj2`虚拟环境中只存在`demo-go`服务，因此访问途径为：

```
                    +-----------+                   +-----------+
dev                 |  demo-js  |                   | demo-java |
                    +-----------+                   +-----------+
                                    +-----------+
dev-proj2                           |  demo-go  |
                                    +-----------+
```

其余情况以此类推。

### 验证路由

```bash
# 首先创建一个在集群中的容器用于访问服务
kubectl create deployment sleep --image=virtualenvironment/sleep --dry-run -o yaml \
        | istioctl kube-inject -f - | kubectl apply -f -
# 启动演示的服务实例
deploy/app.sh apply

# 进入集群中的容器
kubectl exec -it $(kubectl get pod -l app=sleep -o jsonpath='{.items[0].metadata.name}') /bin/sh
# 使用dev-proj1标签
> curl -H 'ali-env-mark: dev-proj1' app-js:8080/demo
  [springboot @ dev-proj1] <-dev-proj1
  [go @ dev] <-dev-proj1
  [node @ dev-proj1] <-dev-proj1
# 使用dev-proj1-feature1标签
> curl -H 'ali-env-mark: dev-proj1-feature1' app-js:8080/demo
  [springboot @ dev-proj1-feature1] <-dev-proj1-feature1
  [go @ dev] <-dev-proj1-feature1
  [node @ dev-proj1] <-dev-proj1-feature1
# 使用dev-proj2标签
> curl -H 'ali-env-mark: dev-proj2' app-js:8080/demo
  [springboot @ dev] <-dev-proj2
  [go @ dev-proj2] <-dev-proj2
  [node @ dev] <-dev-proj2
# 不带任何标签访问
> curl app-js:8080/demo
  [springboot @ dev] <-dev
  [go @ dev] <-dev
  [node @ dev] <-empty
  
# 清理演示使用的服务实例
deploy/app.sh delete
kubectl delete deployment sleep
```

### 本地服务加入隔离域

通过[kt-connect](https://github.com/alibaba/kt-connect)可将本地网络与集群打通，并将本地服务直接加入任意指定隔离域。

- 用`ktctl connect`实现本地直接访问集群

```bash
# 注意label参数指定了加入的隔离域名称
sudo ktctl --label virtual-env=dev-proj2 --namespace default connect
```

现在无需进入集群中的容器，本地就可以直接访问`app-js`服务了

```bash
# 由于开启了envHeader.autoInject配置，本地发出的请求会被自动加上隔离域Header
$ curl app-js:8080/demo
  [springboot @ dev] <-dev-proj2
  [go @ dev-proj2] <-dev-proj2
  [node @ dev] <-dev-proj2
```

- 用`ktctl mesh`实现本地服务加入隔离域

```bash
# 本地启动一个app-java服务实例，监听8080端口
# 注意此处设置envMark环境变量只是为了让本地服务输出隔离域名称，与实际路由控制无关
cd examples/springboot
envMark=local mvn spring-boot:run

# label参数指定本地服务加入的隔离域
# app-java-dev是app-java服务在集群公共隔离域的Deployment实例名
sudo ktctl --label virtual-env=dev-proj2 --namespace default mesh app-java-dev --expose 8080
```

此时本地`app-java`服务在`dev-proj2`隔离域，因而新的访问途径应当为：


```
            +-----------+
dev         |  demo-js  |
            +-----------+
                            +-----------+   +------------------+
dev-proj2                   |  demo-go  |   | demo-java(local) |
                            +-----------+   +------------------+
```

再次调用，本地发请求经过集群的app-js和app-go服务之后，路由到了本地的app-java实例

```bash
$ curl app-js:8080/demo
  [springboot @ local] <-dev-proj2
  [go @ dev-proj2] <-dev-proj2
  [node @ dev] <-dev-proj2
```
