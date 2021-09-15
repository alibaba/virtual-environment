# Best practice pattern

Cooperating with the code branch management mode used by the project, virtual environment can be used in certain patterns to form an optimal R&D practice.

Take Gitflow as an example. Firstly, use a CI pipeline to automatically deploy the `develop branch` code to a test environment, and add a `dev` tag to the Pob template of the corresponding Deployment resources. These Pods will be used as the final fallback route targets, shared with all users in the team.

Developers who don't want to use service deployed by `develop branch`, can deploy a service instance with the image created by a `feature branch` and tag it with a child-level tag, such as `dev.proj-1`. These environments are named "project environment", they can also be shared with other developers who rely on the specified version of the service (usually are developers working on the same specific project).

Service instances on Developer's local machine can also be added to the isolated project environment (instances in the same project environment will be able to call these local service instances). You can also create a child virtual environment separately based on a project environment or shared environment, using tags such as `dev.proj-1.jinji` or` dev.jinji`, which are called "personal environments" (can be usefully when testing or debugging local service against other services in cluster, and ensure that other developers' requests will not enter your own local service instances).

Above practice will finally form a three-tier virtual environment pattern, as shown in the figure below.

![best-practice-en-us.jpg](https://img.alicdn.com/imgextra/i1/O1CN01VgbFUv1X7iBfS3iUb_!!6000000002877-0-tps-2308-1324.jpg)

- **Shared environment**: Contains all required services, as final route targets, should use stable service version
- **Project environment**: Contains service instances used by certain project, usually for cooperative development and debugging, use service version specified by project
- **Personal environment** (include local service instances): Contains minimal services, for personal debugging

This three-tier virtual environment pattern is also applicable to other branch management models.
