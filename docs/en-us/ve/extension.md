# Function extension

In order to support other Service Mesh controller, you should implement [IsolationRouter](https://github.com/alibaba/virtual-environment/blob/master/pkg/component/router/router_interface.go) interface and define a new Router class, and register it to the pool in [router_pool.go](https://github.com/alibaba/virtual-environment/blob/master/pkg/component/router/router_pool.go). Then select route type when configuring the virtual environment instance (**Notice**: this configure item currently not exist yet, please raise a github issue when required)
