# 使用SDK

虚拟环境采用Service Mesh机制完成服务之间的自动路由隔离规则，此规则依赖于请求所携带的"环境标签"Header，这就要求此Header的值需要从最初的请求发出后，沿着调用链被全程保留下来。但实际情况是，对于大多数的服务而言，默认情况下子级调用并不会自动包含父级调用的上下文Header，这与许多链路追踪工具所面临的问题是一致的。

在项目的[examples目录](https://github.com/alibaba/virtual-environment/tree/master/examples)里，我们提供了Java、Nodejs、Golang三种语言的代码示例，这些示例中都包含了显式的从前级调用取出"环境标签"Header值，然后设置到下级调用Header的操作，这个操作的过程比较重复且易被忘记，因此在诸如Zipkin、SkyWalking等链路追踪工具中，通常会提供语言相关的SDK来辅助用户自动完成这些行为。采用同样的原理，我们也可以通过SDK的方式来简化"环境标签"在调用链上透传的过程。

在项目的[sdk目录](https://github.com/alibaba/virtual-environment/tree/master/sdk/java)里，已经提供了一个Java语言的SDK示例，采用Spring框架的切面机制来自动化"环境标签"的传递。需要注意的是，这个SDK仅仅作为一个参考思路的示例，因此并未被上传到任何公共的Maven仓库中，也无法在其他项目里直接依赖使用，用户可根据实际情况对此示例SDK进行修改并上传到自己的私有依赖仓库中使用。
