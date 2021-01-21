# SimpleFTP

要求实现一个FTP服务器，支持cd，ls，put，get等命令，目前实现了用户身份简单确认，获取家目录后可以进行cd，ls，mkdir以及上传和下载文件。

>TODO：
>- 未实现输入密码时不回显（类似C里的getpass函数）；
>- 不支持文件夹的上传与下载；
>- 未实现与linux用户权限管理保持一致。


## FTP
Data Format  
```
Format:   |TotalLength|              content             |
[]byte:   |           |                                  |
Size:     |     4B    |                                  |
ByteOrder: binary.LittleEndian
```

Read/Write Function
```
binary.Read((r io.Reader, order ByteOrder, data interface{})  
binary.Write((w io.Writer, order ByteOrder, data interface{})
```

## Server
服务端

- init
```
sync.Once.Do(initUsers()) // 读取../ex8.2/server/ftp/users文件缓存用户信息
```

- validate
```
// 验证用户名和密码，返回验证结果true/false和验证通过后的用户家目录
func Validate(name string, pwd string) (pass bool, home string) {...}
```
> pass == false: return ErrorCode  
> pass == true:
```
ftpCon := ftp.FtpConn{
    Con:  con,
    Home: cwd,
    Cwd:  cwd,
}
ftpServer := server.FtpServer{
    ftpCon,
}
// 循环监听命令请求
for !ftpServer.Exit{...}
```

- handle  
> cd
```
func (ftpCon *FtpServer) HandleCd(args []byte) error {...}
```

> ls
```
 func (ftpCon *FtpServer) HandleLs(args []byte) error {...}
 ```

> exit
```
func (ftpCon *FtpServer) HandleExit(args []byte) error {...}
```

> mkdir
```
func (ftpCon *FtpServer) HandleMkdir(args []byte) error {...}
```

> put
```
func (ftpCon *FtpServer) HandlePut(args []byte) error {...}
```

> get
```
func (ftpCon *FtpServer) HandleGet(args []byte) error {...}
```

## Client
客户端

- login
> Account:
>```
>go run -user user
>```

```
// 连接到ftp服务器
con, err := net.Dial("tcp", host)
```

- op
```
// 监听命令行输入
for input.Scan() && !ftpClient.Exit {...}
```