package repository

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

// Connection Information on how to connect to the MySQL database
type Connection struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Name     string `yaml:"name"`
	Username string `yaml:"user"`
	Password string `yaml:"pass"`
}

// DBConnections Defines the master and slave connections to a replicated database. Slaves may be empty.
type DBConnections struct {
	Master Connection   `yaml:"master"`
	Slaves []Connection `yaml:"slaves"`
}

// DBConnect Initializes the connection to the database
func DBConnect(connections DBConnections) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(connections.Master.DSN()), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	replicas := make([]gorm.Dialector, len(connections.Slaves))
	for _, slave := range connections.Slaves {
		replicas = append(replicas, mysql.Open(slave.DSN()))
	}

	err = db.Use(dbresolver.Register(dbresolver.Config{
		Replicas: replicas,
		Policy:   dbresolver.RandomPolicy{},
	}))

	return db, err
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
