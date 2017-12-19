package main

import (
	"github.com/gorilla/mux"
	"github.com/gorilla/rpc"
	"github.com/gorilla/rpc/json"
	"net/http"
	"github.com/satori/go.uuid"
	"time"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"log"
	"flag"
	"errors"
)

var (
	appName = "JSON-RPC User service" // название сервиса
	version = "1.0" // версия
	date    = "2017-12-19" // дата сборки
	host    = ":8080" // адрес сервера и порт
)

type User struct {
	Uuid string `gorm:"primary_key;unique_index;not null"`
	Login string `gorm:"not null;unique_index"`
	RegDate string `gorm:"not null"`
}

/*Создание нового пользователя

Пример запроса:
{"jsonrpc":"2.0","method":"User.Create","id":1,"params":[{"login":"TestUser"}]}

Ответ:
{"result":{"Uuid":"dad5475c-ef25-4682-bf4c-40575d2a1b6a","Login":"TestUser","RegDate":"2017-12-19"},"error":null,"id":1}
*/
func (t *User) Create(r *http.Request, args *User, result *User) error {
	if (args.Login!="") {
		//При создании пользователя ему автоматически выдается uuid и устанавливается дата регистрации.
		newUser := User{Uuid: uuid.NewV4().String(), Login: args.Login, RegDate: time.Now().Local().Format("2006-01-02")}
		if err := db.Create(newUser).Error; err != nil {
			return err
		} else {
			*result = newUser
			return nil
		}
	} else {
		return errors.New("Login should not be empty")
	}
}

/*Получение пользователя (или пользователей)

Пример запроса одного пользователя по uuid:
{"jsonrpc":"2.0","method":"User.Get","id":1,"params":[{"uuid":"dad5475c-ef25-4682-bf4c-40575d2a1b6a"}]}

Ответ:
{"result":[{"Uuid":"dad5475c-ef25-4682-bf4c-40575d2a1b6a","Login":"TestUser","RegDate":"2017-12-19"}],"error":null,"id":1}

Пример запроса пользователей по дате регистрации:
{"jsonrpc":"2.0","method":"User.Get","id":1,"params":[{"regdate":"2017-12-19"}]}

Ответ:
{"result":[{"Uuid":"19a8cfe2-a934-4835-8a47-7046c2487a63","Login":"Paladinhoney","RegDate":"2017-12-19"},{"Uuid":"d69934a0-a305-4567-8b1f-bcd9129702e7","Login":"Crestspark","RegDate":"2017-12-19"}],"error":null,"id":1}
*/
func (t *User) Get(r *http.Request, args *User, result *[]User) error {
	var users []User

	//Если в запросе есть дата регистрации, то выводятся все пользователи с данной датой регистрации.
	if args.RegDate!= "" {
		if err := db.Where("reg_date = ?", args.RegDate).Find(&users).Error; err != nil {
			return err
		} else {
			if len(users) == 0 {
				return errors.New("No users found")
			}
			*result = users
			return nil
		}
	//Если даты регистрации в запросе нет, то выводятся пользователи с соответсвующим uuid или логином.
	} else {
		if err := db.Where("uuid = ?", args.Uuid).Or("login = ?", args.Login).First(&users).Error; err != nil {
			return err
		} else {
			if len(users) == 0 {
				return errors.New("No user found")
			}
			*result = users
			return nil
		}
	}
}

/*Изменение пользователя

Пользователь до изменений:
{"result":[{"Uuid":"dac1dd0b-f5ca-4864-a59c-c3c5d7260aac","Login":"Hidejewel","RegDate":"2017-12-19"}],"error":null,"id":1}

Пример запроса:
{"jsonrpc":"2.0","method":"User.Update","id":1,"params":[{"uuid":"dac1dd0b-f5ca-4864-a59c-c3c5d7260aac","login":"Hidejewel-new"}]}

Ответ:
{"result":{"Uuid":"dac1dd0b-f5ca-4864-a59c-c3c5d7260aac","Login":"Hidejewel-new","RegDate":"2017-12-19"},"error":null,"id":1}
*/
func (t *User) Update(r *http.Request, args *User, result *User) error {
	var user User
	//Для изменения пользователя необходимо в запросе указать его uuid. Изменить можно логин и дату регистрации пользователя.
	if err:= db.Where("uuid = ?", args.Uuid).First(&user).Error; err != nil {
		return err
	} else {
		if (args.Login!="") {
			user.Login = args.Login
		}
		if (args.RegDate!="") {
			user.RegDate = args.RegDate
		}
		if err := db.Save(&user).Error; err != nil {
			return err
		} else {
			*result = user
			return nil
		}
	}
}

var db *gorm.DB

func main() {
	flag.StringVar(&host, "host", host, "Main server host name")
	flag.Parse()

	var err error
	db, err = gorm.Open("sqlite3", "users.db")
	if err!=nil {
		log.Fatal("Can't open DB")
	}
	defer db.Close()
	if db.HasTable(&User{}) != true {
		db.CreateTable(&User{})
	}

	s := rpc.NewServer()
	s.RegisterCodec(json.NewCodec(), "application/json")
	s.RegisterCodec(json.NewCodec(), "application/json;charset=UTF-8")
	user := new(User)
	s.RegisterService(user, "User")
	r := mux.NewRouter()
	r.Handle("/rpc", s)
	if err := http.ListenAndServe(host, r); err!=nil {
		log.Fatal("Can't start server. Error: ", err)
	}
}