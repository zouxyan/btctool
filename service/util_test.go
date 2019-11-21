package service

import "testing"

func TestBuildData(t *testing.T) {
	addr := "0x7dD16c0c71F71A123c4BDAF0a468aBC60Db41C0C"
	_, err := buildData(1, 0, addr, addr)
	if err != nil {
		t.Fatal(err)
	}
}
