# 配置虚拟环境

虚拟环境实例通过`VirtualEnvironment`类型的Kubernetes资源定义。其内容结构示例如下（来自[env.alibaba.com_v1alpha1_virtualenvironment_cr.yaml](https://github.com/alibaba/virtual-environment/blob/master/deploy/crds/env.alibaba.com_v1alpha1_virtualenvironment_cr.yaml)）：

```yaml
apiVersion: env.alibaba.com/v1alpha2
kind: VirtualEnvironment
metadata:
  name: example-virtualenv
spec:
  envHeader:
    name: X-Virtual-Env
    autoInject: true
  envLabel:
    name: virtual-env
    splitter: .
    defaultSubset: dev
```

参数作用如表所示：

| 配置参数                | 默认值         | 说明  |
| :--------              | :-----:       | :---- |
| envHeader.name         | X-Virtual-Env | 用于记录虚拟环境名的HTTP头名称（虽然有默认值，建议显性设置） |
| envHeader.autoInject   | false         | 是否为没有虚拟环境HTTP头记录的请求自动注入HTTP头（建议开启） |
| envLabel.name          | virtual-env   | Pod上标记虚拟环境名用的标签名称（除非确实必要，建议保留默认值） |
| envLabel.splitter      | .             | 虚拟环境名中用于划分环境默认路由层级的字符（只能是单个字符） |
| envLabel.defaultSubset |               | 请求未匹配到任何存在的虚拟环境时，进行兜底虚拟环境名（默认为随机路由） |

**注意**：VirtualEnvironment实例只对其所在的Namespace有效。如有需要，可以通过在多个Namespace分别创建相同配置的实例，实现[跨Namespace和跨集群的隔离](zh-cn/ve/cross-cluster.md)。
