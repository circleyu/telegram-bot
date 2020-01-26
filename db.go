package main

import (
	"fmt"
	"log"
	"os"

	"github.com/go-xorm/xorm"

	_ "github.com/lib/pq"
)

var (
	host     = "35.221.173.100"
	port     = "5432"
	user     = "postgres"
	password = "shrimp17"
	dbName   = "kua-mei-a"
)

func initConfig() {
	host = os.Getenv("DB_HOST")
	if host == "" {
		log.Panic("DB_HOST is null")
	}
	port = os.Getenv("DB_PORT")
	if host == "" {
		log.Panic("DB_PORT is null")
	}
	user = os.Getenv("DB_USER")
	if host == "" {
		log.Panic("DB_USER is null")
	}
	password = os.Getenv("DB_PASSWORD")
	if host == "" {
		log.Panic("DB_PASSWORD is null")
	}
	dbName = os.Getenv("DB_NAME")
	if host == "" {
		log.Panic("DB_NAME is null")
	}
}

func getDBEngine() *xorm.Engine {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbName)
	//格式
	engine, err := xorm.NewEngine("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	engine.ShowSQL() //菜鸟必备

	err = engine.Ping()
	if err != nil {
		log.Fatal(err)
		return nil
	}
	log.Println("connect postgresql success")
	return engine
}

//TokenTbl table name 为token_tbl
type TokenTbl struct {
	Id    int64
	Name  string
	Token string
}

//SelectToken 条件查询
func SelectToken(name string) *TokenTbl {
	var tokens []TokenTbl
	engine := getDBEngine()
	err := engine.Where("token_tbl.name=?", name).Find(&tokens)
	if err != nil {
		log.Println(err)
		return nil
	} else if len(tokens) == 0 {
		return nil
	}
	return &tokens[0]
}

//InsertToken 添加
func InsertToken(token *TokenTbl) bool {
	engine := getDBEngine()
	rows, err := engine.Insert(token)
	if err != nil {
		log.Println(err)
		return false
	}
	if rows == 0 {
		return false
	}
	return true
}

//DeleteToken 删除(根据名称删除)
func DeleteToken(name string) bool {
	token := TokenTbl{
		Name: name,
	}
	engine := getDBEngine()
	rows, err := engine.Delete(&token)
	if err != nil {
		log.Println(err)
		return false
	}
	if rows == 0 {
		return false
	}
	return true
}
