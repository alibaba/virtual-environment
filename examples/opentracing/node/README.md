## 使用示例

### 构建镜像

```
docker build -t node-demo:latest .
```

### 启动镜像

Linux

```
ip=`ip addr show eth0 | grep -oP '(?<=inet )[0-9.]+'`
docker run -p 8001:8080 -e envMark="node" -e url="http://${ip}:8002/demo" node-demo:latest
```

Mac

```
ip=`ip addr show en0 | grep 'inet ' | sed 's/.*inet \([0-9.]*\).*/\1/g'`
docker run -p 8001:8080 -e envMark="node" -e url="http://${ip}:8002/demo" node-demo:latest
```

- `envMark`环境标识，默认为`dev`
- `url`此程序会get调用一个地址来测试透传效果，默认不调用