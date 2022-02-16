# 配置Webhook组件

KtEnv产品中包含一个全局的Admission Webhook组件，他的主要作用是将Pod上的"环境标签"信息通过环境变量注入到Sidecar容器里，便于Sidecar为出口流量的Header添加恰当的环境标签。倘若集群中无需使用流量自动染色功能（即创建VirtualEnvironment资源时，`envHeader.autoInject`值始终为`false`），则可以无需部署此组件。

## 配置参数

Webhook组件的配置参数位于KtEnv发布包`global`子目录内`ktenv_webhook.yaml`文件的`Deployment`对象内，包含两项可配置变量。

**envLabel环境变量**

参数`envLabel`的值需要与集群中存在的VirtualEnvironment资源的`envLabel.name`值匹配。倘若集群中含有多种不同`envLabel.name`取值的VirtualEnvironment资源，则应该将这些值全部列出来，用逗号“,”分隔。例如：

```yaml
env:
  - name: envLabel
    value: virtual-env,custom-virtual-env
```

**logLevel环境变量**

参数`logLevel`通常不需要修改，他会影响Webhook输出日志的密集程度，可选值如下：

- ERROR: 只记录错误信息，输出的日志量最少
- INFO: 输出异常错误和正常情况下的自动加标记录（默认值）
- DEBUG: 输出包括访问记录在内的所有日志，通常只在排查问题的时候使用

可以直接修改`ktenv_webhook.yaml`文件并通过`kubectl apply`使之生效；或直接通过`kubectl edit`命令修改`kt-virtual-environment`Namespace中名为`webhook-server`的Deployment对象完成配置的修改。

## TLS证书和秘钥

由于Kubernetes要求Admission Webhook组件必须采用HTTPS协议监听，为了便于实现该约束，VirtualEnvironment的Webhook组件默认采用自签名的TLS证书，同时将证书与签名的值暴露在资源配置中，以便于自助替换。在较正式的使用场景中，强烈建议用户使用自己的企业签名证书或自签名证书替换默认的证书内容。以下将以Linux或Mac系统下生成自签名证书为例，介绍替换证书的方法。

首先需要按照OpenSSL工具，然后创建一个用于存放证书和秘钥的目录并进入到目录中。

```bash
mkdir keys
cd keys
```

在此目录下创建一个配置文件`ssl.conf`，内容如下：

```text
[req_ext]
subjectAltName = @alt_names

[alt_names]
DNS.1 = webhook-server.kt-virtual-environment.svc
```

使用OpenSSL命令行生成新的根证书和一组自签名秘钥文件。

```bash
# Generate the CA cert and private key
openssl req -nodes -new -x509 -days 3650 -keyout ca.key -out ca.crt -subj "/CN=Virtual Environment Admission Controller Webhook CA"
# Generate the private key for the webhook server
openssl genrsa -out webhook-server-tls.key 2048
# Generate a Certificate Signing Request (CSR) for the private key, and sign it with the private key of the CA.
openssl req -new -key webhook-server-tls.key -subj "/CN=webhook-server.kt-virtual-environment.svc" \
    | openssl x509 -req -days 3650 -CA ca.crt -CAkey ca.key -CAcreateserial -out webhook-server-tls.crt -extensions req_ext -extfile ssl.conf
```

将证书和秘钥的内容经过bash64编码后保存到环境变量中待用。

```bash
tls_crt_b64="$(openssl base64 -A < webhook-server-tls.crt)"
tls_key_b64="$(openssl base64 -A < webhook-server-tls.key)"
ca_pem_b64="$(openssl base64 -A < ca.crt)"
```

进入部署包中的`global`目录（见[部署文档](zh-cn/doc/deployment.md?id=部署KtEnv组件)），然后使用以下命令替换配置文件中的相应属性值。

```bash
cd global
sed -i "s/tls.crt: .*/tls.crt: ${tls_crt_b64}/" ktenv_webhook.yaml
sed -i "s/tls.key: .*/tls.key: ${tls_key_b64}/" ktenv_webhook.yaml
sed -i "s/caBundle: .*/caBundle: ${ca_pem_b64}/" ktenv_webhook.yaml
```

部署或重新部署Webhook组件使修改生效：

```bash
kubectl apply -f ktenv_webhook.yaml
```
