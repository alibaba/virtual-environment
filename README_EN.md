KtEnv (Virtual Environment Operator)
===========

Isolate kubernetes pod communication into virtual groups according to pod label, with route fallback support.

![diagram-en-us.jpg](https://img.alicdn.com/imgextra/i1/O1CN01NNA5Cm1XV4NwiFqJ2_!!6000000002928-0-tps-2160-884.jpg)

Check ☞[document](https://alibaba.github.io/virtual-environment/#/en-us/)☜ for more information.

[中文介绍](./README.md)

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

## Support

Contact us with DingTalk:

<img src="https://img.alicdn.com/imgextra/i4/O1CN01sTW3D61NzAFgUCNqz_!!6000000001640-0-tps-573-657.jpg" alt="dingtalk-group-en-us.jpg" width="50%"></img>
