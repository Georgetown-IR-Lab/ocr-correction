package filereader

import "testing"
import "strings"
import log "github.com/cihub/seelog"
import "github.com/wwwjscom/ocr_engine/logging"

func TestTrecFileReader(t *testing.T) {
  logging.SetupTestLogging()

  log.Debugf("Creating FileReader")
  fr := new(TrecFileReader)

  log.Debugf("Opening file")
  fr.Init("test/testfile1.txt")

  log.Debugf("Reading file")
  doc := fr.Read()
  if id := doc.OrigIdent(); id != "12345" {
    t.Error("Failed to parse document id")
  }

  tokens := doc.Tokens()
  exp_tokens :=  expected()
  i := 0
  for tok := range tokens {
    exp, ok := <- exp_tokens
    if ! ok {
      if tok.Type != NullToken {
        t.Error("Read all of 'expected' before reaching end of tokens")
      } else {
      
      }
    }
    i += 1

    if pos := tok.Position; pos != i {
      t.Errorf("'%s' was not at position %d as expected", tok, i)
    }

    if id := tok.DocId; id != doc.Identifier() {
      t.Errorf("Token '%s' did not have DocId matching '%s'",tok, id)
    }

    if tok.Text != exp {
      t.Errorf("%s did not match %s in position %d", tok, exp, i )
      break
    } else {
    }
  }
}

func expected() <-chan string {

  expected := `
DEPARTMENT OF AGRICULTURE
Agricultural Marketing Service
7 CFR Part 28
CN-94-003
RIN 0581-AB06
Cotton Classification Services for Cotton Producers Withdrawal of Proposed Rule 
AGENCY
Agricultural Marketing Service USDA
ACTION
 Proposed rule withdrawal 
SUMMARY
 This document withdraws a proposed rule that would have amended regulations governing cotton classification services
provided to cotton producers by establishing a module averaging method of cotton classification That rule would
have changed the present classification system by adding the new procedure
DATES
 This proposed rule is withdrawn effective April 5 1994
FOR FURTHER INFORMATION CONTACT
 Craig Shackelford 202-720-2259
SUPPLEMENTARY INFORMATION
 The proposed rule was issued as amendments to regulations governing Cotton Classification Services for Producers
7 CFR part 28 The proposal was issued on January 31 1994 and published in the
Federal Register
 59 FR 4257 It proposed the implementation of module averaging a method by which the accuracy of fiber quality measurements
can be improved The module averaging procedure would use all the bales from a module or trailer as the testing unit
rather than using a single bale as the test unit The module averaging procedure has been offered to growers on a voluntary
basis for the past three crop years 
The Secretary of Agricultures Advisory Committee on Cotton Marketing recommended that if no significant problems
were encountered during the 1993 classing season the module averaging procedure be expanded to include all cotton
classed in 1994 and subsequent crop years For the 1993 expanded voluntary program there were 242 gins participating
and the production from these gins totaled 3,053,716 bales This represented 20 percent of the 1993 cotton crop No
problems of any significance are known to have developed during the 1993 project In keeping with the advisory committees
recommendation AMS proposed that module averaging be applied to all bales classed effective with the 1994 cotton
crop 
Written comments regarding this proposal were accepted from January 31 1994 through March 2 1994 Comments were
received from individuals and organizations representing several segments of the cotton industry including producers
ginners warehousers merchants cooperatives national and international trade associations textile manufacturers
and others This broad cross section of the cotton industry together with a significant number of comments 61 indicates
a strong interest in the module averaging concept throughout the industry 
The textile manufacturing segment submitted four comments one from a national organization representing domestic
textile manufacturers and three from individual firms All of these comments expressed support for the proposal
The national organization representing textile manufacturers favored the implementation of the proposal provided
four conditions were met 1 That at least two tests per bale be made 2 that the integrity of the module be maintained
by preventing the intermingling of cotton between module test units 3 that when the module averaging procedure
is used only those measurements that fall within three standard deviations of the average be included in the module
average and 4 that for the purposes of review classification all bales from the module be retested The remaining
three comments all reiterated the suggestion for the use of three standard deviations as a determination for including
bales in the module average The Agency is currently and will continue studying how best to determine the inclusion
or exclusion of bales from the module average 
Thirteen comments were received from the producer segment including two from national organizations One national
organization was in favor of the proposal The other recommended the delay of mandatory implementation until 1995
and the application of module averaging to length strength and micronaire measurements only The remaining comments
submitted by regional producer organizations and individual producers were nearly equally divided among those
favoring the proposal and those supporting the module averaging concept but suggesting the continuation of the
voluntary program 
The cotton ginning segment through national and state organizations and individual ginners submitted 14 written
comments The national organization representing cotton ginners encouraged the comprehensive industry review
of the 1993 and 1994 module averaging results and continuation of the voluntary module averaging program The remaining
comments expressed general support for the module averaging concept but suggested that module averaging be continued
on a voluntary basis 
A leading national association representing cotton merchants shippers and exporters of raw cotton opposed the
mandatory implementation of module averaging while supporting the continuation of the voluntary program The Agency
is currently responding to a separate request from this organization for statistical data pertaining to module averaged
cotton This organization maintains that further evaluation of the data is necessary prior to the implementation
of mandatory module averaging The remaining 22 comments received from the merchant segment also requested that
module averaging be continued on a voluntary basis so that the effects of module averaging on cotton marketing can
be further evaluated 
Three foreign cotton trade organizations submitted comments These comments all stated that the international
cotton trade was not yet sufficiently knowledgeable about the module averaging concept to favor its use on anything
more than a voluntary basis 
In light of the views expressed in comments submitted from the various segments of the cotton industry the Agency
has determined that it is in the public interest to continue with the module averaging program on a voluntary basis
Such action will provide the cotton industry more time to evaluate the effects of module averaging on cotton marketing
Accordingly the proposed rule is withdrawn effective April 5 1994
Dated March 30 1994
Lon Hatamiya
Administrator
FR Doc 94-8028 Filed 4-4-94 8:45 am
BILLING CODE 3410-02-P
`

  c := make(chan string)

  // Load the channel
  go func (c chan string, arr []string ){
    for _, s := range arr {
      c <- s
    }
    close(c)
  }(c, strings.Fields(expected))

  return c
}
