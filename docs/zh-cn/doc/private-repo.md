# 使用私有镜像仓库部署

KtEnv的两个主要组件（Webhook和Operator）的镜像默认托管在Docker Hub公共仓库，该仓库对未登录Docker Hub的客户端有每天5000次的拉取限制，若您遇到了错误 "ERROR: toomanyrequests: Too Many Requests."，则说明当天使用此镜像的用户数已经超过限值了。

为了让部署更加稳定，建议您在部署KtEnv前，通过（已登录过Docker Hub的）本地Docker客户端将镜像拉取并推送到自己的私有镜像仓库，并将KtEnv组件的镜像改为从私有仓库拉取。

## 转存镜像到私库

首先将KtEnv组件的镜像拉到本地（如遇到拉取次数超限的错误，可先用`docker login`登录您的Docker Hub账号）

```bash
docker pull virtualenvironment/virtual-env-operator:v0.6.0
docker pull virtualenvironment/virtual-env-admission-webhook:v0.6.0
```

然后将镜像重新命名，并推到自己的私有仓库（将以下命令中`<您的仓库地址>`替换为实际仓库地址）

```bash
docker tag virtualenvironment/virtual-env-operator:v0.6.0 <您的仓库地址>/virtual-env-operator:v0.6.0
docker tag virtualenvironment/virtual-env-admission-webhook:v0.6.0 <您的仓库地址>/virtual-env-admission-webhook:v0.6.0
docker push <您的仓库地址>/virtual-env-operator:v0.6.0
docker push <您的仓库地址>/virtual-env-admission-webhook:v0.6.0
```

## 更新部署组件镜像

编辑部署包中的`global/ktenv_webhook.yaml`和`ktenv_operator.yaml`文件，将其中的`image:`参数值替换为私有仓库镜像地址，在Linux下可通过以下命令完成：

```bash
sed -i 's#virtualenvironment/virtual-env-operator:#<您的仓库地址>/virtual-env-operator:#' ktenv_operator.yaml
sed -i 's#virtualenvironment/virtual-env-admission-webhook:#<您的仓库地址>/virtual-env-admission-webhook:#' global/ktenv_webhook.yaml
```

> 此命令也适用于Mac系统，但需要将命令中的`sed -i`改为`sed -i ''`
