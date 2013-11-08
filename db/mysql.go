package db

import (
    //"os"
    log "github.com/cihub/seelog"
    "github.com/ziutek/mymysql/mysql"
    _ "github.com/ziutek/mymysql/native" // Native engine
)

// Struct for managing connections to the database
type Mysql struct {
    user string
    pass string
    name string

    conn mysql.Conn
}

func (db *Mysql) New(user, pass, name string) *Mysql {
    db = new(Mysql)

    db.user = user
    db.pass = pass
    db.name = name
    //db = Mysql{
    //    user: user,
    //    pass: pass,
    //    name: name,
    //}

    db.conn = mysql.New("tcp", "", "127.0.0.1:3306", db.user, db.pass, db.name)

    err := db.conn.Connect()
    if err != nil {
        log.Criticalf("Cannot connect to %s using username: %s and pass: %s", db.name, db.user, db.pass)
        panic(err)
    }

    return db
}

func (db *Mysql) Query(query string) ([]mysql.Row) {
    // "select * from X where id > %d", 20
    log.Tracef("Searching the db for %s", query)
    rows, _, err := db.conn.Query(query)
    if err != nil {
        panic(err)
    }
    return rows
}
