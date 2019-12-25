Virtual Environment Operator
===========

Isolate kubernetes pod communication into virtual groups according to pod label, with route fallback support.

[中文介绍](./README_CN.md)

## Instruction

This project inspired by the "project environment" practice widely used in Alibaba Group,
witch allow developers create virtual testing environment cheaply and quickly by reusing many service instances from
another testing environment.

```
           +----------+   +----------+   +----------+
ParentEnv  | ServiceA |   | ServiceB |   | ServiceC |
           +----------+   +-^------+-+   +----------+
                            |      |
           +-----------+    |      |    +-----------+
SubEnv    -> ServiceA' +----+      +----> ServiceC' |
           +-----------+                +-----------+
```

This is a simple example of 2 virtual environments. A ParentEnv contains instances of all the 3 services,
while a SubEnv inherited from ParentEnv only contains instance of ServiceA and ServiceC (record as `ServiceA'` and `ServiceC'`).

The arrow shows a call sequence "ServiceA->ServiceB->ServiceC", when the request comes from instance `ServiceA'`,
it should firstly fallback to hit `ServiceB` in the ParentEnv, and finally turn to `ServiceC'` in SubEnv.

## How it works

In this implementation, a pod label is used to identify the name and inherit relationship of pods.
Service mesh (currently only support Istio) rules will be dynamically generated to fit the isolation requirement.

For example, the identify label key is set as `virtual-env`, and symbol `.` is used to split inherit level. Then,

- All pods contain label `virtual-env: dev` and `virtual-env: dev.proj1` would logically belong to two virtual environments named `dev` and `dev.proj1`
- In any condition, when traffic from `dev.proj1` virtual environment visit a service which has no pod available in that environment, pod from the upper level (i.e. `dev`) virtual environment would response
- Even the traffic have been passed to a upper level pod, in the following sequence, if target service has pods in the original environment (i.e. `dev.proj1`), those pods should respond for the call
- If the virtual environment label contains multiple inherit symbol, e.g. `dev.proj1.user1`, then the fallback logic should goes up through the hierarchy, until find matched pod or reach the most top level

## Limitation

This project currently only support [Istio](http://istio.io) service mesh control plane, the limitations are mostly come with Istio.

- Only support HTTP protocol
- Only support service discovery via kubernetes internal DNS, framework with external service discovery (e.g. Dubbo, SpringCloud) would not work
- The application MUST pass down a HTTP header contains original traffic environment name by itself. This can be done by OpenTracing's baggage, or implement manually

## Usage

**[0]** Premise: Kubernetes cluster with Istio installed

**[1]** Deploy CRD and Operator

```bash
kubectl apply -f deploy/crds/env.alibaba.com_virtualenvironments_crd.yaml
kubectl apply -f deploy/operator.yaml
```

If the cluster has RBAC enabled, please also apply Role and ServiceAccount

```bash
kubectl apply -f deploy/service_account.yaml
kubectl apply -f deploy/role.yaml
kubectl apply -f deploy/role_binding.yaml
```

**[2]** Add mechanism for passing down environment name header in application (the default header key is `X-Virtual-Env`)

**[3]** Build docker image, when deploy to Kubernetes put a pod label to identify which virtual environment it belongs to (the default label key is `virtual-env`)

**[4]** Create a resource with kind `VirtualEnvironment` (e.g. `deploy/crds/env.alibaba.com_v1alpha1_virtualenvironment_cr.yaml`), change the spec values, then use `kubectl apply` to take effect

## VirtualEnvironment configuration

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

| Config item            | Default value | Description  |
| :--------              | :-----:       | :---- |
| envHeader.name         | X-Virtual-Env | Name of header to keep env name in trace (recommend to set expressly) |
| envHeader.autoInject   | false         | Whether auto inject env header via sidecar (recommend to enable) |
| envLabel.name          | virtual-env   | Name of pod label to mark virtual environment name (recommend to leave as default) |
| envLabel.splitter      | .             | Symbol to split virtual environment levels (single symbol only) |
| envLabel.defaultSubset |               | Default subset to route when env header matches nothing (default means random) |

## Support

Contact us with DingTalk:

![image](https://github.com/alibaba/kt-connect/raw/master/docs/_media/dingtalk-group.png)
