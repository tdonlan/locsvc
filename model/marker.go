package model

import (
	_ "database/sql"
	"errors"
	_ "github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"time"
)

type Marker struct {
	Id          int
	Text        string
	UserId      int
	Lat         float64
	Lon         float64
	TimeCreated time.Time `db:"timeCreated"`
}

type CreateMarker struct {
	SessionId string
	Text      string
	Lat       float64
	Lon       float64
}

type SearchMarkers struct {
	SessionId string
	Lat       float64
	Lon       float64
}

func (dal *ModelDAL) CreateMarker(marker *CreateMarker, userId int) (*Marker, error) {

	sql := `insert into markers (text,userid,lat,lon,timecreated) values (?,?,?,?,?)`
	result, err := dal.db.Exec(sql, marker.Text, userId, marker.Lat, marker.Lon, time.Now())
	if err != nil {
		log.Print(err)
		return nil, err
	}

	if rows, _ := result.RowsAffected(); rows == 0 {
		return nil, errors.New("Failed to insert row")
	}

	id, _ := result.LastInsertId()
	newSession, err := dal.GetMarkerById(int(id))

	return newSession, err
}

//TODO: Delete Session
func (dal *ModelDAL) GetMarkerById(id int) (*Marker, error) {
	var retval Marker

	sql := `select id,text,userid,lat,lon, timecreated from markers where id = ?`
	err := dal.db.Get(&retval, sql, id)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	return &retval, nil
}

/*
Latitude: 1 deg = 110.574 km
Longitude: 1 deg = 111.320*cos(latitude) km
*/
//Defaulting to searching within 1km
func (dal *ModelDAL) SearchMarkersByLoc(lat, lon float64) ([]*Marker, error) {
	var retval []*Marker

	sql := `select id,text,userid,lat,lon, timecreated from markers where abs(lat-?) < .001 AND abs(lon-?) < .001`
	err := dal.db.Select(&retval, sql, lat, lon)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	return retval, nil
}

//TODO: GetMarkersByUser //return list of all markers the user has created

//TODO: Delete markers

/*
func (dal *ModelDAL) GetSessionBySessionId(sessionid string) (*Session, error) {
	var retval Session

	sql := `select id,text,userid,lat,lon, timecreated from markers where id = ?`
	err := dal.db.Get(&retval, sql, sessionid)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	return &retval, nil
}
*/
