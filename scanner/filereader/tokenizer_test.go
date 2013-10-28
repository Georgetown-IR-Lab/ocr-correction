package filereader

import "testing"
import "strings"
import "math/rand"
import "fmt"
import log "github.com/cihub/seelog"
import "github.com/wwwjscom/ocr_engine/logging"

type testcase struct {
	test     string
	expected []Token
}

var (
	RandInts []int
	tests    []testcase
)

func init() {
	rand.Seed(0)
	RandInts = make([]int, 20)
	for x := 0; x < 20; x++ {
		RandInts[x] = rand.Intn(1000)
    log.Infof("RandInt[%d] = %d", x, RandInts[x])
	}

	tests = []testcase{
		{
			"(7 CFR) 8c(15)(A)",
			[]Token{
				Token{Text: "7", Type: TextToken, PhraseId: RandInts[1]},
				Token{Text: "CFR", Type: TextToken, PhraseId: RandInts[1]},
				Token{Text: "8c-15-A", Type: TextToken, PhraseId: RandInts[2]},
				Token{Type: NullToken},
			},
		},
		{
			"welcome; this is jim's house. Not jims' house ``Act.'' after &hyph; me ",
			[]Token{
				Token{Text: "welcome", Type: TextToken, PhraseId: RandInts[0]},
				Token{Text: "this", Type: TextToken, PhraseId: RandInts[1]},
				Token{Text: "is", Type: TextToken, PhraseId: RandInts[1]},
				Token{Text: "jims", Type: TextToken, PhraseId: RandInts[1]},
				Token{Text: "house", Type: TextToken, PhraseId: RandInts[1]},
				Token{Text: "Not", Type: TextToken, PhraseId: RandInts[2]},
				Token{Text: "jims", Type: TextToken, PhraseId: RandInts[2]},
				Token{Text: "house", Type: TextToken, PhraseId: RandInts[2]},
				Token{Text: "Act", Type: TextToken, PhraseId: RandInts[4]},
				Token{Text: "after", Type: TextToken, PhraseId: RandInts[5]},
				Token{Text: "-", Type: SymbolToken, PhraseId: 0},
				Token{Text: "me", Type: TextToken, PhraseId: RandInts[6]},
				Token{Type: NullToken},
			},
		},
		{
			"Mayag&uuml;ez A&ntilde;asco Gu&aacute;nica D&iacute;az Rinc&oacute;n",
			[]Token{
				Token{Text: "Mayag\u00fcez", Type: TextToken, PhraseId: RandInts[0]},
				Token{Text: "A\u00f1asco", Type: TextToken, PhraseId: RandInts[0]},
				Token{Text: "Gu\u00e1nica", Type: TextToken, PhraseId: RandInts[0]},
				Token{Text: "D\u00edaz", Type: TextToken, PhraseId: RandInts[0]},
				Token{Text: "Rinc\u00f3n", Type: TextToken, PhraseId: RandInts[0]},
				Token{Type: NullToken},
			},
		},
		{
			"8:43pm 100.242 100,000,1.10",
			[]Token{
				Token{Text: "8:43pm", Type: TextToken, PhraseId: RandInts[0]},
				Token{Text: "100.242", Type: TextToken, PhraseId: RandInts[0]},
				Token{Text: "100,000,1.10", Type: TextToken, PhraseId: RandInts[0]},
				Token{Type: NullToken},
			},
		},
		{
			"<PARENT> FR940405-1-00001 </PARENT>",
			[]Token{
				Token{Text: "PARENT", Type: XMLStartToken, PhraseId: 0},
				Token{Text: "FR940405-1-00001", Type: TextToken, PhraseId: RandInts[1]},
				Token{Text: "PARENT", Type: XMLEndToken, PhraseId: 0},
				Token{Type: NullToken},
			},
		},
		{
			"<DOCNO>DEADBEEF</DOCNO>",
			[]Token{
				Token{Text: "DOCNO", Type: XMLStartToken, PhraseId: 0},
				Token{Text: "DEADBEEF", Type: TextToken, PhraseId: RandInts[1]},
				Token{Text: "DOCNO", Type: XMLEndToken, PhraseId: 0},
				Token{Type: NullToken},
			},
		},
		{
			"<CFRNO>7 $CFR Part£ <!-- blah elsld --> 28 </CFRNO>",
			[]Token{
				Token{Text: "CFRNO", Type: XMLStartToken, PhraseId: 0},
				Token{Text: "7", Type: TextToken, PhraseId: RandInts[1]},
				Token{Text: "$CFR", Type: TextToken, PhraseId: RandInts[1]},
				Token{Text: "Part£", Type: TextToken, PhraseId: RandInts[1]},
				Token{Text: "28", Type: TextToken, PhraseId: RandInts[2]},
				Token{Text: "CFRNO", Type: XMLEndToken, PhraseId: 0},
				Token{Type: NullToken},
			},
		},
		{
			"<RINDOCK>[CN&hyph;94&hyph;003] </RINDOCK>",
			[]Token{
				Token{Text: "RINDOCK", Type: XMLStartToken, PhraseId: 0},
				Token{Text: "CN-94-003", Type: TextToken, PhraseId: RandInts[2]},
				Token{Text: "RINDOCK", Type: XMLEndToken, PhraseId: 0},
				Token{Type: NullToken},
			},
		},
		{
			"&blank;/&blank;Vol. 59, No. 2&blank;/&blank;Tuesday, January 4, 1994&blank;/&blank;Rules and Regulations",
			[]Token{
				Token{Text: "Vol", Type: TextToken, PhraseId: RandInts[1]},
				Token{Text: "59", Type: TextToken, PhraseId: RandInts[2]},
				Token{Text: "No", Type: TextToken, PhraseId: RandInts[3]},
				Token{Text: "2", Type: TextToken, PhraseId: RandInts[4]},
				Token{Text: "Tuesday", Type: TextToken, PhraseId: RandInts[5]},
				Token{Text: "January", Type: TextToken, PhraseId: RandInts[6]},
				Token{Text: "4", Type: TextToken, PhraseId: RandInts[6]},
				Token{Text: "1994", Type: TextToken, PhraseId: RandInts[7]},
				Token{Text: "Rules", Type: TextToken, PhraseId: RandInts[8]},
				Token{Text: "and", Type: TextToken, PhraseId: RandInts[8]},
				Token{Text: "Regulations", Type: TextToken, PhraseId: RandInts[8]},
				Token{Type: NullToken},
			},
		},
	}
}

func TestTokenizer(t *testing.T) {
	logging.SetupTestLogging()

	for _, test := range tests {
		run_testcase(test, t)
	}
}

func run_testcase(test testcase, t *testing.T) {

	reader := strings.NewReader(test.test)
	rand.Seed(0) //Seed so it matches expectations
	tz := BadXMLTokenizer_FromReader(reader)

	i := 0
	for tok := range tz.Tokens() {
		expected := test.expected[i]

		log.Debugf("TEST %d: '%s' == '%s'\n", i, tok, expected)

		if ! tok.Eql(&expected) {
			t.Error(fmt.Sprintf("%s != %s at %d\n", tok, expected, i))
		}
		if tok.PhraseId != expected.PhraseId {
			t.Errorf("PhraseId mismatch for %s. Expected %d. Got %d. ", tok, expected.PhraseId, tok.PhraseId)
		}

		i += 1
		if i > len(test.expected) {
			t.Error(fmt.Sprintf("Tokenizer has more tokens (%d) than expected (%d)",
				i, len(test.expected)))
			break
		}
	}

}
