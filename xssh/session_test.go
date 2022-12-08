package xssh

import "testing"

func TestRemoteSession_Exists(t *testing.T) {
	session := LocalSession{}
	output, err := session.Output("cmd", "/c", "dir", "/b", "D:\\iso\\H3C_UIS-E0770-centos-x86_64-AUTO.iso")
	if err != nil {
		t.Fatal(err)
	}
	if len(output) != 0 {
		t.Log("true")
	} else {
		t.Log("false")
	}
}
