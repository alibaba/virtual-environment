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

其余以此类推。

### 验证路由

```bash
# 首先创建一个在集群中的容器用于访问服务
kubectl create deployment sleep --image=virtualenvironment/sleep
# 启动演示的服务实例
cd examples/deploy/
./app.sh apply

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
  [node @ dev] <-
  
# 清理演示使用的服务实例
./app.sh delete
kubectl delete deployment sleep
```
