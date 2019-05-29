## syncthing

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

+ 显示debug日志，在idea下environment设置STTRACE=model，然后debug go build github.com/syncthing/syncthing/cmd/syncthing 

### 依赖列表
+ 关于Erlang语言中的supervisor trees，参照http://diaocow.iteye.com/blog/1762895
+ 

### 连接步骤
+ 添加deviceID
+ 计算要连接设备的IP和port
    + 1、静态绑定IP、port
    + 2、使用本地local discovery,因为syncthing会广播自己的信息，如果发现到这个deviceID，则就能找到IP和port
    + 3、使用global discovery
+ 建立tcp连接，执行tls握手，需要执行一些验证操作

### 关于[BEP协议](BEP.md)
+ 连接监听代码：
```$xslt
①func (m *model) AddConnection
↓
②conn.Start() = (rawConnection.Start)
↓
③调用收发消息核心
func (c *rawConnection) Start() {
	go func() {
		err := c.readerLoop()
		c.internalClose(err)
	}()
	go c.writerLoop()
	go c.pingSender()
	go c.pingReceiver()
}
↓
④下面就是writerLoop、readerLoop按照事件去处理
```

+ supervisor服务启动,folder实现suture的service接口。服务启动时就会调用各个实现Serve()方法
```$xslt
① folder.go
case <-f.pullScheduled:
    pullFailTimer.Stop()
    select {
    case <-pullFailTimer.C:
    default:
    }
    pull()
↓

```

### 文件夹三种同步方式：
+ 默认：sendreceive
+ sendonly
+ recvonly
```$xslt
①folder.go //监控文件夹，通过文件系统和定时任务实现
if f.FSWatcherEnabled && f.CheckHealth() == nil {
    f.startWatch()
}
↓
②filesystem.go //实现了Watch方法的 basicfs_watch.go
eventChan, err := f.Filesystem().Watch(".", f.ignores, ctx, f.IgnorePerms)
↓
③basicfs_watch.go
func (f *BasicFilesystem) Watch(name string, ignore Matcher, ctx context.Context, ignorePerms bool) (<-chan Event, error) {

case outChan <- Event{Name: relPath, Type: evType}://发送文件Chan


```