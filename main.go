package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model   //Isso aqui adiciona toda a estrutura de Model ID, DELETEDAT, ETC
	Name         string
	Email        string
	Age          uint8
	Birthday     time.Time
	MemberNumber sql.NullString
	ActivatedAt  sql.NullTime
	Registers    []Register // Associação com registros de usuários
}

type Register struct {
	gorm.Model
	UserID     uint // Chave estrangeira para a tabela 'users'
	User       User `gorm:"foreignKey:UserID"` // Associação com o usuário correspondente
	BirthLocal string
	//Aqui o embedded basicamente faz herdar o User, e o Prefix é para adicionar User antes
	//Aqui a Doc de configuração: https://gorm.io/docs/models.html#embedded_struct
}

func main() {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Register{})

	// Create
	newUser := db.Create(&User{Name: "Juninho", Email: "brizollajr@gmail.com", Age: 20, Birthday: time.Date(2002, 8, 1, 20, 34, 58, 651387237, time.UTC)})

	// Read
	var user User
	db.Find(&user, newUser)
	fmt.Println(&user)

	c := db.Create(&Register{User: user, BirthLocal: "Estância Velha"})
	var r Register
	db.Find(&r, c)
	registerJson, _ := json.MarshalIndent(&r, "", "  ")

	fmt.Printf("Register struct: %+v\n", string(registerJson))

	// Update
	db.Model(&user).Update("Name", "João")

	//Delete Register
	db.Delete(&r, 1)
}
