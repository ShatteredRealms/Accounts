package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

func main() {
	dsn := "admin:MJrLs5z^kTS4!3AZHfH2@tcp(sro-db-1.cluster-cmtxtfoc35iu.us-east-1.rds.amazonaws.com:3306)/accounts?charset=utf8&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Printf("Error: %v", err)
		return
	}

	replicas := []gorm.Dialector{mysql.Open("admin:MJrLs5z^kTS4!3AZHfH2@tcp(sro-db-1.cluster-ro-cmtxtfoc35iu.us-east-1.rds.amazonaws.com:3306)/accounts?charset=utf8&parseTime=True&loc=Local")}

	err = db.Use(dbresolver.Register(dbresolver.Config{
		Replicas: replicas,
		Policy:   dbresolver.RandomPolicy{},
	}))

	if err != nil {
		fmt.Printf("Error: %v", err)
		return
	}

	if db != nil {
		fmt.Printf("Success\n")
	}

	db.AutoMigrate(&Test{})
	db.Create(&Test{Test: "asdf"})

}

type Test struct {
	gorm.Model
	Test string
}
