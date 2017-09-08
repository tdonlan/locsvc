package model

import (
	_ "database/sql"
	"errors"
	_ "github.com/mattn/go-sqlite3"
	"github.com/satori/go.uuid"
	"log"

	"time"
)

type Session struct {
	Id          int
	SessionId   string
	UserId      int
	TimeCreated time.Time `db:"timeCreated"`
}

func (dal *ModelDAL) CreateSession(userid int) (*Session, error) {
	//todo: verify user doesn't have a session.
	newSessionId := uuid.NewV4()

	sql := `insert into sessions (sessionid,userid,timecreated) values (?,?,?)`
	result, err := dal.db.Exec(sql, newSessionId, userid, time.Now())
	if err != nil {
		log.Print(err)
		return nil, err
	}

	if rows, _ := result.RowsAffected(); rows == 0 {
		return nil, errors.New("Failed to insert row")
	}

	id, _ := result.LastInsertId()
	newSession, err := dal.GetSessionById(int(id))

	return newSession, err
}

//TODO: Delete Session

func (dal *ModelDAL) GetSessionById(id int) (*Session, error) {
	var retval Session

	sql := `select id, sessionid, userid, timecreated from sessions where id = ?`
	err := dal.db.Get(&retval, sql, id)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	return &retval, nil
}

func (dal *ModelDAL) GetSessionBySessionId(sessionid string) (*Session, error) {
	var retval Session

	sql := `select id,sessionid, userid, timecreated from sessions where sessionid = ?`
	err := dal.db.Get(&retval, sql, sessionid)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	return &retval, nil
}
