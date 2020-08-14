# SimpleFTP

Ҫ��ʵ��һ��ftp��������֧��cd��ls��put��get�����Ŀǰʵ�����û���ݼ�ȷ�ϣ���ȡ��Ŀ¼����Խ���cd��ls��mkdir�Լ��ϴ��������ļ���

>TODO��
>- δʵ����������ʱ�����ԣ�����C���getpass��������
>- ��֧���ļ��е��ϴ������أ�
>- δʵ����linux�û�Ȩ�޹�����һ�¡�


## ftp
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

## server
�����

- init
```
sync.Once.Do(initUsers()) // ��ȡ../ex8.2/server/ftp/users�ļ������û���Ϣ
```

- validate
```
// ��֤�û��������룬������֤���true/false����֤ͨ������û���Ŀ¼
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
// ѭ��������������
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

## client
�ͻ���

- login
> Account: username@host
> Password: password
```
// ���ӵ�ftp������
con, err := net.Dial("tcp", host)
```

- op
```
// ��������������
for input.Scan() && !ftpClient.Exit {...}
```