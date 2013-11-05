package tagger

import (
    //"os"
    log "github.com/cihub/seelog"
    "github.com/ziutek/mymysql/mysql"
    _ "github.com/ziutek/mymysql/native" // Native engine
)

type DB_conn struct {
    user string
    pass string
    name string

    conn mysql.Conn
}

func New(user, pass, name string) *DB_conn {
    db := DB_conn {
        user: user,
        pass: pass,
        name: name,
    }

    db.conn = mysql.New("tcp", "", "127.0.0.1:3306", db.user, db.pass, db.name)

    err := db.conn.Connect()
    if err != nil {
        log.Criticalf("Cannot connect to %s using username: %s and pass: %s", db.name, db.user, db.pass)
        panic(err)
    }

    return &db
}
