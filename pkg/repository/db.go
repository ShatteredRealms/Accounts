package repository

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"io/ioutil"
	"log"
)

// Connection Information on how to connect to the MySQL database
type Connection struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Name     string `yaml:"name"`
	Username string `yaml:"user"`
	Password string `yaml:"pass"`
}

// DBConnect Initializes the connection to the database
func DBConnect(filePath string) (*gorm.DB, error) {
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatalf("reading db file: %v", err)
	}

	c := &Connection{}
	err = yaml.Unmarshal(file, c)
	if err != nil {
		log.Fatalf("yaml: %v", err)
	}

	// refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		c.Username,
		c.Password,
		c.Host,
		c.Port,
		c.Name,
	)

	return gorm.Open(mysql.Open(dsn), &gorm.Config{})
}

func (c Connection) DSN() string {
	// refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		c.Username,
		c.Password,
		c.Host,
		c.Port,
		c.Name,
	)
}
