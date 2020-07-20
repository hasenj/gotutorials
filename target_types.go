package main

import "time"

type UserLoginInfo struct {
	UserID int
	Name string
	IsActive bool
}

type UserProfile struct {
	LoginInfo UserLoginInfo
	LastLogin time.Time `ts:"string"`
	Bio string
}

type UserBasicInfo struct {
	UserLoginInfo
	Email string
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
