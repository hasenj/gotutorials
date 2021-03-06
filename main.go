package main

import "os"

import "time"

// Some example types to test the go-ts type bridge

type UserLoginInfo struct {
	UserID   int
	Name     string
	IsActive bool
}

type UserBasicInfo struct {
	Email string
	Bio   string
}

type UserProfile struct {
	UserLoginInfo
	Basic     UserBasicInfo
	LastLogin time.Time `ts:"string"`
	Friends   	 []UserLoginInfo
}

func main() {
	var bridge = ts_bridge.NewTypeBridge()
	var inst UserProfile
	ts_bridge.QueueInstance(bridge, inst)
	ts_bridge.Process(bridge)
	ts_bridge.DescribeTypes(bridge, os.Stdout)
}
