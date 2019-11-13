## 使用示例
####构建应用
```
go mod init main.go
go build main.go
```

####构建镜像
```
docker build -t go-demo:latest .
```
####启动镜像

```
docker run -p 9090:9090  go-demo:latest --envMark=dev --url=http://127.0.0.1:8888/demo
```

`envMark`环境标识，默认为`dev`

`url`此程序会get调用一个地址来测试透传效果，默认不调用