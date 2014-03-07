package suggester

import(
  "sync"
  "github.com/wwwjscom/ocr_engine/db"
  "fmt"
)

type EditDistanceSuggester struct {
  mysql_chan chan *db.Mysql
  tables_to_search []string
  suggestions_chan chan *EditDistanceSuggestion
  suggestions []*EditDistanceSuggestion
}

func NewEditDistanceSuggester(mysql_chan chan *db.Mysql, tables_to_search []string) *EditDistanceSuggester {  
  ed := new(EditDistanceSuggester)
  ed.suggestions = make([]*EditDistanceSuggestion, 0)
  ed.mysql_chan = mysql_chan
  ed.tables_to_search = tables_to_search
  return ed
}

// Implements the Suggester interface
func (ed *EditDistanceSuggester) Suggest(word string) []*EditDistanceSuggestion {
  ed.suggestions_chan = make(chan *EditDistanceSuggestion)
  
  wg := new(sync.WaitGroup)

  // Add all the suggestions from the go routiunes to our suggestions array
  go ed.addSuggestions()  

  for _, table_name := range ed.tables_to_search {
    wg.Add(1)
    go ed.search(word, table_name, wg)
  }

  // Wait for the wait group to empty
  wg.Wait()
  
  close(ed.suggestions_chan)
    
  return ed.suggestions
}

func (ed *EditDistanceSuggester) addSuggestions() {
  for s := range ed.suggestions_chan {
    ed.suggestions = append(ed.suggestions, s)
  }
}

// Executes the edit distance searching logic
func (ed *EditDistanceSuggester) search(word string, table_name string, wg *sync.WaitGroup) {
  min_length, max_length := determineSubspaceLengths(word)
  
  q1 := "DROP TEMPORARY TABLE IF EXISTS table2;"
  q2 := fmt.Sprintf("CREATE TEMPORARY TABLE IF NOT EXISTS table2 AS(SELECT * FROM %s WHERE name RLIKE '^.{%d,%d}$');", table_name, min_length, max_length)
  q3 := fmt.Sprintf("SELECT name, levenshtein('%s', `name`) FROM table2 WHERE levenshtein('%s', `name`) BETWEEN 0 AND 1;", word, word)
  
  mysql_conn := <- ed.mysql_chan
  mysql_conn.Query(q1)
  mysql_conn.Query(q2)
  results := mysql_conn.Query(q3)
  ed.mysql_chan <- mysql_conn
  
  for _, row := range results {
    newSuggestion := new(EditDistanceSuggestion)
    newSuggestion.OriginalText = word
    newSuggestion.Text = row.Str(0)
    newSuggestion.Confidence = row.Int(1)
    ed.suggestions_chan <- newSuggestion
  }

  wg.Done()
}

// Determines the size of all words in the subspace that edit distance will evaluate against.
func determineSubspaceLengths(word string) (int, int){
  min_length := len(word) - 2
  if min_length < 0 {
    min_length = 0
  }
  
  max_length := len(word) + 2
  return min_length, max_length
}

//func (ed *EditDistanceSuggester) mysqlToSuggestions(mysqlResults []Interface) []*EditDistanceSuggestion {
//  suggestions := make([]*EditDistanceSuggestion)
//  
//  for _, result := range results {
//    newSuggestion := new(EditDistanceSuggestion)
//    newSuggestion.suggestion = result.value
//    suggestions = append(suggestions, newSuggestion)
//  }
//  
//  return suggestions
//}



// Implements the suggestion interface
type EditDistanceSuggestion struct {
  OriginalText string
  Confidence int 
  Text string
}