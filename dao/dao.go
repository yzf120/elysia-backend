package dao

import (
	"database/sql"
	"fmt"
	"github.com/yzf120/elysia-backend/config"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

// DB 数据库连接实例
var DB *sql.DB

// InitDB 初始化数据库连接
func InitDB() error {
	// 加载配置
	cfg := config.LoadConfig()

	// 连接数据库
	var err error
	DB, err = sql.Open("mysql", cfg.GetDSN())
	if err != nil {
		return fmt.Errorf("数据库连接失败: %v", err)
	}

	// 设置连接池参数
	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(5)
	DB.SetConnMaxLifetime(5 * 60) // 5分钟

	// 测试连接
	err = DB.Ping()
	if err != nil {
		return fmt.Errorf("数据库连接测试失败: %v", err)
	}

	log.Printf("数据库连接成功: %s@%s:%s/%s",
		cfg.Database.User, cfg.Database.Host, cfg.Database.Port, cfg.Database.Name)
	return nil
}

// CloseDB 关闭数据库连接
func CloseDB() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}

// User 用户模型
type User struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// CreateUser 创建用户
func CreateUser(user *User) error {
	query := `INSERT INTO users (username, email) VALUES (?, ?)`
	result, err := DB.Exec(query, user.Username, user.Email)
	if err != nil {
		return fmt.Errorf("创建用户失败: %v", err)
	}

	user.ID, err = result.LastInsertId()
	if err != nil {
		return fmt.Errorf("获取用户ID失败: %v", err)
	}

	return nil
}

// GetUserByID 根据ID获取用户
func GetUserByID(id int64) (*User, error) {
	query := `SELECT id, username, email, created_at, updated_at FROM users WHERE id = ?`
	row := DB.QueryRow(query, id)

	var user User
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("查询用户失败: %v", err)
	}

	return &user, nil
}

// GetUserByUsername 根据用户名获取用户
func GetUserByUsername(username string) (*User, error) {
	query := `SELECT id, username, email, created_at, updated_at FROM users WHERE username = ?`
	row := DB.QueryRow(query, username)

	var user User
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("查询用户失败: %v", err)
	}

	return &user, nil
}
