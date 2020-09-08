# gnmi server

## 编译二进制文件

### 测试用
```shell script
$ cd cmd/gnmi
$ go build -tags debug
```

当前目录下生成可执行文件gnmi
```shell script
$ ls
gnmi  gnmi.go
```

### 发布用
```shell script
$ cd cmd/gnmi
$ go build -tags release
```

当前目录下生成可执行文件gnmi
```shell script
$ ls 
gnmi  gnmi.go
```

## 运行

1.将编译好的可执行文件复制到SONiC上
```shell script
$ scp gnmi admin@192.168.200.47:/home/admin
```

2.运行
```shell script
$ sudo ./gnmi run --address 0.0.0.0 --port 5002
```

## 测试

### 安装测试用客户端
1.执行下面的命令
```shell script
$ go get github.com/aristanetworks/goarista/cmd/gnmi
```
命令执行完成后在go目录下bin文件夹内生成gnmi可执行文件

2.将$GOPATH/bin添加到系统的PATH中
```shell script
$ export PATH=$PATH:`go env GOPATH`/bin
```

3.测试
```shell script
$ gnmi -addr 192.168.200.47:5005 get /sonic-port/port/port-state-list/counters
```

