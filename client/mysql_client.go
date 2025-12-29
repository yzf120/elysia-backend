package client

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/yzf120/elysia-backend/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// MySQLClient MySQL客户端
type MySQLClient struct {
	DB     *sql.DB
	GormDB *gorm.DB
}

var defaultMySQLClient *MySQLClient

// InitMySQLClient 初始化MySQL客户端
func InitMySQLClient() error {
	cfg := config.LoadConfig()

	// 初始化原生 sql.DB
	db, err := sql.Open("mysql", cfg.GetDSN())
	if err != nil {
		return fmt.Errorf("数据库连接失败: %v", err)
	}

	// 设置连接池参数
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// 测试连接
	if err := db.Ping(); err != nil {
		return fmt.Errorf("数据库连接测试失败: %v", err)
	}

	// 初始化 GORM
	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn: db,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return fmt.Errorf("GORM 初始化失败: %v", err)
	}

	defaultMySQLClient = &MySQLClient{
		DB:     db,
		GormDB: gormDB,
	}

	log.Printf("MySQL客户端初始化成功: %s@%s:%s/%s",
		cfg.Database.User, cfg.Database.Host, cfg.Database.Port, cfg.Database.Name)
	return nil
}

// GetMySQLClient 获取默认MySQL客户端
func GetMySQLClient() *MySQLClient {
	return defaultMySQLClient
}

// Close 关闭数据库连接
func (c *MySQLClient) Close() error {
	if c.DB != nil {
		return c.DB.Close()
	}
	return nil
}

// QueryRow 执行单行查询
func (c *MySQLClient) QueryRow(query string, args ...interface{}) *sql.Row {
	return c.DB.QueryRow(query, args...)
}

// Query 执行多行查询
func (c *MySQLClient) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return c.DB.Query(query, args...)
}

// Exec 执行SQL语句
func (c *MySQLClient) Exec(query string, args ...interface{}) (sql.Result, error) {
	return c.DB.Exec(query, args...)
}

// Begin 开启事务
func (c *MySQLClient) Begin() (*sql.Tx, error) {
	return c.DB.Begin()
}
