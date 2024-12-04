package utils

import (
	"database/sql"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

// contains 函数检查字符串数组中是否包含特定的字符串
func contains(arr []int, str int) bool {
	for _, v := range arr {
		if v == str {
			return true
		}
	}
	return false
}

// 读取 SQLite 数据库文件并执行 SQL 查询
func ReadSQLiteDB_url(dbPath string, query string) (string, error) {
	// 创建一个 strings.Builder 对象，用于构建最终的字符串结果
	var builder strings.Builder

	// 创建一个临时文件
	tempFile, err := os.CreateTemp("", "temp_sqlite-*.db")
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
func ReadSQLiteDB_url2(dbPath string, query string) (string, []string) {
	// 创建一个 strings.Builder 对象，用于构建最终的字符串结果
	var builder strings.Builder

	// 创建一个临时文件
	tempFile, err := os.CreateTemp("", "temp_history-*.db")
	if err != nil {
		return "", nil
	}
	defer os.Remove(tempFile.Name()) // 确保临时文件在函数结束时被删除

	// 将数据库文件复制到临时文件
	if err := copyFile(dbPath, tempFile.Name()); err != nil {
		return "", nil
	}

	// 打开 SQLite 数据库
	db, err := sql.Open("sqlite3", tempFile.Name())
	if err != nil {
		return "", nil
	}

	// 执行 SQL 查询
	rows, err := db.Query(query)
	if err != nil {
		return "", nil
	}
	defer rows.Close()
	var list []int
	// 遍历查询结果
	for rows.Next() {
		var fk sql.NullInt64
		if err := rows.Scan(&fk); err != nil {
			fmt.Println("扫描行失败:", err)
			return "", nil
		}
		if fk.Valid && fk.Int64 != 0 {
			list = append(list, int(fk.Int64))
		}
	}

	rows, err = db.Query("SELECT id, url FROM moz_places ")
	if err != nil {
		return "", nil
	}
	defer rows.Close()
	// 获取查询结果的列数
	columns, err := rows.Columns()
	if err != nil {
		return "", nil
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
			return "", nil
		}
		var id string = "0"
		for i, value := range values {
			var fieldValue string
			switch v := value.(type) {
			case int:
				fieldValue = strconv.Itoa(v)
			case string:
				fieldValue = v
			default:
				fieldValue = fmt.Sprintf("%v", v)
			}
			if columns[i] == "id" {
				num, _ := strconv.Atoi(fieldValue)
				if contains(list, num) {
					//fmt.Printf("The array contains '%s'.\n", fieldValue)
					id = fieldValue
					//list2 = append(list2, )
				}
			}
			if columns[i] == "url" {
				if id != "0" {
					builder.WriteString(fieldValue + "\n")
				}
			}

		}

		// 将结果添加到 strings.Builder 中
		//builder.WriteString(url + "\n")
	}

	// 返回构建的字符串结果
	return builder.String(), nil
}

// ReadSQLiteDB 读取 SQLite 数据库并返回查询结果
func ReadSQLiteDB(dbPath string, query string, CSVData [][]string) ([][]string, error) {

	// 创建一个临时文件
	tempFile, err := os.CreateTemp("", "chrome-history-*.db")
	if err != nil {
		return nil, err
	}
	defer os.Remove(tempFile.Name()) // 确保临时文件在函数结束时被删除

	// 将数据库文件复制到临时文件
	if err := copyFile(dbPath, tempFile.Name()); err != nil {
		return nil, err
	}

	// 打开 SQLite 数据库
	db, err := sql.Open("sqlite3", tempFile.Name())
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// 执行 SQL 查询
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// 获取查询结果的列数
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
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
			return nil, nil
		}
		var host string
		var name string
		var cookieValue string
		// 将结果添加到 strings.Builder 中
		for i, value := range values {
			var fieldValue string
			switch v := value.(type) {
			case string:
				fieldValue = v

			}
			if columns[i] == "host" {
				host = fieldValue
			}
			if columns[i] == "name" {
				name = fieldValue
			}
			if columns[i] == "value" {
				cookieValue = fieldValue
			}
		}

		CSVData = append(CSVData, []string{host, name, cookieValue})

	}

	// 返回构建的字符串结果
	return CSVData, nil
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
