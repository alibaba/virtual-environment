# Principle

Essentially, the Virtual Environment Operator is a Service Mesh rule controller, currently implementation based on Istio open-sourced version.

After the Operator started, it would traverse all label in pod template of Deployments in its Namespace, and calculate Service subset visiting rule for each possible virtual-environment Header on HTTP request. Then continually watch for events on all Services and Deployments, dynamically create or adjust VirtualService and DestinationRule instances.

![calculate-rule-en-us.jpg](https://img.alicdn.com/imgextra/i2/O1CN01Szvy3O1fRsYLPnVQb_!!6000000004004-0-tps-1620-440.jpg)
