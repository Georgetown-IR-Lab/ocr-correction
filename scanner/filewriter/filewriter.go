package filewriter


type FileWriter interface {
    Init(string)
    Write()
    WriteAllTokens()
}
