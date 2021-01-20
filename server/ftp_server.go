package main

import (
	"SimpleFTP/common"
	"bufio"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
)

type FTPServer struct {
	common.FtpConn
}

// Handle cd
func (server *FTPServer) HandleCd(args []byte) error {
	cwd := common.Sbyte2str(args)
	if strings.HasPrefix(cwd, "/") {
		cwd = path.Join(server.Cwd, cwd)
	}

	f, err := os.Open(cwd)
	if err != nil {
		return server.WriteContent(common.Str2sbyte(err.Error()))
	}
	defer f.Close()

	fileInfo, err := f.Stat()
	if err != nil {
		return server.WriteContent(common.Str2sbyte(err.Error()))
	}
	if !fileInfo.IsDir() {
		return server.WriteContent(common.Str2sbyte(common.CDArgsErr.Error()))
	}

	server.Cwd = cwd

	return server.WriteContent(common.Str2sbyte(cwd))
}

// Handle ls
func (server *FTPServer) HandleLs(args []byte) error {
	cwd := common.Sbyte2str(args)
	if strings.HasPrefix(cwd, "/") {
		cwd = path.Join(server.Cwd, cwd)
	}
	if len(cwd) == 0 {
		cwd = server.Cwd
	}

	f, err := os.Open(cwd)
	if err != nil {
		return server.WriteContent(common.Str2sbyte(err.Error()))
	}

	fileInfo, err := f.Stat()
	if err != nil {
		return server.WriteContent(common.Str2sbyte(err.Error()))
	}
	if fileInfo.IsDir() {
		fileInfos, err := f.Readdir(0)
		if err != nil {
			server.WriteContent(common.Str2sbyte(err.Error()))
		}
		var res string
		res = fmt.Sprintf("Total:%d\n", len(fileInfos))
		for _, info := range fileInfos {
			res = res + fmt.Sprintf("%.30s\t%.10d\t%s\n", info.Name(), info.Size(), info.ModTime())
		}
		err = server.WriteContent(common.Str2sbyte(res))
	} else {
		res := fmt.Sprintf("%.30s\t%.10d\t%s\n", fileInfo.Name(), fileInfo.Size(), fileInfo.ModTime())
		err = server.WriteContent(common.Str2sbyte(res))
	}

	if err != nil {
		return server.WriteContent(common.Str2sbyte(err.Error()))
	}
	return nil
}

// Handle exit
func (server *FTPServer) HandleExit(args []byte) error {
	server.Exit = true
	return server.WriteContent(common.Str2sbyte("Byebye."))
}

// Handle mkdir
func (server *FTPServer) HandleMkdir(args []byte) error {
	dir := common.Sbyte2str(args)
	if strings.HasPrefix(dir, "/") {
		dir = path.Join(server.Home, dir)
	} else {
		dir = path.Join(server.Cwd, dir)
	}

	err := os.Mkdir(dir, os.ModePerm)
	if err != nil {
		return err
	}
	return server.WriteContent(common.Str2sbyte(common.NoErr.Error()))
}

// Handle put
func (server *FTPServer) HandlePut(args []byte) error {
	fileName := common.Sbyte2str(args)
	f, err := os.Create(path.Join(server.Cwd, fileName))
	if err != nil {
		return err
	}
	defer f.Close()

	var length int64
	err = server.Read(&length)
	if err != nil {
		return err
	}
	var total, bufSize int64
	if length > 4096 {
		bufSize = 4096
	} else {
		bufSize = length
	}

	buf := make([]byte, bufSize)
	for total < length {
		err = server.Read(buf)
		if err != nil {
			return err
		}
		n, err := f.Write(buf)
		if err != nil {
			return err
		}
		total += int64(n)
		if (length - total) < bufSize {
			buf = buf[:length-total]
		}
	}

	return server.WriteContent(common.Str2sbyte(common.NoErr.Error()))
}

// Handle get
func (server *FTPServer) HandleGet(args []byte) error {
	filePath := common.Sbyte2str(args)
	if strings.HasPrefix(filePath, "/") {
		filePath = path.Join(server.Home, filePath)
	} else {
		filePath = path.Join(server.Cwd, filePath)
	}

	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	fileInfo, err := f.Stat()
	if err != nil {
		return err
	}
	// TODO 暂不支持下载文件夹
	if fileInfo.IsDir() {
		return server.Write(0)
	}

	err = server.Write(fileInfo.Size())
	if err != nil {
		return err
	}
	bufReader := bufio.NewReader(f)
	buf := make([]byte, 4096)
	for {
		n, err := bufReader.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return err
			}
		}
		err = server.Write(buf[:n])
		if err != nil {
			return err
		}
	}

	return server.WriteContent(common.Str2sbyte(common.NoErr.Error()))
}
