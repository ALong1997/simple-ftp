package main

import (
	"SimpleFTP/common"
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"os"
	"path"
	"strings"
)

type FTPClient struct {
	*common.FtpConn
}

// New FTPClient
func NewClient(con net.Conn, cwd, home string) *FTPClient {
	return &FTPClient{
		&common.FtpConn{
			Con:  con,
			Cwd:  cwd,
			Home: home,
			Exit: false,
		},
	}
}

// New Dail
func NewDail(address string) (*FTPClient, error) {
	con, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return NewClient(con, "", ""), nil
}

// Close client
func (c *FTPClient) Close() error {
	return c.Con.Close()
}

// Identity Authentication
func (c *FTPClient) Auth(user, pwd string) (err error) {
	err = c.WriteContent(common.Str2sbyte(user))
	if err != nil {
		fmt.Println(err)
		return err
	}

	err = c.WriteContent(common.Str2sbyte(pwd))
	if err != nil {
		fmt.Println(err)
		return err
	}

	var res uint32
	err = c.Read(&res)
	if err != nil {
		fmt.Println(err)
		return err
	}
	if res == 0 {
		return common.AuthenticationErr
	}

	cwd := make([]byte, res)
	err = c.Read(cwd)
	if err != nil {
		fmt.Println(err)
		return err
	}
	c.Cwd = common.Sbyte2str(cwd)
	c.Home = c.Cwd

	return nil
}

// handle command
func (c *FTPClient) handleCommand(command string, args []string) error {
	cmdId, err := common.GetCommandId(command)
	if err != nil {
		return err
	}

	err = c.sendCommand(cmdId, args)
	if err != nil {
		return err
	}

	if cmdId == common.GetId {
		err = c.handleGet(args[0])
		if err != nil {
			return err
		}
	}

	var length uint32
	err = c.Read(&length)
	if err != nil {
		return err
	}
	if length == 0 {
		fmt.Printf("\n%s:", c.Cwd)
		return nil
	}

	res := make([]byte, length-uint32(binary.Size(length)))
	err = c.Read(res)
	if err != nil {
		return err
	}
	if cmdId == common.CdId {
		c.Cwd = common.Sbyte2str(res)
		fmt.Printf("\n%s:", c.Cwd)
		return nil
	}
	if cmdId == common.ExitId {
		c.Exit = true
		fmt.Printf("%s\n", common.Sbyte2str(res))
		return nil
	}

	fmt.Printf("%s\n%s:", common.Sbyte2str(res), c.Cwd)
	return nil
}

// send command
func (c *FTPClient) sendCommand(cmdId uint8, args []string) error {
	if cmdId == common.PutId {
		return c.handlePut(cmdId, args[0])
	}

	argStr := strings.Join(args, "")
	return c.WriteCommand(cmdId, argStr)
}

// handle put
func (c *FTPClient) handlePut(cmdId uint8, filePath string) error {
	filePath = strings.Replace(filePath, "\\", "/", -1)
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	// 发送命令与文件名
	fileName := path.Base(filePath)
	err = c.WriteCommand(cmdId, fileName)
	if err != nil {
		return err
	}

	// 发送文件长度
	fileInfo, err := f.Stat()
	if err != nil {
		return err
	}
	if fileInfo.IsDir() {
		return common.PutDirErr
	} else {
		err = c.Write(fileInfo.Size())
		if err != nil {
			return err
		}
	}

	// 发送文件内容
	buf := make([]byte, maxBinarySize)
	bufReader := bufio.NewReader(f)
	for {
		n, err := bufReader.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		err = c.Write(buf[0:n])
		if err != nil {
			return err
		}
	}
	return nil
}

// handle get
func (c *FTPClient) handleGet(filePath string) error {
	fileName := path.Base(filePath)
	f, err := os.Create(fileName)
	if err != nil {
		if os.IsExist(err) {
			err = f.Truncate(0)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	defer f.Close()

	var length int64
	err = c.Read(&length)
	if err != nil {
		return err
	}
	var total, bufSize int64
	if length > maxBinarySize {
		bufSize = maxBinarySize
	} else {
		bufSize = length
	}
	buf := make([]byte, bufSize)
	for total < length {
		err = c.Read(buf)
		if err != nil {
			return err
		}
		n, err := f.Write(buf)
		if err != nil {
			return err
		}
		total += int64(n)
		if length-total < bufSize {
			buf = buf[0 : length-total]
		}
	}
	return nil
}
