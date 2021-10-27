Java SDK 示例
---

使用SDK方式在上下游HTTP请求之间自动透传环境标签的演示代码。

您可以直接通过Maven或Gradle从官方仓库直接使用，或自行将代码修改为符合您使用情况的工具包。

Maven依赖：

```xml
<dependency>
    <groupId>com.alibaba.kt</groupId>
    <artifactId>trace-sdk</artifactId>
    <version>0.1.0</version>
</dependency>
```

Gradle依赖：

```groovy
dependencies {
    implementation('com.alibaba.kt:trace-sdk:0.1.0')
}
```

添加此依赖后，所有通过`@Controller`（或`@RestController`）注解方式创建的API，`RestTemplate`

子目录说明：
- `lib`
- `demo`
