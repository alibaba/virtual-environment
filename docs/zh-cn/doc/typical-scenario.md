# 典型场景

本质而言，环境隔离机制的核心在于识别染色流量的`环境标`，其用途并不局限于层级化的测试环境管理，以下列举几种VirtualEnvironment项目可以使用的场景。

#### 1. 单服务调试

个人开发者使用时，结合kt-connect工具和浏览器的[ModHeader插件](https://github.com/bewisse/modheader)，可将本地开发中的服务混入集群里，同时保证自己打开浏览器访问的流量始终流经本地的服务进程，不会进入集群里存在的其他同类Pod实例，而其他开发者的正常调用不会进入该本地实例。

<img src="https://virtual-environment.oss-cn-zhangjiakou.aliyuncs.com/image/typical-scenario-1.jpg" height="300px"/>

#### 2. 多服务联调

在团队开发的场景下，结合kt-connect工具可使多个开发者共享相同的流量标签，形成可以相互调用的联调小圈子，个人还可以基于这个标签再创建子级标签，用来本地单步调试，同时不会影响其他未进入该虚拟环境的开发者的正常流量。为此我们总结了一种更普适的[最佳实践模式](https://alibaba.github.io/virtual-environment/#/zh-cn/doc/best-practice)。

<img src="https://virtual-environment.oss-cn-zhangjiakou.aliyuncs.com/image/typical-scenario-2.jpg" height="300px"/>

#### 3. 局部替换服务版本

在进行功能验证时，当遇到某些公共环境部署了有缺陷或是破坏兼容性的版本，而导致自己开发分支上的代码无法正常运行时，可以用服务的主干版本部署一个临时实例放到隔离环境中，避免把更多的时间浪费在等待修复和提前合并其他尚未进入主干的修改。这个场景可以融入到前两个场景中结合使用。

<img src="https://virtual-environment.oss-cn-zhangjiakou.aliyuncs.com/image/typical-scenario-3.jpg" height="300px"/>

#### 4. 集成测试链路隔离

虚拟环境也可以成为自动化测试的好帮手。通过隔离出要进行集成测试的特定版本实例，无需额外的计算资源就能搞定待测服务的运行依赖问题（复用公共环境实例）。

<img src="https://virtual-environment.oss-cn-zhangjiakou.aliyuncs.com/image/typical-scenario-4.jpg" height="200px"/>

#### 5. 快速多版本对比

利用ModHeader插件给浏览器访问请求任意“上色”的方法，还可以将虚拟环境与AB灰度测试结合，快速切换请求链路流经某个后台服务新旧版本，对比端到端的运行效果。

<img src="https://virtual-environment.oss-cn-zhangjiakou.aliyuncs.com/image/typical-scenario-5.jpg" height="300px"/>
