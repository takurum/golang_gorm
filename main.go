package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"os"
)

type Channel struct {
	gorm.Model
	Name        string
	Description string
}

type User struct {
	gorm.Model
	Email    string
	Username string
}

type Message struct {
	gorm.Model
	Content   string
	UserID    uint
	ChannelID uint
	User      User
	Channel   Channel
}

func getDB() *gorm.DB {
	os.Getenv("REDIS_DB")
	args := fmt.Sprintf(
		"%s:%s@tcp(127.0.0.1:5576)/%s?charset=utf8&parseTime=True",
		os.Getenv("MYSQL_USER"),
		os.Getenv("MYSQL_USER_PASS"),
		os.Getenv("MYSQL_NAME"),
	)

	db, err := gorm.Open("mysql", args)
	if err != nil {
		panic(err.Error())
	}

	return db
}

func createDB(db *gorm.DB) {
	db.AutoMigrate(&Channel{}, &User{}, &Message{})
	db.Model(&Message{}).AddForeignKey("user_id", "users(id)", "RESTRICT", "RESTRICT")
}

func fillDB(db *gorm.DB) {
	channels := []Channel{
		{Name: "Home", Description: "Home discussions"},
		{Name: "Work", Description: "Off discussions"},
	}
	for _, c := range channels {
		db.Create(&c)
	}

	users := []User{
		{Email: "fred@mail.com", Username: "Freddy"},
		{Email: "lina@mail.com", Username: "Elina"},
	}
	for _, u := range users {
		db.Create(&u)
	}

	var homeChat, workChat Channel
	db.First(&homeChat, "Name=?", "Home")
	db.First(&workChat, "Name=?", "Work")

	var joe, bob User
	db.First(&joe, "Username=?", "Freddy")
	db.First(&bob, "Username=?", "Elina")

	messages := []Message{
		{Content: "Thank you!", Channel: homeChat, User: joe},
		{Content: "How are you?", Channel: homeChat, User: bob},
		{Content: "What about smt?", Channel: workChat, User: joe},
	}

	for _, m := range messages {
		db.Create(&m)
	}
}

func main() {
	db := getDB()
	defer db.Close()

	createDB(db)
	fillDB(db)

	var users []User
	db.Find(&users)
	for _, u := range users {
		fmt.Println("Email: ", u.Email, "Username: ", u.Username)
	}

	var messages []Message
	db.Model(users[0]).Related(&messages)
	for _, m := range messages {
		fmt.Println("Message:", m.Content, "Sender:", m.UserID)
	}
}
