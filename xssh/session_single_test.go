package xssh

import (
	"testing"
)

func TestSingleSession_Exists(t *testing.T) {
	singleSession, err := NewSingleSession(Config{})
	if err != nil {
		t.Fatal(err)
	}
	exists, err := singleSession.Exists("D:\\iso\\H3C_UIS-E0770-centos-x86_64-AUTO.iso")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(exists)
	exists, err = singleSession.Exists("D:\\iso\\H3C_UIS-E0770-centos-x86_64-AUTO")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(exists)
}
