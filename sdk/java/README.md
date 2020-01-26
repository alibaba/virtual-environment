# HTTP头参数透传插件

## 简介

通过HTTP头透传数据的插件演示，依赖Spring 3.0以后引入的`RestTemplate`类型（Spring Cloud框架默认支持）。

## 使用方法

### 1.基于Spring Boot的应用

采用Spring Boot构建的应用只需加入插件依赖即可，以Maven项目的pom.xml文件为例：

```xml
<dependency>
	<groupId>com.alibaba.aone</groupId>
	<artifactId>demax-trace-plugin</artifactId>
	<version>1.0.0-SNAPSHOT</version>
</dependency>
```

### 2.普通Spring MVC应用

若是使用了Spring MVC，但非Spring Boot的应用，首先引入插件依赖，方法同上。
然后手工增加Bean扫描路径，可通过xml或配置类实现。

xml配置方式：

```xml
<beans>  
    <context:component-scan base-package="com.alibaba.aone.demax.trace.plugin"/>  
</beans> 
```

配置类方式:

```java
package <任意包路径>;

import org.springframework.context.annotation.ComponentScan;
import org.springframework.context.annotation.Configuration;

@Configuration
@ComponentScan("com.alibaba.aone.demax.trace.plugin")
public class TraceConfig {
}
```

## 功能验证

创建3个通过HTTP依次串联调用的服务，分别引入`demax-trace-plugin`包，
然后访问第一个服务服务的入口API。例如：

```bash
curl -H "X-TB-ENV-ID: FEATURE-123" "http://127.0.0.1:8081/api/trace"
```

在第三个服务的相应方法中读取`X-TB-ENV-ID`请求头，并打印读取结果。从而验证环境ID能自动沿链路进行传递。

也可以创建大量并发请求，观察不同请求传递数据是否互窜：

```bash
for i in {1..1000}; do
    curl -H "X-TB-ENV-ID: $i" "http://127.0.0.1:8081/api/trace" >>/tmp/q &
done
```

等待所有curl进程结束，检查`/tmp/q`文件内容，请求返回的数值没有重复，从而证实每次数据透传都是独立的。
