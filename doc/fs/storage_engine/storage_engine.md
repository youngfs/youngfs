**### 存储引擎接口

对于任何其他的存储引擎，其接口通过组合逻辑需要满足 `interface StorageEngine (storage_engine/storage_engine.go)`

对于存储引擎和文件系统，可以通过多种协议进行交互，如 `http`, `grpc`等等

对于新增加的存储引擎，存储引擎出现了错误需要在 `errors/errors.go` 下面新增对应的错误类型进行返回

对于存储引擎接口，需要维护每个文件的链接数，具体维护方式可以是存储引擎自己维护，也可以在文件系统当中利用数据库进行维护

总共有如下的接口：

#### 上传文件

```
    PutObject(ctx context.Context, size uint64, reader io.Reader, hosts ...string) (string, error)
```

其中参数为

|   参数   |    类型     |                                    说明                                     |
|:------:|:---------:|:-------------------------------------------------------------------------:|
|  size  |  uint64   |                传入文件的大小，单位为byte，用于提前告知存储引擎传输文件的大小，从而更好的进行分配                |
| reader | io.Reader |                                 实际文件的io读流                                 |
| hosts  | []string  |        表示上传到哪几台主机当中，如果为空则表示上传到任意一台主机均可，如果长度不为1则表示上传到hosts中任意一台主机均可        |

其返回一个 `fid` (类型:  `string`)， `fid`的格式任意，在存储引擎端自己进行解析，但需要保证任何文件拥有全局唯一的 `fid`

#### 删除文件

```
    DeleteObject(ctx context.Context, fid string) error
```

其中参数为

| 参数  |    类型     |     说明      |
|:---:|:---------:|:-----------:|
| fid |  string   | 对应fid的对象删除  |

#### 下载文件

```
	GetObject(ctx context.Context, fid string, writer io.Writer) error
```

其中参数为

|   参数    |    类型     |       说明        |
|:-------:|:---------:|:---------------:|
|   fid   |  string   |  请求获取文件URL的fid  |
| writer  | io.Writer | 将下载的数据信息写入的io写流 |

#### 获得存储引擎下的所有HOST

```
    GetHosts(ctx context.Context) ([]string, error)
```

获得所有存储引擎下控制的主机的`host`