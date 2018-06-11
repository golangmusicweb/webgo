package main

import (
	//"fmt"
	"webgo/restful"
	"webgo/entity"
	"fmt"
	"time"
)

func main() {
	// Add a logger
	var logger restful.Logging
	logger.GetLogger()
	defer logger.Close()

	// Load config
	var config restful.Config
	config.LoadConfig()

	// get orm engine
	dbEngine := restful.GetDbEngine("default")
	dbEngine.TZLocation, _ = time.LoadLocation("Asia/Shanghai")
	user := new(entity.UserProfile)
	/*
	user.Username = "dongxy"
	user.Password = "sk927312*"
	user.Email = "870428371@qq.com"
	user.Address = "beijing"
	birthday, _ := time.Parse("2006-01-02", "1992-03-13")
	user.Birthday = birthday
	user.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Println(user)
	affected, err := dbEngine.Insert(user)
	fmt.Println(affected)
	*/

	total, err := dbEngine.Where("id >?", 1).Count(user)
	fmt.Println(total, err)
}
