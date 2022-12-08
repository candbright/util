package xssh

import (
	"bytes"
	"errors"
	"fmt"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var (
	ErrNilSshClient = errors.New("ssh client is nil")
)

type Session interface {
	IsLocal() bool
	Connect() error
	Close() error
	Run(name string, arg ...string) error
	Output(name string, arg ...string) ([]byte, error)
	CombinedOutput(name string, arg ...string) ([]byte, error)
	OutputGrep(cmdList []struct {
		name string
		arg  []string
	}) ([]byte, error)
	Exists(path string) (bool, error)
	ReadFile(fileName string) ([]byte, error)
	ReadDir(dir string) ([]FileInfo, error)
	MakeDirAll(path string, perm os.FileMode) error
	Remove(name string) error
	RemoveAll(path string) error
	Create(name string) error
	WriteString(name string, data string, mode ...string) error
}

func NewSession(config Config) (Session, error) {
	var session Session
	if config.IsLocal() {
		session = &LocalSession{}
	} else {
		session = &RemoteSession{}
	}
	err := session.Connect()
	if err != nil {
		return nil, err
	}
	return session, nil
}

type LocalSession struct {
}

func (s *LocalSession) IsLinux() bool {
	return runtime.GOOS == "linux"
}

func (s *LocalSession) IsLocal() bool {
	return true
}

func (s *LocalSession) Connect() error {
	return nil
}

func (s *LocalSession) Close() error {
	return nil
}

func (s *LocalSession) Run(name string, arg ...string) error {
	if s.IsLinux() {
		return exec.Command(name, arg...).Run()
	} else {
		args := make([]string, len(arg)+5)
		args[0] = "/c"
		args[1] = name
		copy(args[2:], arg)
		return exec.Command("cmd", args...).Run()
	}
}

func (s *LocalSession) Output(name string, arg ...string) ([]byte, error) {
	if s.IsLinux() {
		return exec.Command(name, arg...).Output()
	} else {
		args := make([]string, len(arg)+5)
		args[0] = "/c"
		args[1] = name
		copy(args[2:], arg)
		return exec.Command("cmd", args...).Output()
	}
}

func (s *LocalSession) CombinedOutput(name string, arg ...string) ([]byte, error) {
	if s.IsLinux() {
		return exec.Command(name, arg...).CombinedOutput()
	} else {
		args := make([]string, len(arg)+5)
		args[0] = "/c"
		args[1] = name
		copy(args[2:], arg)
		return exec.Command("cmd", args...).CombinedOutput()
	}
}

func (s *LocalSession) OutputGrep(cmdList []struct {
	name string
	arg  []string
}) ([]byte, error) {
	if cmdList == nil {
		return nil, errors.New("execute cmd grep failed, cmd list is nil")
	}
	cmdStrList := make([]string, len(cmdList))
	for i, cmd := range cmdList {
		cmdStrList[i] = cmd.name
		for _, arg := range cmd.arg {
			cmdStrList[i] += " " + arg
		}
	}
	if s.IsLinux() {
		return exec.Command("cmd", "/c", strings.Join(cmdStrList, " | ")).Output()
	} else {
		return exec.Command("bash", "-c", strings.Join(cmdStrList, " | ")).Output()
	}
}

func (s *LocalSession) Exists(path string) (bool, error) {
	return Exists(path), nil
}

func (s *LocalSession) ReadFile(fileName string) ([]byte, error) {
	return ioutil.ReadFile(fileName)
}

func (s *LocalSession) ReadDir(dir string) ([]FileInfo, error) {
	dirs, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	files := make([]FileInfo, len(dirs))
	for i, fileInfo := range dirs {
		files[i] = FileInfo{Name: fileInfo.Name(), Path: fileInfo.Name()}
	}
	return files, nil
}

func (s *LocalSession) MakeDirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}

func (s *LocalSession) Remove(name string) error {
	return os.Remove(name)
}

func (s *LocalSession) RemoveAll(path string) error {
	return os.RemoveAll(path)
}

func (s *LocalSession) Create(name string) error {
	_, err := os.Create(name)
	return err
}

func (s *LocalSession) WriteString(name string, data string, mode ...string) error {
	flag := os.O_RDWR | os.O_CREATE | os.O_TRUNC
	if len(mode) == 1 && mode[0] == ">>" {
		flag = os.O_RDWR | os.O_CREATE | os.O_APPEND
	}
	fileHandler, err := os.OpenFile(name, flag, os.ModePerm)
	if err != nil {
		return err
	}
	_, err = fileHandler.WriteString(data)
	if err != nil {
		return err
	}
	return nil
}

type RemoteSession struct {
	Config
	Client *ssh.Client
}

func (s *RemoteSession) IsLinux() bool {
	err := s.Run("cat", "/etc/os-release")
	if err != nil {
		return false
	}
	return true
}

func (s *RemoteSession) IsLocal() bool {
	return false
}

func (s *RemoteSession) Connect() error {
	if s.Client != nil {
		err := s.Client.Close()
		if err != nil {
			return err
		}
	}
	var auth []ssh.AuthMethod
	if s.Config.Password() != "" {
		auth = []ssh.AuthMethod{ssh.Password(s.Config.Password())}
	} else {
		sshKeyPath := "/root/.ssh/id_rsa"
		keyAuth, err := publicKeyAuth(sshKeyPath)
		if err != nil {
			return err
		}
		auth = []ssh.AuthMethod{keyAuth}
	}
	sshCfg := &ssh.ClientConfig{
		Timeout:         time.Second * 3,
		User:            s.Config.User(),
		Auth:            auth,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	addr := fmt.Sprintf("%s:%d", s.Config.Host(), s.Config.Port())
	session, err := ssh.Dial("tcp", addr, sshCfg)
	if err != nil {
		return err
	}
	s.Client = session
	return nil
}

func publicKeyAuth(kPath string) (ssh.AuthMethod, error) {
	key, err := ioutil.ReadFile(kPath)
	if err != nil {
		return nil, err
	}
	singer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, err
	}
	return ssh.PublicKeys(singer), nil
}

func (s *RemoteSession) Close() error {
	if s.Client != nil {
		return s.Client.Close()
	}
	return nil
}

func (s *RemoteSession) Run(name string, arg ...string) error {
	if s.Client == nil {
		return ErrNilSshClient
	}
	session, err := s.Client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()
	session.Stderr = &bytes.Buffer{}
	err = session.Run(Command(name, arg...))
	if err != nil {
		return err
	}
	return nil
}

func (s *RemoteSession) Output(name string, arg ...string) ([]byte, error) {
	if s.Client == nil {
		return nil, ErrNilSshClient
	}
	session, err := s.Client.NewSession()
	if err != nil {
		return nil, err
	}
	defer session.Close()
	session.Stderr = &bytes.Buffer{}
	output, err := session.Output(Command(name, arg...))
	if err != nil {
		return nil, err
	}
	return output, nil
}

func (s *RemoteSession) CombinedOutput(name string, arg ...string) ([]byte, error) {
	if s.Client == nil {
		return nil, ErrNilSshClient
	}
	session, err := s.Client.NewSession()
	if err != nil {
		return nil, err
	}
	defer session.Close()
	output, err := session.CombinedOutput(Command(name, arg...))
	if err != nil {
		return nil, err
	}
	return output, nil
}

func (s *RemoteSession) OutputGrep(cmdList []struct {
	name string
	arg  []string
}) ([]byte, error) {
	if cmdList == nil {
		return nil, errors.New("execute cmd grep failed, cmd list is nil")
	}
	if s.IsLinux() {
		cmdStrList := make([]string, len(cmdList))
		for i, cmd := range cmdList {
			cmdStrList[i] = cmd.name
			for _, arg := range cmd.arg {
				cmdStrList[i] += " " + arg
			}
		}
		return s.Output("bash", "-c", "'"+strings.Join(cmdStrList, " | ")+"'")
	} else {
		//TODO
		return nil, nil
	}
}

func (s *RemoteSession) Exists(path string) (bool, error) {
	if s.Client == nil {
		return false, ErrNilSshClient
	}
	var err error
	var output []byte
	if s.IsLinux() {
		if strings.HasSuffix(path, "/") {
			//dir
			output, err = s.Output("find", path)
		} else {
			//file
			output, err = s.Output("find", Dir(path), "-name", FileName(path))
		}
		if err != nil {
			return false, err
		}
		if len(output) != 0 {
			return true, nil
		} else {
			return false, nil
		}
	} else {
		output, err = s.Output("dir", "/b", path)
		if err != nil {
			return false, err
		}
		if len(output) != 0 {
			return true, nil
		} else {
			return false, nil
		}
	}
}

func (s *RemoteSession) ReadFile(fileName string) ([]byte, error) {
	if s.Client == nil {
		return nil, ErrNilSshClient
	}
	if s.IsLinux() {
		return s.Output("cat", fileName)
	} else {
		//TODO
		return nil, nil
	}
}

func (s *RemoteSession) ReadDir(dir string) ([]FileInfo, error) {
	if s.Client == nil {
		return nil, ErrNilSshClient
	}
	if s.IsLinux() {
		output, err := s.Output("ls", "-AF", dir)
		if err != nil {
			return nil, err
		}
		files := make([]FileInfo, 0)
		for _, file := range strings.Split(string(output), "\n") {
			if strings.HasSuffix(file, "/") {
				files = append(files, FileInfo{Name: file, Path: dir + file})
			} else {
				files = append(files, FileInfo{Name: file, Path: dir + "/" + file})
			}
		}
		return files, nil
	} else {
		//TODO
		return nil, nil
	}
}

func (s *RemoteSession) MakeDirAll(path string, perm os.FileMode) error {
	if s.Client == nil {
		return ErrNilSshClient
	}
	if s.IsLinux() {
		return s.Run("mkdir", "-p", path, "-m", strconv.FormatUint(uint64(perm), 10))
	} else {
		//TODO
		return nil
	}
}

func (s *RemoteSession) Remove(name string) error {
	if s.Client == nil {
		return ErrNilSshClient
	}
	if s.IsLinux() {
		return s.Run("rm", "-f", name)
	} else {
		//TODO
		return nil
	}
}

func (s *RemoteSession) RemoveAll(path string) error {
	if s.Client == nil {
		return ErrNilSshClient
	}
	if s.IsLinux() {
		return s.Run("rm", "-r", "-f", path)
	} else {
		//TODO
		return nil
	}
}

func (s *RemoteSession) Create(name string) error {
	if s.Client == nil {
		return ErrNilSshClient
	}
	if s.IsLinux() {
		exists, err := s.Exists(Dir(name))
		if err != nil {
			return err
		}
		if !exists {
			err = s.MakeDirAll(Dir(name), 0666)
			if err != nil {
				return err
			}
		}
		return s.Run("touch", name)
	} else {
		//TODO
		return nil
	}
}

func (s *RemoteSession) WriteString(name string, data string, mode ...string) error {
	if s.Client == nil {
		return ErrNilSshClient
	}
	if s.IsLinux() {
		flag := ">"
		if len(mode) == 1 && mode[0] == ">>" {
			flag = ">>"
		}
		return s.Run("echo", fmt.Sprintf(`"%s"`, data), flag)
	} else {
		//TODO
		return nil
	}
}
