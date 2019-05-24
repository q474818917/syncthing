## syncthing

+ 关于Erlang语言中的supervisor trees，参照http://diaocow.iteye.com/blog/1762895
+ 调试将
```$xslt
if innerProcess || options.noRestart {
    syncthingMain(options)
} else {
    monitorMain(options)
}
替换成
syncthingMain(options)
```

### 连接步骤
+ 添加deviceID
+ 计算要连接设备的IP和port
    + 1、静态绑定IP、port
    + 2、使用本地local discovery,因为syncthing会广播自己的信息，如果发现到这个deviceID，则就能找到IP和port
    + 3、使用global discovery
+ 建立tcp连接，执行tls握手，需要执行一些验证操作

