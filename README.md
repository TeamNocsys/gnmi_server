# gnmi server

> 更新openconfig子项目

```
$ git submodule update --init
```

> yang文件处理

```
$ cd third_party
$ go generate
```

> 启动服务

```
$ go run cmd/gnmi/main.go
```