# 配置Webhook组件

VirtualEnvironment产品中包含一个全局的Admission Webhook组件，他的主要作用是将Pod上的`环境标`信息通过环境变量注入到Sidecar容器里，便于Sidecar为出口流量的Header添加恰当的环境标。倘若集群中无需使用流量自动染色功能（即创建VirtualEnvironment资源时，`envHeader.autoInject`值始终为`false`），则可以无需部署此组件。

Webhook的配置位于发布包的`webhooks`子目录内，名称是`virtualenvironment_tag_injector_webhook.yaml`，其中包含两项可配置内容。

## envLabel变量

在配置文件的`Deployment`资源内，有一个名为`envLabel`环境变量，它的值需要与集群中存在的VirtualEnvironment资源的`envLabel.name`值匹配。倘若集群中含有多种不同`envLabel.name`取值的VirtualEnvironment资源，则应该将这些值全部列出来，用逗号“,”分隔。例如：

```yaml
env:
  - name: envLabel
    value: virtual-env,custom-virtual-env
```

## TLS证书和秘钥

由于Kubernetes要求Admission Webhook组件必须采用HTTPS协议监听，为了便于实现该约束，VirtualEnvironment的Webhook组件默认采用自签名的TLS证书，同时将证书与签名的值暴露在资源配置中，以便于自助替换。在较正式的使用场景中，强烈建议用户使用自己的企业签名证书或自签名证书替换默认的证书内容。以下将以Linux或Mac系统下生成自签名证书为例，介绍替换证书的方法。

首先需要按照OpenSSL工具，然后创建一个用于存放证书和秘钥的目录并进入到目录中。

```bash
mkdir keys
cd keys
```

使用OpenSSL命令行生成新的根证书和一组自签名秘钥文件。

```bash
# Generate the CA cert and private key
openssl req -nodes -new -x509 -keyout ca.key -out ca.crt -subj "/CN=Virtual Environment Admission Controller Webhook CA"
# Generate the private key for the webhook server
openssl genrsa -out webhook-server-tls.key 2048
# Generate a Certificate Signing Request (CSR) for the private key, and sign it with the private key of the CA.
openssl req -new -key webhook-server-tls.key -subj "/CN=webhook-server.kt-virtual-environment.svc" \
    | openssl x509 -req -CA ca.crt -CAkey ca.key -CAcreateserial -out webhook-server-tls.crt
```

将证书和秘钥的内容经过bash64编码后保存到环境变量中待用。

```bash
tls_crt_b64="$(openssl base64 -A < webhook-server-tls.crt)"
tls_key_b64="$(openssl base64 -A < webhook-server-tls.key)"
ca_pem_b64="$(openssl base64 -A < ca.crt)"
```

进入部署包中的`webhooks`目录（见[部署文档](deployment.md)），然后使用以下命令替换配置文件中的相应属性值。

```bash
cd webhooks
sed -i "s/tls.crt: .*/tls.crt: ${tls_crt_b64}/" virtualenvironment_tag_injector_webhook.yaml
sed -i "s/tls.key: .*/tls.key: ${tls_key_b64}/" virtualenvironment_tag_injector_webhook.yaml
sed -i "s/caBundle: .*/caBundle: ${ca_pem_b64}/" virtualenvironment_tag_injector_webhook.yaml
```
