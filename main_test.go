package main

import "testing"

func TestCheckIPIsListedFunc(t *testing.T) {
	ips, isListed := IPIsListed("127.0.0.2")
	if isListed != true {
		t.Errorf("Expected blacklist to be true, got : %v, ips : %v", isListed, ips)
	}
}

func TestCheckGetReverseIPFunc(t *testing.T) {
	ip, _ := GetReverseIP("1.2.3.4")
	if ip != "4.3.2.1" {
		t.Errorf("Incorrect reverse IP, got : %s", ip)
	}
}
