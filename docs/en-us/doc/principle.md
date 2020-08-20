# Principle

Essentially, the Virtual Environment Operator is a Service Mesh rule controller, currently implementation based on Istio open-sourced version.

After the Operator started, it would traverse all label in pod template of Deployments in its Namespace, and calculate Service subset visiting rule for each possible virtual-environment Header on HTTP request. Then continually watch for events on all Services and Deployments, dynamically create or adjust VirtualService and DestinationRule instances.

![calculate-rule](https://virtual-environment.oss-cn-zhangjiakou.aliyuncs.com/image/calculate-rule-en-us.jpg)
