package db

import (
    "github.com/ziutek/mymysql/mysql"
    _ "github.com/ziutek/mymysql/native" // Native engine
)

type Lexicon struct {
    mysql *Mysql
    table_name string
    field_name string
}

func (l *Lexicon) Init(mysql *Mysql, table_name string) {
    l.mysql = mysql
    l.table_name = table_name
}

func (l *Lexicon) Query(word *string) []mysql.Row {
    q := "SELECT * FROM " + l.table_name + " WHERE " + l.field_name + " = \"" + *word + "\""
    return l.mysql.Query(q)
}

func (l *Lexicon) Find(word *string) bool {
  if *word == "foo" {
    return false
  } else {
    return true
  }
}