## device exchange example
# | A | B
---|---|---
1 | ClusterConfiguration->  | <-ClusterConfiguration
2 | Index->  | <-Index
3 | IndexUpdate->   | <-IndexUpdate
4 | IndexUpdate->   | 
5 | Request-> 	 |
6 | Request-> 	 |
7 | Request-> 	 |
8 | Request-> 	 |
9 |  |	<-Response
10 |  |	<-Response
11 |  |	<-Response
12 |  |	<-Response
13 | Index Update-> |
... |   |
14 |  | <-Ping
15 | Ping-> |  

ps:A 发现有4个block缺失，然后发送4个request给B，请求4个block。而经过测试并未发现此场景。

### ClusterConfig
①model.go
func (m *model) ClusterConfig(deviceID protocol.DeviceID, cm protocol.ClusterConfig)
↓
go sendIndexes(conn, folder.ID, fs, startSequence, dropSymlinks)
↓
sub := events.Default.Subscribe(events.LocalIndexUpdated | events.DeviceDisconnected)
↓
func sendIndexTo(prevSequence int64, conn protocol.Connection, folder string, fs *db.FileSet, dropSymlinks bool) (int64, error)

### Index
①model.go
func (m *model) Index(deviceID protocol.DeviceID, folder string, fs []protocol.FileInfo) 

### Index Update
①model.go
func (m *model) IndexUpdate(deviceID protocol.DeviceID, folder string, fs []protocol.FileInfo)

