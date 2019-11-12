## 使用示例

####构建镜像
```
docker build -t node-demo:latest .
```
####启动镜像

```
docker run -p 8888:8888 -e envMark="node" -e url="http://127.0.0.1:9090/demo" node-demo:latest
```

`envMark`环境标识，默认为`dev`

`url`此程序会get调用一个地址来测试透传效果，默认不调用