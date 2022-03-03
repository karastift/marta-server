package main

import (
	"encoding/json"
	"fmt"
)

type ClientInfo struct {
	Username      string
	Name          string
	Uid           string
	HomeDir       string
	Os            string
	Device        string
	MacAdress     string
	Administrator bool
	LocalAddress  string
}

func NewClientInfo(jsonStr string) (*ClientInfo, error) {
	info := ClientInfo{}

	err := json.Unmarshal([]byte(jsonStr), &info)

	return &info, err
}

func (info *ClientInfo) String() string {
	return fmt.Sprintf(`Username	%s
Name		%s
Uid		%s
HomeDir		%s
Os		%s
Device		%s
MacAdress	%s
Administrator	%t
LocalAddress	%s`, info.Username, info.Name, info.Uid, info.HomeDir, info.Os, info.Device, info.MacAdress, info.Administrator, info.LocalAddress)
}
