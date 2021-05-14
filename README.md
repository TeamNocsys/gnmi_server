# gnmi server

## 编译二进制文件

### 测试用
```shell script
$ cd cmd/gnmi
### 根据SONIC版本选择不同编译参数，目前仅支持broadcom和edgecore
$ go build -tags "debug broadcom"
$ go build -tags "debug ec"
```

当前目录下生成可执行文件gnmi
```shell script
$ ls
gnmi  gnmi.go
```

### 发布用
```shell script
### 根据SONIC版本选择不同编译参数，目前仅支持broadcom和edgecore
$ cd cmd/gnmi
$ GOOS=linux go build -tags "release broadcom"
$ GOOS=linux go build -tags "release ec"
### 支持远程DEBUG
$ GOOS=linux go build --gcflags "all=-N -l" -tags "release broadcom"
$ GOOS=linux go build --gcflags "all=-N -l" -tags "release ec"

### 制作deb包，交换机需要安装daemonize才能正常运行
### 在Debian主机上执行
### 其中GO指令可以在任意平台执行，然后将gnmi拷贝到build/deb/usr/local/bin，即可注释掉脚本的GO指令，默认编译broadcom版本
$ sudo apt install -y dh-make dpkg-dev devscripts
$ ./build_deb.sh
```

当前目录下生成可执行文件gnmi
```shell script
$ ls 
gnmi  gnmi.go
```

## 运行

1.将编译好的可执行文件复制到SONiC上
```shell script
### 1.1 普通二进制安装(用于调试)
$ scp gnmi admin@192.168.200.47:/home/admin

### 1.2 deb包安装
$ scp gnmi.deb admin@192.168.200.47:/home/admin
```

2.运行
```shell script
### 2.1 普通二进制安装(用于调试，根据实际SSH认证信息修改用户名密码)
### 日志默认路径为/var/log/gnmi_server.log
$ sudo ./gnmi run --address 0.0.0.0 --port 5002 --username admin --password YourPaSsWoRd
### 指定数据库配置文件
$ sudo ./gnmi run --address 0.0.0.0 --port 5002 -v --config /var/run/redis/sonic-db/database_config.json --username admin --password YourPaSsWoRd
### 远程DEBUG模式，注意ROOT用户运行
# dlv --listen=:2345 --headless=true --api-version=2 --accept-multiclient exec ./gnmi -- run --address 0.0.0.0 --port 5002 --path . --username admin --password YourPaSsWoRd

### 2.2 deb包安装(根据实际SSH认证信息修改gnmi.service)
$ sudo dpkg -i gnmi.deb
$ sudo systemctl start gnmi
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

