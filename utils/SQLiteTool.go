package utils

import (
	"database/sql"
	"fmt"
	"io"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

// ReadSQLiteDB 读取 SQLite 数据库文件并执行 SQL 查询
func ReadSQLiteDB_url(dbPath string, query string) (string, error) {
	// 创建一个 strings.Builder 对象，用于构建最终的字符串结果
	var builder strings.Builder

	// 创建一个临时文件
	tempFile, err := os.CreateTemp("", "chrome-history-*.db")
	if err != nil {
		return "", err
	}
	defer os.Remove(tempFile.Name()) // 确保临时文件在函数结束时被删除

	// 将数据库文件复制到临时文件
	if err := copyFile(dbPath, tempFile.Name()); err != nil {
		return "", err
	}

	// 打开 SQLite 数据库
	db, err := sql.Open("sqlite3", tempFile.Name())
	if err != nil {
		return "", err
	}
	defer db.Close()

	// 执行 SQL 查询
	rows, err := db.Query(query)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	// 遍历查询结果
	for rows.Next() {

		var url string
		if err := rows.Scan(&url); err != nil {
			return "", err
		}
		// 将结果添加到 strings.Builder 中
		builder.WriteString(url + "\n")
	}

	// 返回构建的字符串结果
	return builder.String(), nil
}

// ReadSQLiteDB 读取 SQLite 数据库并返回查询结果
func ReadSQLiteDB(dbPath string, query string) (string, error) {
	// 创建一个 strings.Builder 对象，用于构建最终的字符串结果
	var builder strings.Builder

	// 创建一个临时文件
	tempFile, err := os.CreateTemp("", "chrome-history-*.db")
	if err != nil {
		return "", err
	}
	defer os.Remove(tempFile.Name()) // 确保临时文件在函数结束时被删除

	// 将数据库文件复制到临时文件
	if err := copyFile(dbPath, tempFile.Name()); err != nil {
		return "", err
	}

	// 打开 SQLite 数据库
	db, err := sql.Open("sqlite3", tempFile.Name())
	if err != nil {
		return "", err
	}
	defer db.Close()

	// 执行 SQL 查询
	rows, err := db.Query(query)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	// 获取查询结果的列数
	columns, err := rows.Columns()
	if err != nil {
		return "", err
	}

	// 预分配切片以存储每一行的数据
	values := make([]interface{}, len(columns))
	valuePtrs := make([]interface{}, len(columns))

	// 遍历查询结果
	for rows.Next() {
		for i := range columns {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return "", err
		}

		// 将结果添加到 strings.Builder 中
		for i, value := range values {
			switch v := value.(type) {
			case []byte:
				builder.WriteString(string(v))
			case string:
				builder.WriteString(v)
			default:
				builder.WriteString(fmt.Sprintf("%v", v))
			}
			if i < len(values)-1 {
				builder.WriteString("\t") // 使用制表符分隔字段
			}
		}
		builder.WriteString("\n")
	}

	// 返回构建的字符串结果
	return builder.String(), nil
}

// 辅助函数：复制文件
func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}
