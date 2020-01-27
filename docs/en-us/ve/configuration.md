# Configuration guide

Virtual environment instances are defined through Kubernetes resources of `VirtualEnvironment` kind. This is an example of its content structure (from [env.alibaba.com_v1alpha1_virtualenvironment_cr.yaml](https://github.com/alibaba/virtual-environment/blob/master/deploy/crds/env.alibaba.com_v1alpha1_virtualenvironment_cr.yaml))

```yaml
apiVersion: env.alibaba.com/v1alpha1
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

The parameter are shown in below table:

| Config item            | Default value | Description  |
| :--------              | :-----:       | :---- |
| envHeader.name         | X-Virtual-Env | Name of header to keep env name in trace (recommend to set expressly) |
| envHeader.autoInject   | false         | Whether auto inject env header via sidecar (recommend to enable) |
| envLabel.name          | virtual-env   | Name of pod label to mark virtual environment name (recommend to leave as default) |
| envLabel.splitter      | .             | Symbol to split virtual environment levels (single symbol only) |
| envLabel.defaultSubset |               | Default subset to route when env header matches nothing (default means random) |

A VirtualEnvironment instance is only valid for the namespace it is in. If needed, you can create instances with the same configuration in multiple namespaces to achieve [isolation across namespaces and clusters](cross-cluster.md).
