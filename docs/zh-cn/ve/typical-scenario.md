# 典型场景

#### 1. 单服务调试

开发者将本地服务通过KtConnect连接到集群进行调试， VirtualEnvironment能确保开发者自己始终访问本地实例，而其他开发者的正常调用不会进入该本地实例。

<img src="https://virtual-environment.oss-cn-zhangjiakou.aliyuncs.com/image/typical-scenario-1.jpg" height="300px"/>

#### 2. 多服务联调

多个开发者将本地服务通过KtConnect添加到同一个虚拟环境中， VirtualEnvironment能确保这些开发者之间的调用互通，从而进行项目联调，同时不影响未进入该虚拟环境的其他开发者正常使用测试环境。

<img src="https://virtual-environment.oss-cn-zhangjiakou.aliyuncs.com/image/typical-scenario-2.jpg" height="300px"/>

#### 3. 局部替换服务版本

在进行功能验证时，需要某个服务使用指定的不稳定版本，为了不影响其他开发者使用公共日常环境，可将指定版本部署在隔离的虚拟环境中自己使用

<img src="https://virtual-environment.oss-cn-zhangjiakou.aliyuncs.com/image/typical-scenario-3.jpg" height="300px"/>

#### 4. 集成测试链路隔离

集成测试时，将服务的待测版本放到隔离环境，复用公共日常环境的其他服务实例，从而无需创建全量服务集群，就能快速验证调用链路上特定服务的特定版本功能

<img src="https://virtual-environment.oss-cn-zhangjiakou.aliyuncs.com/image/typical-scenario-4.jpg" height="200px"/>

#### 5. 快速多版本对比

使用浏览器通过插件设置不同的Header值，快速切换访问属于不同虚拟环境中的服务实例，进行前后效果对比。

<img src="https://virtual-environment.oss-cn-zhangjiakou.aliyuncs.com/image/typical-scenario-5.jpg" height="300px"/>
