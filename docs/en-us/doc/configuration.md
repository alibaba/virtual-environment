# Configuration guide

Virtual environment instances are defined through Kubernetes resources of `VirtualEnvironment` kind. This is an example of its content structure (from [env.alibaba.com_v1alpha1_virtualenvironment_cr.yaml](https://github.com/alibaba/virtual-environment/blob/master/deploy/crds/env.alibaba.com_v1alpha1_virtualenvironment_cr.yaml))

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

The parameters are described in below table:

| Config item            | Default value | Description  |
| :--------              | :---          | :--- |
| envHeader.name         | X-Virtual-Env | Name of header to keep env name in trace (recommend to set expressly) |
| envHeader.autoInject   | false         | Whether auto inject env header via sidecar (recommend to enable) |
| envHeader.aliases      |               | Suppletive header names used for keeping env name (usually not necessary) |
| envLabel.name          | virtual-env   | Name of pod label to mark virtual environment name (recommend to leave as default) |
| envLabel.splitter      | .             | Symbol to split virtual environment levels (single symbol only) |
| envLabel.defaultSubset |               | Default subset to route when env header matches nothing (default means random) |

The `envHeader.aliases` configuration mainly used when extracting part of the text from the header content as the environment name required, e.g. [Coloring of callback traffic](https://github.com/alibaba/virtual-environment/issues/14).
This feature will significantly increase the amount of Istio rules generated, please do not use it if not necessary.

The parameters are as follows:

| Config item                   | Default value | Description  |
| :--------                     | :---          | :--- |
| envHeader.aliases.name        |               | Name of header to keep env name in trace (required) |
| envHeader.aliases.pattern     |               | Use regular expressions to match the environment name of the HTTP header content, if it is empty, it means use full text matching |
| envHeader.aliases.placeholder | (@)           | Placeholder for environment name in regular expression |

The value of `envHeader.aliases.placeholder` should exist inside `envHeader.aliases.pattern` text, this placeholder text will be replaced with environment name when generating route rules.

**Notice**: A VirtualEnvironment instance is only valid for the namespace it is in. If needed, you can create instances with the same configuration in multiple namespaces to achieve [isolation across namespaces and clusters](en-us/doc/cross-cluster.md).
