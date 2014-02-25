package suggester

import(
  "fmt"
  "testing"
  "github.com/wwwjscom/ocr_engine/db"
)

type testcase struct {
  query string
  expected_min int
  expected_max int
  expected_result_size int
}

var (
  tests []testcase
  ed *EditDistanceSuggester
)

func init() {
  // setup
  table_names := []string{"names", "dict", "geo"}
  mysql_conn := db.NewMySQLConn("root", "", "ocr_research")
  mysql_chan := make(chan *db.Mysql, 1)
  go func() { mysql_chan <- mysql_conn }()

  // create
  ed = NewEditDistanceSuggester(mysql_chan, table_names)
  
  tests = []testcase {
    {
      "Cat",
      1,
      5,
      68,
    },
    {
      "Bo",
      0,
      4,
      144,
    },
    {
      "A",
      0,
      3,
      223,
    },
  }  
}


func TestDetermineSubspaceLengths(t *testing.T) {
  for _, tc := range tests {
    min, max := determineSubspaceLengths(tc.query)
    if min != tc.expected_min {
      t.Fail()
    }
    if max != tc.expected_max {
      t.Fail()      
    }
    
  }
}

func TestSuggester(t *testing.T) {  
  for _, tc := range tests {
    suggestions := ed.Suggest(tc.query)    
    if suggestions == nil {
      t.Fail()
    }  
    
    if len(suggestions) != tc.expected_result_size {
      t.Errorf(fmt.Sprintf("Too few results returned.  Expected %d, but received %d", tc.expected_result_size, len(suggestions)))
    }
  }  
}