# 演示

### 目录

- `deploy`: 演示使用的VirtualEnvironment、Service和Deployment资源定义
- `go`: 使用Golang语言编写的简单示例项目
- `node`: 使用Javascript语言编写的简单示例项目
- `springboot`: 使用Java语言编写的简单示例项目

### 环境拓扑结构

下图展示了所有资源部署后，该Namespace中的各虚拟环境拓扑示意（横坐标为服务名，纵坐标为环境标签）

```
                    +----------+   +----------+   +----------+
dev                 |  app-js  |   |  app-go  |   | app-java |
                    +----------+   +----------+   +----------+

                    +----------+                  +----------+
dev.proj1           |  app-js  |                  | app-java |
                    +----------+                  +----------+

                                                  +----------+
dev.proj1.feature1                                | app-java |
                                                  +----------+

                                   +----------+
dev.proj2                          |  app-go  |
                                   +----------+
```

### 预期结果

测试`app-js`->`app-go`->`app-java`调用路径。

- 来源请求包含`dev.proj1`头标签。
由于`dev.proj1`虚拟环境中只存在`app-js`和`app-java`服务，因此访问途径为：

```
                                    +----------+                
dev                                 |  app-go  |                
                                    +----------+                
                    +----------+                  +----------+
dev.proj1           |  app-js  |                  | app-java |
                    +----------+                  +----------+
```

- 来源请求包含`dev.proj1.feature1`头标签。
由于`dev.proj1.feature1`虚拟环境中只存在`app-java`服务，因此访问途径为：

```
                                   +----------+                
dev                                |  app-go  |                
                                   +----------+                
                    +----------+                                
dev.proj1           |  app-js  |                                
                    +----------+                                
                                                  +----------+
dev.proj1.feature1                                | app-java |
                                                  +----------+
```

- 来源请求包含`dev.proj2`头标签。
由于`dev.proj2`虚拟环境中只存在`app-go`服务，因此访问途径为：

```
                    +----------+                  +----------+
dev                 |  app-js  |                  | app-java |
                    +----------+                  +----------+
                                   +----------+
dev.proj2                          |  app-go  |
                                   +----------+
```

其余情况以此类推。

### 功能验证

见[快速开始](https://alibaba.github.io/virtual-environment/#/zh-cn/doc/quickstart)文档。
