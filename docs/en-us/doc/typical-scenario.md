# Typical scenarios

#### 1. Single user debugging

Developer connect local services to the cluster for debugging through KtConnect. VirtualEnvironment can ensure that developers always access the local instance themselves, and normal calls by other developers will not enter the local instance.

<img src="https://virtual-environment.oss-cn-zhangjiakou.aliyuncs.com/image/typical-scenario-1.jpg" height="300px"/>

#### 2. Multi-user cooperation

Multiple developers add local services to the same virtual environment through KtConnect. VirtualEnvironment can ensure that the calls between these developers are interoperable, so that the team can be coordinated without affecting the normal use of other developers who have not entered the virtual environment.

<img src="https://virtual-environment.oss-cn-zhangjiakou.aliyuncs.com/image/typical-scenario-2.jpg" height="300px"/>

#### 3. Replace unstable service

When performing functional verification, a certain service is required to use the specified unstable version. In order not to affect other developers' use of the shared environment, the specified version can be deployed in an isolated virtual environment and used by itself.

<img src="https://virtual-environment.oss-cn-zhangjiakou.aliyuncs.com/image/typical-scenario-3.jpg" height="300px"/>

#### 4. Integration test isolation

During the integration test, the tested version of the service is placed in an isolated environment, and other service instances in the shared environment are reused, so that it is possible to quickly verify the function of a specific version of a specific service on the cluster without creating a full service cluster.

<img src="https://virtual-environment.oss-cn-zhangjiakou.aliyuncs.com/image/typical-scenario-4.jpg" height="200px"/>

#### 5. Quick multi-version comparison

Use a browser to set different header values ​​through plug-ins, quickly switch to access service instances that belong to different virtual environments, and compare the effects before and after.

<img src="https://virtual-environment.oss-cn-zhangjiakou.aliyuncs.com/image/typical-scenario-5.jpg" height="300px"/>
