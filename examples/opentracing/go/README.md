## 使用示例

### 构建镜像

```
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build main.go
docker build -t go-demo:latest .
```

### 启动镜像

Linux

```
ip=`ip addr show eth0 | grep -oP '(?<=inet )[0-9.]+'`
docker run -p 8002:8080 -e envMark=dev -e url=http://${ip}:8003/demo go-demo:latest
```

Mac

```
ip=`ip addr show en0 | grep 'inet ' | sed 's/.*inet \([0-9.]*\).*/\1/g'`
docker run -p 8002:8080 -e envMark=dev -e url=http://${ip}:8003/demo go-demo:latest
```

- `envMark`环境标识，默认为`dev`
- `url`此程序会get调用一个地址来测试透传效果，默认不调用