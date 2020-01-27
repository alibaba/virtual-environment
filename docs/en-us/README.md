# Introduction

KtVirtualEnvironment is a member of KT-series projects, supported by cloud development department of Alibaba Inc.

This project implement an isolation mechanism named virtual environment based on traffic tagging, fit for Kubernetes clusters. It can be used independently or combined with [KtConnect](https://alibaba.github.io/kt-connect/) tools to implement local-to-cluster traffic routing control, discussed at [typical scenario](ve/typical-scenario.md) in detail.

## Intention

For microservices developers, having a clean, dedicated and full-filled testing environment can undoubtedly improve the efficiency of debugging and troubleshooting in the software development process.

However, considering both economic costs and management costs, in any medium and large teams, maintaining a dedicated test clusters of all services for each developer is not a realistic idea. Therefore, Alibaba's R&D team has adopted a "virtual environment" approach based on route isolation.

![isolation](https://virtual-environment.oss-cn-zhangjiakou.aliyuncs.com/image/diagram-en-us.jpg)

According to the specified "virtual environment name" header on the request, this mechanism will form isolated traffic groups which composed of several specific service instances with other shared service instances, into full-filled testing environments from each developer's perspective.

This practice was also mentioned at [this](https://medium.com/hackernoon/lower-cost-with-higher-stability-how-do-we-manage-test-environments-at-alibaba-f7bd444ff6d2) article.

## Features

- Division of virtual isolation groups based on traffic tagging (request header)
- Reuse service instances from shared environment in different isolated groups
- Developer can redeploy or debug service in their project environment anytime, without affecting others
- Support for local running services directly join isolation group, without deploying to the cluster

## Connect us

Please join `kt-dev` DingTalk group：

<img src="https://virtual-environment.oss-cn-zhangjiakou.aliyuncs.com/image/dingtalk-group-en-us.jpg" width="40%"></img>
