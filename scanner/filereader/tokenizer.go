package filereader

import "text/scanner"
import "io"
import "math/rand"
import log "github.com/cihub/seelog"
import "fmt"
import "bytes"
import "unicode"
import "regexp"

type TokenType int

var garbageRegexp = regexp.MustCompile(`^[-\$\.,0-9]+$`)

const (
    NullToken = 0
    TextToken TokenType = iota
    XMLStartToken
    XMLEndToken
    SymbolToken
)

type Token struct {
    Text string
    Type TokenType
    DocId DocumentId
    Position int
    Final bool
    // An unique identifier to denote phrases which do not cross
    // a punctuation barrier.
    // This allows indexers to identify phrases without having to
    // process punctuation itself
    PhraseId int
}

func (t *Token) Clone() *Token {
  newtok := NewToken(t.Text, t.Type)
  newtok.DocId = t.DocId
  newtok.Position = t.Position
  newtok.Final = t.Final
  newtok.PhraseId = t.PhraseId
  return newtok
}

func (t *Token) Equal(other *Token) (equal bool) {
  equal = true

  equal = t.Eql(other)

  if t.Position != other.Position {
  equal = false
}

  return
}


func (t *Token) Eql(other *Token) (equal bool) {
  equal = true

  if t.Type != other.Type {
    equal = false
  }

  if t.Type == NullToken {
    return true
  }

  if t.Text != other.Text {
    equal = false
  }


  return
}

func NewToken(text string, ttype TokenType) (*Token) {
  t := new(Token)
  t.Text = text
  t.Type = ttype
  t.Final = false
  t.PhraseId = 0
  return t
}


type Tokenizer interface {
    Next() (*Token, error)
    Tokens() <-chan *Token
    Reset()
}

func (t TokenType) String() string {
    switch t {
        case TextToken: return "TEXT"
        case XMLStartToken: return "XMLSTART"
        case XMLEndToken: return "XMLEND"
        case SymbolToken: return "SYMBOL"
        default: return "NULL"
    }
}

func (t Token) String() string {
  return fmt.Sprintf("%s [%s@%d:%d]", t.Text, t.Type, t.Position, t.PhraseId)
}

type BadXMLTokenizer struct{
    tok_start, tok_end int
    scanner  *scanner.Scanner
    rd io.ReadSeeker
    current_phrase_id int
}

func BadXMLTokenizer_FromReader(rd io.ReadSeeker) (Tokenizer){
    t := new(BadXMLTokenizer)
    t.rd = rd
    t.scanner = new(scanner.Scanner).Init(rd)
    t.scanner.Whitespace = 0
    t.scanner.Error = func(s *scanner.Scanner, msg string) { panic(msg)}
    //t.scanner.Mode = scanner.ScanStrings // Original scan mode
    t.scanner.Mode = 0
    t.current_phrase_id = rand.Intn(1000)
    return t
}

func (tz *BadXMLTokenizer) Reset() {
    _, err := tz.rd.Seek(0, 0)
    if err != nil {
        panic(err)
    }
    tz.scanner = new(scanner.Scanner).Init(tz.rd)
    tz.scanner.Whitespace = 0
    tz.scanner.Error = func(s *scanner.Scanner, msg string) { panic(msg)}
    //tz.scanner.Mode = scanner.ScanStrings // Original scan mode
    tz.scanner.Mode = 0
    tz.current_phrase_id = rand.Intn(1000)
}

var alnum = []*unicode.RangeTable{unicode.Digit, unicode.Letter,
unicode.Dash, unicode.Hyphen}

var symbols = []*unicode.RangeTable{unicode.Symbol,
unicode.Punct}

func (tz *BadXMLTokenizer) Next() (*Token, error) {

    for {
        tok := tz.scanner.Peek()
        log.Tracef("Scanner found: %v", tok)

        if tok == scanner.EOF {
            return nil, io.EOF
        }

        switch  {
        case unicode.IsPrint(tok) == false:
            log.Tracef("Skipping unprintable character")
            fallthrough
        case unicode.IsSpace(tok):
            tz.scanner.Scan()
            continue

        case tok == '<':
            log.Tracef("parsing XML")
            token, ok := parseXML(tz.scanner)
            // We actually bump the phrase no matter what. It's
            // either a comment, an xml token, or something weird
            tz.current_phrase_id = rand.Intn(1000)
            if ok {
                log.Tracef("Returning XML Token: %s", token)
                return token, nil
            }

        case tok == '&':
            log.Tracef("parsing HTML")
            if token := parseHTMLEntity(tz.scanner); token != nil {
                tz.current_phrase_id = rand.Intn(1000)
                return token, nil
            }

        case tok == '`': // Handle this speciallly - technically it's a 'grave accent'
            fallthrough
        case unicode.Is(unicode.Punct, tok):
            log.Tracef("Ignoring punctuation: %v", tok)
            tz.current_phrase_id = rand.Intn(1000)
            tok = tz.scanner.Scan()

        default:
            /* Catch special things in words */
            log.Tracef("Found '%s' . Parsing Text", string(tok))
            token, ok := tz.parseCompound()
            if ok {
              log.Debugf("Returing Text Token: %s", token)
              return token, nil
            } else {
              tz.scanner.Scan()
            }
        }
    }
}

func (t *BadXMLTokenizer) parseCompound() (*Token, bool) {
    var entity = new(bytes.Buffer)
    var compoundPhraseId = t.current_phrase_id

    for {
      next := t.scanner.Peek()
      log.Tracef("Next is '%v'. Text entity is %s",
        next, entity.String())

      switch {

      case next == '&':
        log.Tracef("parsing HTML")
        if token := parseHTMLEntity(t.scanner); token != nil {
          entity.WriteString(token.Text)
        } else {
          if entity.Len() > 0 {
            tok := NewToken(entity.String(), TextToken)
            tok.PhraseId = compoundPhraseId
            return tok, true
          }
        }

      case next == '<':
        if entity.Len() > 0 {
          tok := NewToken(entity.String(), TextToken)
          tok.PhraseId = compoundPhraseId
          return tok, true
        } else {
          return nil, false
        }

    case next == '/':
        if entity.Len() > 0 {
            log.Tracef("Found a '%c'.  Splitting on it. Returning %s.", next, entity.String())
            tok := NewToken(entity.String(), TextToken)
            tok.PhraseId = compoundPhraseId
            return tok, true
        } else {
            return nil, false
        }

      case unicode.IsOneOf(alnum, next):
        t.scanner.Scan()
        entity.WriteString(t.scanner.TokenText())

      case unicode.IsOneOf(symbols, next):
        t.scanner.Scan()
        part2, ok := t.parseCompound()

        log.Tracef("Parsing symbol %c. part2 is %s", next, part2)
        switch {

        case ok && next == '\'':
          if ok {
             entity.WriteString(part2.Text)
          }

        case ok && next == '(':
          entity.WriteRune('-')
          log.Tracef("subparse got %s. appending with '-'", part2.Text)
          entity.WriteString(part2.Text)

        case ok && next == ')':
          log.Tracef("Skipping ending paren. After the paren is %s. Entity is %s", part2.Text, entity.String())
          entity.WriteString(part2.Text)
          //Don't keep ending parens

        case ok:
          entity.WriteRune(next)
          entity.WriteString(part2.Text)

        case unicode.Is(unicode.Sc, next): //currency
          entity.WriteRune(next)

        case next == '\'':
          // Trailing single punctuation should not make a new phrase

        default:
          t.current_phrase_id = rand.Intn(1000)
        }

      default:
        if entity.Len() > 0 && garbageRegexp.FindString(entity.String()) == "" {
          tok := NewToken(entity.String(), TextToken)
          tok.PhraseId = compoundPhraseId
          return tok, true
        } else {
          return nil, false
        }
      }
    }
}

func (t *BadXMLTokenizer) parseParenthetical() (string, bool) {
  buf := new(bytes.Buffer)

  for {
    next := t.scanner.Next()
    switch {

    case next == ')':
      return buf.String(), true

    case unicode.IsSpace(next):
      return buf.String(), false

    default:
      buf.WriteRune(next)
    }
  }
}

func decodeEntity(entity string) (string, bool) {

  switch entity {
  case "&hyph;":
    return "-", true
  case "&blank;":
    return "", false
  case "&lt;":
    return "<", true
  case "&gt;":
    return ">", true

  case "&eacute;":
    return "\u00E9", true
  case "&uuml;":
    return "\u00fc", true
  case "&ntilde;":
    return "\u00f1", true
  case "&aacute;":
    return "\u00e1", true
  case "&iacute;":
    return "\u00ed", true
  case "&oacute;":
    return "\u00f3", true

  case "&rsquo;":
    fallthrough
  case "&lsquo;":
    return "'", true

  case "&mu;":
    fallthrough
  case "&para;":
    fallthrough
  case "&reg;":
    fallthrough
  case "&sect;":
    fallthrough
  case "&cir;":
    fallthrough
  case "&bull;":
    fallthrough
  case "&racute;":
    fallthrough
  case "&lacute;":
    fallthrough
  case "&tilde;":
    fallthrough
  case "&amp;":
    return "", false // no good, but don't warn

  default:
    log.Warnf("Invalid character escape sequence: %s", entity)
    return "", false
  }
}

// Return a token representing the HTML entity, or nil if 
// this decodes to something that should not be kept (and which
// breaks words
func parseHTMLEntity(sc *scanner.Scanner) (*Token) {

    var entity = new(bytes.Buffer)
    log.Tracef("ParseHTML. Starting with '%s'", entity.String())

    for {
        tok := sc.Scan()
        log.Tracef("Parse HTML. Reading token %c", tok)

        switch {
        case unicode.IsSpace(tok):
            if entity.Len() > 1 {
                token := NewToken(entity.String(), SymbolToken)
                log.Tracef("ParseHTML. Returning non-HTML '%s'",
            token.Text)
                return token
            } else {
                return nil
            }

        case tok == ';':
            entity.WriteRune(tok)
            log.Tracef("Attempting to decode %s", entity.String())
            if decoded, ok := decodeEntity(entity.String()); ok {
                token := NewToken(decoded, SymbolToken)
                log.Tracef("ParseHTML. Returning HTML '%s'", entity.String())
                return token
            } else {
              return nil
            }

        default:
            entity.WriteString(sc.TokenText())
        }
    }

}

func parseXML(sc *scanner.Scanner) (*Token, bool) {

    var entity = new(bytes.Buffer)
    token := new(Token)

    // Skip the '<'
    sc.Scan()

    switch sc.Peek() {
    case '/':
        token.Type = XMLEndToken
        sc.Next()
    case '!':
        log.Tracef("parseXML skipping comment")
        next := sc.Next()
        for next != '>' {
            next = sc.Next()
        }
        return nil, false
    default:
        token.Type = XMLStartToken
    }

    log.Tracef("parseXML creating %s element", token.Type )

    for {
        tok := sc.Scan()
        log.Tracef("parseXML found %s. Token is %v. Entity is: '%s'",
          sc.TokenText(),
          tok,
          entity.String())

        switch {
        case tok == '>':
            token.Text = entity.String()
            return token, true

        case unicode.IsSpace(tok):
            return nil, false

        default:
            log.Tracef("parseXML appending %s to string",
            sc.TokenText())
            entity.WriteString(sc.TokenText())

        }
    }
}

func (tz *BadXMLTokenizer) Tokens() (<- chan *Token) {

    token_channel := make(chan *Token)
    log.Debugf("Created channel %v as part of Tokens(), with" +
              " Scanner = %v", token_channel, tz)

    go func(ret chan *Token, tz *BadXMLTokenizer) {
        for {
            log.Tracef("Scanner calling Next()")
            tok, err := tz.Next()
            log.Tracef("scanner.Next() returned %s, %v", tok, err)
            switch err {
            case nil:
                log.Debugf("Pushing %s into token channel %v",
                tok, ret)
                ret <- tok
            case io.EOF:
                log.Debugf("received EOF, closing channel")
                close(ret)
                log.Debugf("Closed.")
                log.Flush()
                return
                panic("I should have exited the goroutine but " +
                "didn't")
            }
        }
    }(token_channel, tz)

    return token_channel
}
