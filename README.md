## go-tcp-proxy

### 简介
TCP代理服务，可以实现对端通讯，并且是TLS加密

### 编译和安装
```shell
# make
```
编译后的文件在bin/tcp-proxy

### 使用方法

```toml
debug=true

[proxy]
    [proxy.openapi]
        enabled = true
    [proxy.openapi.listen]
        addr = "0.0.0.0:81"
        tls = true
        ca = "./certs/client.pem"
        privFile = "./certs/server.pem"
        pubFile = "./certs/server.key"
    [proxy.openapi.remote]
        addr = "127.0.0.1:80"```

* proxy.openapi.listen 表示本地监控的信息
* proxy.openapi.remote 表示远端信息

`注：`两端均支持TLS加密