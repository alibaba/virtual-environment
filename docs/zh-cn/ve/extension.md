# 功能扩展

如需支持其他类型的Service Mesh控制面，只需要实现[IsolationRouter](https://github.com/alibaba/virtual-environment/blob/master/pkg/component/router/router_interface.go)接口的方法，自定义新的Router类型，并注册到[router_pool.go](https://github.com/alibaba/virtual-environment/blob/master/pkg/component/router/router_pool.go)的pool列表中，然后在创建虚拟环境实例时选择路由类型即可（**注意**：此配置当前尚未实现，如有需要请留意提出）。
