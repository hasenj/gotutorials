package main

import "os"
import "time"
import . "github.com/hasenj/gotutorials/ts_bridge"

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
	var bridge = NewTypeBridge()
	var inst UserProfile
	QueueInstance(bridge, inst)
	Process(bridge)
	DescribeTypes(bridge, os.Stdout)
}
