## BEP协议
+ 关于具体得BEP协议内容，请查看bep.proto文件


## pre-authentication:验证前操作，发送Hello message
+ 建立连接后，在进行各种验证操作前， 设备间交换Hello message
+ Hello message消息组成,包括：Magic、length、Hello消息体
```$xslt
 0                   1
 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|             Magic             |
|           (32 bits)           |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|             Length            |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
/                               /
\             Hello             \
/                               /
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

Hello message消息体

message Hello {
    string device_name    = 1;
    string client_name    = 2;
    string client_version = 3;
}
```

## post-authentication:
```$xslt
 0                   1
 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|         Header Length         |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
/                               /
\            Header             \
/                               /
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|         Message Length        |
|           (32 bits)           |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
/                               /
\            Message            \
/                               /
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

对应的proto

message Header {
    MessageType        type        = 1;
    MessageCompression compression = 2;
}

enum MessageType {
    CLUSTER_CONFIG    = 0;
    INDEX             = 1;
    INDEX_UPDATE      = 2;
    REQUEST           = 3;
    RESPONSE          = 4;
    DOWNLOAD_PROGRESS = 5;
    PING              = 6;
    CLOSE             = 7;
}

enum MessageCompression {
    NONE = 0;
    LZ4  = 1;
}

```

## Header中的MessageType枚举类型，具体子字段，分析源码时用到再作分析

+ Cluster Config：当前连接的集群配置，此类型是必须是第一个发送的post authentication
```$xslt
对应的proto

message ClusterConfig {
    repeated Folder folders = 1;
}

message Folder {
    string id                   = 1;
    string label                = 2;
    bool   read_only            = 3;
    bool   ignore_permissions   = 4;
    bool   ignore_delete        = 5;
    bool   disable_temp_indexes = 6;
    bool   paused               = 7;

    repeated Device devices = 16;
}

message Device {
    bytes           id                         = 1;
    string          name                       = 2;
    repeated string addresses                  = 3;
    Compression     compression                = 4;
    string          cert_name                  = 5;
    int64           max_sequence               = 6;
    bool            introducer                 = 7;
    uint64          index_id                   = 8;
    bool            skip_introduction_removals = 9;
}

enum Compression {
    METADATA = 0;
    NEVER    = 1;
    ALWAYS   = 2;
}

```

+ Index and Index Update：Index包含发送者文件夹的所有内容，而Index Update只包含需要更新的Index
```$xslt
message Index {
    string            folder = 1;
    repeated FileInfo files  = 2;
}

message IndexUpdate {
    string            folder = 1;
    repeated FileInfo files  = 2;
}

message FileInfo {
    string       name           = 1;
    FileInfoType type           = 2;
    int64        size           = 3;
    uint32       permissions    = 4;
    int64        modified_s     = 5;
    int32        modified_ns    = 11;
    uint64       modified_by    = 12;
    bool         deleted        = 6;
    bool         invalid        = 7;
    bool         no_permissions = 8;
    Vector       version        = 9;
    int64        sequence       = 10;
    int32        block_size     = 13;

    repeated BlockInfo Blocks         = 16;
    string             symlink_target = 17;
}

enum FileInfoType {
    FILE              = 0;
    DIRECTORY         = 1;
    SYMLINK_FILE      = 2 [deprecated = true];
    SYMLINK_DIRECTORY = 3 [deprecated = true];
    SYMLINK           = 4;
}

message BlockInfo {
    int64 offset     = 1;
    int32 size       = 2;
    bytes hash       = 3;
    uint32 weak_hash = 4;
}

message Vector {
    repeated Counter counters = 1;
}

message Counter {
    uint64 id    = 1;
    uint64 value = 2;
}
```

+ Request：代表想要接收对等节点文件夹中的一个数据块的请求
```$xslt
message Request {
    int32  id             = 1;
    string folder         = 2;
    string name           = 3;
    int64  offset         = 4;
    int32  size           = 5;
    bytes  hash           = 6;
    bool   from_temporary = 7;
}

```

+ Response：对应着回应request消息
```$xslt
message Response {
    int32     id   = 1;
    bytes     data = 2;
    ErrorCode code = 3;
}

enum ErrorCode {
    NO_ERROR     = 0;
    GENERIC      = 1;
    NO_SUCH_FILE = 2;
    INVALID_FILE = 3;
}

```

+ DownloadProgress：通知远程设备的文件的部分可用性
```$xslt
一个DownloadProgress包含多个FileDownloadProgressUpdateType
message DownloadProgress {
    string                              folder  = 1;
    repeated FileDownloadProgressUpdate updates = 2;
}

message FileDownloadProgressUpdate {
    FileDownloadProgressUpdateType update_type   = 1;
    string                         name          = 2;
    Vector                         version       = 3;
    repeated int32                 block_indexes = 4;
}

enum FileDownloadProgressUpdateType {
    APPEND = 0;
    FORGET = 1;
}

```

+ Ping：确保长连接，nat和防火墙之类下的连接存活
```$xslt
message Ping {
}

```

+ Close
```$xslt
message Close {
    string reason = 1;
}

```