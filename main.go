package main

import "os"

func main() {
	var bridge = NewTypeBridge()
	var inst UserProfile
	QueueInstance(bridge, inst)
	Process(bridge)
	DescribeTypes(bridge, os.Stdout)
}
