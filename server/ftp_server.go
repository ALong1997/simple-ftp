package main

import (
	"SimpleFTP/common"
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"path"
	"strings"
)

type FTPServer struct {
	*common.FtpConn
}

// New FTPServer
func NewServer(con net.Conn, cwd, home string) *FTPServer {
	return &FTPServer{
		&common.FtpConn{
			Con:  con,
			Cwd:  cwd,
			Home: home,
			Exit: false,
		},
	}
}

// Close Server
func (s *FTPServer) Close() error {
	return s.Con.Close()
}

// Handle cd
func (s *FTPServer) HandleCd(args []byte) error {
	cwd := common.Sbyte2str(args)
	if strings.HasPrefix(cwd, "/") {
		cwd = path.Join(s.Cwd, cwd)
	}

	f, err := os.Open(cwd)
	if err != nil {
		return s.WriteContent(common.Str2sbyte(err.Error()))
	}
	defer f.Close()

	fileInfo, err := f.Stat()
	if err != nil {
		return s.WriteContent(common.Str2sbyte(err.Error()))
	}
	if !fileInfo.IsDir() {
		return s.WriteContent(common.Str2sbyte(common.CDArgsErr.Error()))
	}

	s.Cwd = cwd

	return s.WriteContent(common.Str2sbyte(cwd))
}

// Handle ls
func (s *FTPServer) HandleLs(args []byte) error {
	cwd := common.Sbyte2str(args)
	if strings.HasPrefix(cwd, "/") {
		cwd = path.Join(s.Cwd, cwd)
	}
	if len(cwd) == 0 {
		cwd = s.Cwd
	}

	f, err := os.Open(cwd)
	if err != nil {
		return s.WriteContent(common.Str2sbyte(err.Error()))
	}

	fileInfo, err := f.Stat()
	if err != nil {
		return s.WriteContent(common.Str2sbyte(err.Error()))
	}
	if fileInfo.IsDir() {
		fileInfos, err := f.Readdir(0)
		if err != nil {
			s.WriteContent(common.Str2sbyte(err.Error()))
		}
		var res string
		res = fmt.Sprintf("Total:%d\n", len(fileInfos))
		for _, info := range fileInfos {
			res = res + fmt.Sprintf("%.30s\t%.10d\t%s\n", info.Name(), info.Size(), info.ModTime())
		}
		err = s.WriteContent(common.Str2sbyte(res))
	} else {
		res := fmt.Sprintf("%.30s\t%.10d\t%s\n", fileInfo.Name(), fileInfo.Size(), fileInfo.ModTime())
		err = s.WriteContent(common.Str2sbyte(res))
	}

	if err != nil {
		return s.WriteContent(common.Str2sbyte(err.Error()))
	}
	return nil
}

// Handle exit
func (s *FTPServer) HandleExit(args []byte) error {
	s.Exit = true
	return s.WriteContent(common.Str2sbyte("Byebye."))
}

// Handle mkdir
func (s *FTPServer) HandleMkdir(args []byte) error {
	dir := common.Sbyte2str(args)
	if strings.HasPrefix(dir, "/") {
		dir = path.Join(s.Home, dir)
	} else {
		dir = path.Join(s.Cwd, dir)
	}

	err := os.Mkdir(dir, os.ModePerm)
	if err != nil {
		return err
	}
	return s.WriteContent(common.Str2sbyte(common.NoErr.Error()))
}

// Handle put
func (s *FTPServer) HandlePut(args []byte) error {
	fileName := common.Sbyte2str(args)
	f, err := os.Create(path.Join(s.Cwd, fileName))
	if err != nil {
		return err
	}
	defer f.Close()

	var length int64
	err = s.Read(&length)
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
		err = s.Read(buf)
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

	return s.WriteContent(common.Str2sbyte(common.NoErr.Error()))
}

// Handle get
func (s *FTPServer) HandleGet(args []byte) error {
	filePath := common.Sbyte2str(args)
	if strings.HasPrefix(filePath, "/") {
		filePath = path.Join(s.Home, filePath)
	} else {
		filePath = path.Join(s.Cwd, filePath)
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
		return s.Write(0)
	}

	err = s.Write(fileInfo.Size())
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
		err = s.Write(buf[:n])
		if err != nil {
			return err
		}
	}

	return s.WriteContent(common.Str2sbyte(common.NoErr.Error()))
}
