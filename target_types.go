package main

import "time"

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
	Friends   map[int]*UserLoginInfo
}

/*
	// Want to automatically
	interface UserLoginInfo {
		UserID: number;
		Name: string;
		IsActive: boolean;
	}

	interface UserProfile {
		LoginInfo: UserLoginInfo;
		LastLogin: string;
		Bio: string;
	}

	interface UserBasicInfo {
		UserID: number;
		Name: string;
		IsActive: boolean;

		Email: string;
	}
*/
