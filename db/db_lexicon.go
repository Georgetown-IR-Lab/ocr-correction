package db

type Lexicon struct {
    mysql Mysql
    table_name string
    field_name string
}

type (*l Lexicon) Init(mysql Mysql, table_name string) {
    l.mysql = mysql
    l.table_name = table_name
}

type (*l Lexicon) Query(word *string) ([]mysql.Row) {
    q := "SELECT * FROM " + l.table_name + " WHERE " + l.field_name + " = \"" + word + "\""
    mysql.Query(q)
}
