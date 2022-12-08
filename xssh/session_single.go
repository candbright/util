package xssh

import (
	"os"
	"sync"
)

type SingleSession struct {
	Lock    *sync.Mutex
	session Session
}

func NewSingleSession(config Config) (Session, error) {
	session := &SingleSession{
		Lock: &sync.Mutex{},
	}
	if config.IsLocal() {
		session.session = &LocalSession{}
	} else {
		session.session = &RemoteSession{}
	}
	err := session.Connect()
	if err != nil {
		return nil, err
	}
	return session, nil
}

func (c SingleSession) Connect() error {
	return c.session.Connect()
}

func (c SingleSession) Close() error {
	return c.session.Close()
}

func (c SingleSession) IsLocal() bool {
	return c.session.IsLocal()
}

func (c SingleSession) Run(name string, arg ...string) error {
	c.Lock.Lock()
	defer c.Lock.Unlock()
	return c.session.Run(name, arg...)
}

func (c SingleSession) Output(name string, arg ...string) ([]byte, error) {
	c.Lock.Lock()
	defer c.Lock.Unlock()
	return c.session.Output(name, arg...)
}

func (c SingleSession) CombinedOutput(name string, arg ...string) ([]byte, error) {
	c.Lock.Lock()
	defer c.Lock.Unlock()
	return c.session.CombinedOutput(name, arg...)
}

func (c SingleSession) OutputGrep(cmdList []struct {
	name string
	arg  []string
}) ([]byte, error) {
	c.Lock.Lock()
	defer c.Lock.Unlock()
	return c.session.OutputGrep(cmdList)
}

func (c SingleSession) Exists(path string) (bool, error) {
	c.Lock.Lock()
	defer c.Lock.Unlock()
	return c.session.Exists(path)
}

func (c SingleSession) ReadFile(fileName string) ([]byte, error) {
	c.Lock.Lock()
	defer c.Lock.Unlock()
	return c.session.ReadFile(fileName)
}

func (c SingleSession) ReadDir(dir string) ([]FileInfo, error) {
	c.Lock.Lock()
	defer c.Lock.Unlock()
	return c.session.ReadDir(dir)
}

func (c SingleSession) MakeDirAll(path string, perm os.FileMode) error {
	c.Lock.Lock()
	defer c.Lock.Unlock()
	return c.session.MakeDirAll(path, perm)
}

func (c SingleSession) Remove(name string) error {
	c.Lock.Lock()
	defer c.Lock.Unlock()
	return c.session.Remove(name)
}

func (c SingleSession) RemoveAll(path string) error {
	c.Lock.Lock()
	defer c.Lock.Unlock()
	return c.session.RemoveAll(path)
}

func (c SingleSession) Create(name string) error {
	c.Lock.Lock()
	defer c.Lock.Unlock()
	return c.session.Create(name)
}

func (c SingleSession) WriteString(name string, data string, mode ...string) error {
	c.Lock.Lock()
	defer c.Lock.Unlock()
	return c.session.WriteString(name, data, mode...)
}
