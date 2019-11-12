## 使用示例
####构建应用
```
mvn clean package
```

####构建镜像
```
docker build -t springboot-demo:latest .
```
####启动镜像

```
docker run -p 8080:8080 -e "envMark=v1" springboot-demo:latest
```

`envMark`环境标识，默认为`dev`

`url`此程序会get调用一个地址来测试透传效果，默认不调用