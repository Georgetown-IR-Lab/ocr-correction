package filewriter

import "testing"
import "os"
import "bufio"
import "fmt"
import log "github.com/cihub/seelog"

func TestTrecFileWriter(t *testing.T) {
    filename := "/tmp/test_file_123"
    fw := new(TrecFileWriter)
    fw.Init(filename)
    go fw.WriteAllTokens()

    // Write these words to disk via the writer channel
    words := [3]string{"word1", "word2", "word3"}
    for i, _ := range words {
        log.Debugf("Adding %s to writer chan", words[i])
        fw.StringChan <- &words[i]
    }

    close(fw.StringChan)
    log.Debugf("Writer channel closed")

    // Verify file contents
    if file, err := os.Open(filename); err != nil {
        panic(fmt.Sprintf("Unable to open %s due to error: %s\n", filename, err))
    } else {
        scanner := bufio.NewScanner(file)
        ctr := 0
        for scanner.Scan() {
            if (words[ctr] != scanner.Text()) {
                t.Errorf("%s found, should have been %s", scanner.Text(), words[ctr])
            }
            ctr++
        }

        file.Close()
    }

}
