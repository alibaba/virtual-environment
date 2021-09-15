# Typical scenarios

#### 1. Single user debugging

Developer connect local services to the cluster for debugging through KtConnect. VirtualEnvironment can ensure that developers always access the local instance themselves, and normal calls by other developers will not enter the local instance.

<img src="https://img.alicdn.com/imgextra/i4/O1CN010Yz4KT1xMX0to68mG_!!6000000006429-0-tps-1036-940.jpg" alt="typical-scenario-1.jpg" height="300px"/>

#### 2. Multi-user cooperation

Multiple developers add local services to the same virtual environment through KtConnect. VirtualEnvironment can ensure that the calls between these developers are interoperable, so that the team can be coordinated without affecting the normal use of other developers who have not entered the virtual environment.

<img src="https://img.alicdn.com/imgextra/i2/O1CN01LzaQQX1Wd1eQWEgRS_!!6000000002810-0-tps-1732-966.jpg" alt="typical-scenario-2.jpg" height="300px"/>

#### 3. Replace unstable service

When performing functional verification, a certain service is required to use the specified unstable version. In order not to affect other developers' use of the shared environment, the specified version can be deployed in an isolated virtual environment and used by itself.

<img src="https://img.alicdn.com/imgextra/i2/O1CN01SdeV921qaXok6uq6z_!!6000000005512-0-tps-1702-986.jpg" alt="typical-scenario-3.jpg" height="300px"/>

#### 4. Integration test isolation

During the integration test, the tested version of the service is placed in an isolated environment, and other service instances in the shared environment are reused, so that it is possible to quickly verify the function of a specific version of a specific service on the cluster without creating a full service cluster.

<img src="https://img.alicdn.com/imgextra/i1/O1CN01d4Bb0H1MvfcItqQG7_!!6000000001497-0-tps-1822-622.jpg" alt="typical-scenario-4.jpg" height="200px"/>

#### 5. Quick multi-version comparison

Use a browser to set different header values ​​through plug-ins, quickly switch to access service instances that belong to different virtual environments, and compare the effects before and after.

<img src="https://img.alicdn.com/imgextra/i2/O1CN01hU3ORL1LKk4FNXjFr_!!6000000001281-0-tps-912-932.jpg" alt="typical-scenario-5.jpg" height="300px"/>
