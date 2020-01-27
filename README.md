Virtual Environment Operator
===========

Isolate kubernetes pod communication into virtual groups according to pod label, with route fallback support.

![isolation](https://virtual-environment.oss-cn-zhangjiakou.aliyuncs.com/image/diagram-en-us.jpg)

Check ☞[document](https://alibaba.github.io/virtual-environment/#/en-us/)☜ for more information.

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

## Support

Contact us with DingTalk:

<img src="https://virtual-environment.oss-cn-zhangjiakou.aliyuncs.com/image/dingtalk-group-en-us.jpg" width="50%"></img>
