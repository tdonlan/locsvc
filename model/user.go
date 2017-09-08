package model

import (
	_ "database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
)

type ModelDAL struct {
	db *sqlx.DB
}

func OpenDB(dataSourceName string) *ModelDAL {
	if _, err := os.Stat(dataSourceName); os.IsNotExist(err) {
		log.Print(dataSourceName + " does not exist.\n")
	}

	db, err := sqlx.Open("sqlite3", dataSourceName)
	log.Print("Opening " + dataSourceName)
	if err != nil {
		log.Print(err)
	}

	if err = db.Ping(); err != nil {
		log.Print(err)
	}

	return &ModelDAL{db}
}

type User struct {
	Id       int
	Name     string
	Password string
}

func (dal *ModelDAL) CreateUser(data *User) (*User, error) {
	//TODO: verify user doesn't exist with matching name

	sql := `insert into users (name, password) values (?,?)`
	result, err := dal.db.Exec(sql, data.Name, data.Password)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	if rows, _ := result.RowsAffected(); rows == 0 {
		return nil, errors.New("Failed to insert row")
	}

	id, _ := result.LastInsertId()
	newUser, err := dal.GetUserById(int(id))

	return newUser, err
}

func (dal *ModelDAL) GetUserById(id int) (*User, error) {
	var retval User

	sql := `select id,name,password from users where id = ?`
	err := dal.db.Get(&retval, sql, id)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	return &retval, nil
}

func (dal *ModelDAL) GetUserByName(name string) (*User, error) {
	var retval User

	sql := `select id,name,password from users where name = ?`
	err := dal.db.Get(&retval, sql, name)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	return &retval, nil
}
