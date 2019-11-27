演示
---

环境拓扑结构
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


```bash
# 首先创建一个在集群中的容器用于访问服务
kubectl create deployment sleep --image=virtualenvironment/sleep
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
```
