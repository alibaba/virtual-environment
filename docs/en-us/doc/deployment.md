# Using virtual environment

**Premise**: Kubernetes cluster already has Istio installed, local machine has kubectl installed and configured, check below links for detail:

- Istio：https://istio.io/docs/setup/install
- kubectl：https://kubernetes.io/docs/tasks/tools/install-kubectl

## Add operator to cluster

Download latest CRD from [release page](https://github.com/alibaba/virtual-environment/releases), use `kubectl apply` command to add the operator into Kubernetes

```bash
wget https://github.com/alibaba/virtual-environment/releases/download/v0.2/kt-virtual-environment-v0.2.zip
unzip kt-virtual-environment-v0.2.zip
cd v0.2/
kubectl apply -f crds/env.alibaba.com_virtualenvironments_crd.yaml
```

Put the operator into any namespaces which require virtual environment, e.g. `default`

```bash
kubectl apply -n default -f operator.yaml
```

If the cluster has RBAC enabled, please also apply Role and ServiceAccount

```bash
kubectl apply -n default -f service_account.yaml
kubectl apply -n default -f role.yaml
kubectl apply -n default -f role_binding.yaml
```

Now, the Kubernetes cluster already has capability to empower virtual environment.

## Check deployment result

KtEnv project contents two components as `Operator CRD` and `Admission Webhook`.
The `Webhook` component is used to write pod's virtual environment label into the runtime environment variable of its sidecar container.
The `CRD` component would listener to VirtualEnvironment resource instance created and dynamically generates routing rules according to the service status in the cluster.

The `Webhook` component is deployed in `kt-virtual-environment` namespace by default, contents a `Service`, a `Deployment` instance and other sub-resources created by them.
You could check the deployment status by below command:

```bash
kubectl -n kt-virtual-environment get all
```

If the output is similar to the following information, it indicates that the Webhook component of KtEnv has been deployed and is running normally.

```
NAME                                  READY   STATUS    RESTARTS   AGE
pod/webhook-server-5dd55c79b5-rf6dl   1/1     Running   0          86s

NAME                     TYPE        CLUSTER-IP     EXTERNAL-IP   PORT(S)   AGE
service/webhook-server   ClusterIP   172.21.0.254   <none>        443/TCP   109s

NAME                             READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/webhook-server   1/1     1            1           109s

NAME                                        DESIRED   CURRENT   READY   AGE
replicaset.apps/webhook-server-5dd55c79b5   1         1         1       86s
```

Check the `AGE` attribute of each resource object in the above output (the time elapsed from creation to now) to determine whether the object was created by the newly deployed Webhook component.

The `CRD` component will add a resource type named `VirtualEnvironment` in the Kubernetes cluster, which we will use in the next step. The installation status can be verified by the following command:

```bash
kubectl get crd virtualenvironments.env.alibaba.com
```

If the output is similar to the following information, it indicates that the CRD component of KtEnv has been deployed correctly.

```
NAME                                  CREATED AT
virtualenvironments.env.alibaba.com   2020-04-21T13:20:35Z
```

Check the `CREATED AT` attribute (resource creation time) in the output to determine whether the object is a newly deployed CRD component.

## Create VirtualEnvironment instance

Create a resource with kind `VirtualEnvironment`, use `kubectl apply` command to take effect

```bash
kubectl apply -n default -f path-to-virtual-environment-cr.yaml
```

After VirtualEnvironment instance created, it would automatically watch all Service, Deployment and StatefulSet resource **in the same Namespace** and generate isolation rule, thus form the virtual environment.

Please refer to [virtualenv.yaml](https://github.com/alibaba/virtual-environment/blob/master/examples/deploy/virtualenv.yaml) as an example of the resource definition file,
and doc [configuration guide](en-us/doc/configuration.md) provide a more detail description of each configurable item,
please modify the parameters according to the actual situation.

## Application adaptation

According to the virtual environment configuration, add a virtual environment label to the Pod and let the service pass the virtual environment header down through the call chain.

- Put a label in pod template of deployment to identify which virtual environment it belongs to (the default label key is `virtual-env`)

- Add mechanism for passing down environment name header in application (the default header key is `X-Virtual-Env`)

Check [quick start](en-us/doc/quickstart.md) for a completed demonstrate.
