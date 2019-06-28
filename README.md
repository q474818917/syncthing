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

+ 拉取过程中会产生.xxx.tmp的中间文件

### 依赖列表
+ 关于Erlang语言中的supervisor trees，参照http://diaocow.iteye.com/blog/1762895
+ 文件系统通知库，github.com/syncthing/notify
+ go项目的leveldb库，github.com/syndtr/goleveldb
+ 从一个字符串中查找给定字符串，github.com/chmduquesne/rollinghash

### 连接步骤
+ 添加deviceID
+ 计算要连接设备的IP和port
    + 1、静态绑定IP、port
    + 2、使用本地local discovery,因为syncthing会广播自己的信息，如果发现到这个deviceID，则就能找到IP和port
    + 3、使用global discovery
+ 建立tcp连接，执行tls握手，需要执行一些验证操作

### 关于[BEP协议](doc/BEP.md), 详细的[messageType交互](doc/exchange.md)
+ 核心连接监听代码：
```$xslt
①model.go
func (m *model) AddConnection
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

### 事件订阅：

### 文件大小：
+ 可以通过以下命令查看文件字节大小
```$xslt
ls -l -k
```

+ wiki上有个区段表，如果目标文件大小在0 - 250 MiB，则每个区块的最大传输为128 KiB
```$xslt
例如：196426字节的文件，按每块128KiB切分，第一块也就是128*1024（131072），第二块用总数-第一块（65354）

```
### goroutine核心拉取逻辑：
+ 逻辑都在folder_sendrecv.go里
```
pull() -> pullerIteration() -> processNeeded() -> handleFile() 塞入 copyChan
                  goroutine -> copierRoutine                    -> 触发 copierRoutine, 塞入pullChan
                  goroutine -> pullerRoutine                            -> 触发 pullerRoutine，塞入finisherChan
                  goroutine -> finisherRoutine                                 -> 触发 finisherRoutine 
```