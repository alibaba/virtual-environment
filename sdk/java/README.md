Java SDK 示例
---

使用SDK方式在上下游HTTP请求之间自动透传环境标签的演示代码。

子目录说明：
- `lib` 环境标签透传的SDK代码
- `demo` 用于展示SDK使用方法的简单示例

## 使用方法

您可以直接通过Maven或Gradle从官方仓库直接使用，或自行将`lib`中的代码修改为符合您实际情况的依赖包。

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

添加SDK依赖后，
1. 代码中所有通过`@Controller`（或`@RestController`）注解方式创建的API将自动从指定Header读取环境标签，并存储到请求的上下文
2. 当收到的请求不包含指定的环境标签Header时，将检查当前服务实例运行时是否包含指示当前虚拟环境的环境变量，如果有则读取并存储到请求上下文
3. 代码中所有通过Spring容器管理的`RestTemplate`在发出请求前将自动读取所在上下文的环境标签，并添加到发出请求的指定Header里

> 默认用于存储环境标签的Header是`X-Virtual-Env`，用于标记当前应用实例所属虚拟环境的变量为`APP_VIRTUAL_ENV`，可通过调用`InterceptorGlobal.setupInterceptors()`方法修改

## 运行示例

TBD