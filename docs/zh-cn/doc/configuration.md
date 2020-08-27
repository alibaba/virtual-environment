# 配置虚拟环境

虚拟环境实例通过`VirtualEnvironment`类型的Kubernetes资源定义。其内容结构示例如下：

```yaml
apiVersion: env.alibaba.com/v1alpha2
kind: VirtualEnvironment
metadata:
  name: example-virtualenv
spec:
  envHeader:
    name: X-Virtual-Env
    autoInject: true
    aliases:
      - name: ALTERNATIVE-NAME
        pattern: ".*;VirtualEnv=(@);.*"
        placeholder: "(@)"
  envLabel:
    name: virtual-env
    splitter: .
    defaultSubset: dev
```

> 注意：以上内容仅做配置结构参考，请根据实际情况设置各参数的值

参数作用如下表所示：

| 配置参数                | 默认值         | 说明  |
| :--------              | :---          | :--- |
| envHeader.name         | X-Virtual-Env | 用于透传虚拟环境名的HTTP头名称（虽然有默认值，建议显性设置） |
| envHeader.autoInject   | false         | 是否为没有虚拟环境HTTP头记录的请求自动注入HTTP头（建议开启） |
| envHeader.aliases      |               | 添加额外的备选透传HTTP头（通常不需要此配置） |
| envLabel.name          | virtual-env   | Pod上标记虚拟环境名用的标签名称（除非确实必要，建议保留默认值） |
| envLabel.splitter      | .             | 虚拟环境名中用于划分环境默认路由层级的字符（只能是单个字符） |
| envLabel.defaultSubset |               | 请求未匹配到任何存在的虚拟环境时，进行兜底虚拟环境名（默认为随机路由） |

其中`envHeader.aliases`配置主要用于需要从Header内容提取部分文本作为环境名的情况，例如[回调流量的染色](https://github.com/alibaba/virtual-environment/issues/14)。此功能会显著增加生成的Istio规则量，但通常不会造成性能问题。

具体参数作用如下：

| 配置参数                       | 默认值 | 说明 |
| :--------                     | :--- | :--- |
| envHeader.aliases.name        |      | 透传虚拟环境名的HTTP头名称 |
| envHeader.aliases.pattern     |      | 使用正则表达式匹配HTTP头内容的环境名称，若为空则表示使用完整匹配 |
| envHeader.aliases.placeholder | (@)  | 正则表达式中的环境名占位符 |

在`envHeader.aliases.pattern`表达式中应该包含`envHeader.aliases.placeholder`所指定的占位符，该占位符会在生成规则时被替换为对应的环境标签值。

**注意**：VirtualEnvironment实例只对其所在的Namespace有效。如有需要，可以通过在多个Namespace分别创建相同配置的实例，实现[跨Namespace和跨集群的隔离](zh-cn/doc/cross-cluster.md)。

