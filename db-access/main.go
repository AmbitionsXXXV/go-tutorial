package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/go-sql-driver/mysql"
)

var db *sql.DB

// Album 结构体定义，用于表示数据库中的专辑记录。
type Album struct {
	ID     int64
	Title  string
	Artist string
	Price  float32
}

func main() {
	// 从环境变量中获取数据库连接信息。
	cfg := mysql.Config{
		User:   os.Getenv("DBUSER"), // 数据库用户名
		Passwd: os.Getenv("DBPASS"), // 数据库密码
		Net:    "tcp",
		Addr:   "127.0.0.1:3306", // 数据库地址
		DBName: "recordings",     // 数据库名
	}

	// 使用配置信息创建数据库连接。
	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())

	if err != nil {
		log.Fatal(err)
	}

	// 测试数据库连接是否成功。
	pingErr := db.Ping()

	if pingErr != nil {
		log.Fatal(pingErr)
	}

	fmt.Println("Connected!")

	// 查询特定艺术家的专辑。
	albums, err := albumsByArtist("ONE OK ROCK")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Albums found: %v\n", albums)

	// 根据ID查询专辑。
	alb, err := albumByID(6)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Album found: %v\n", alb)

	// 添加一个新专辑。
	albID, err := addAlbum(Album{
		Title:  "Tamago",
		Artist: "Aimyon",
		Price:  49.99,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("ID of added album: %v\n", albID)
}

// albumsByArtist 根据艺术家名称查询专辑。
func albumsByArtist(name string) ([]Album, error) {
	var albums []Album

	// 执行SQL查询。
	rows, err := db.Query("SELECT * FROM recordings.album WHERE artist = ?", name)
	if err != nil {
		return nil, fmt.Errorf("albumByArtist %q: %v", name, err)
	}

	// 保证资源释放
	defer func(rows *sql.Rows) {
		err := rows.Close()

		if err != nil {
			log.Fatal(err)
		}
	}(rows)

	// 遍历查询结果。
	for rows.Next() {
		var alb Album

		if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
			return nil, fmt.Errorf("albumByArtist %q: %v", name, err)
		}

		albums = append(albums, alb)
	}

	// 检查查询过程中是否有错误发生。
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("albumByArtist %q: %v", name, err)
	}

	return albums, nil
}

// addAlbum 添加一个新专辑到数据库。
func addAlbum(alb Album) (int64, error) {
	// 执行SQL插入操作。
	result, err := db.Exec("INSERT INTO recordings.album (title, artist, price) VALUES (?, ?, ?)", alb.Title, alb.Artist, alb.Price)

	if err != nil {
		return 0, fmt.Errorf("addAlbum: %v", err)
	}

	// 获取插入记录的ID。
	id, err := result.LastInsertId()

	if err != nil {
		return 0, fmt.Errorf("addAlbum: %v", err)
	}

	return id, nil
}

// albumByID 根据ID查询一个专辑。
func albumByID(id int64) (Album, error) {
	var alb Album

	// 执行SQL查询。
	row := db.QueryRow("SELECT * FROM album WHERE id = ?", id)

	if err := row.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return alb, fmt.Errorf("albumsById %d: no such album", id)
		}

		return alb, fmt.Errorf("albumsById %d: %v", id, err)
	}

	return alb, nil
}
