package storage

import (
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func init() {
	_ = os.Mkdir("logs", 0755)
}

type FileLogStorage struct {
}

func (store *FileLogStorage) Create(id string) {
	errLog, err := os.Create("logs/" + id + "-err.log")
	if err != nil {
		log.Print(err)
	}
	defer errLog.Close()
	outLog, err := os.Create("logs/" + id + "-out.log")
	if err != nil {
		log.Print(err)
	}
	outLog.Close()
}

func (store *FileLogStorage) Insert(id string, stdOut *strings.Builder, stdErr *strings.Builder) {
	outFile, err := os.OpenFile("logs/"+id+"-out.log", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0755)
	if err != nil {
		log.Print(err)
	}
	defer outFile.Close()
	_, err = outFile.WriteString(stdOut.String())
	if err != nil {
		log.Print(err)
	}
	outFile.Sync()

	errFile, err := os.OpenFile("logs/"+id+"-err.log", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0755)
	if err != nil {
		log.Print(err)
	}
	defer errFile.Close()
	_, err = errFile.WriteString(stdErr.String())
	if err != nil {
		log.Print(err)
	}
	errFile.Sync()
}

func (store *FileLogStorage) Get(id string) (*strings.Reader, *strings.Reader) {
	outBytes, err := ioutil.ReadFile("logs/" + id + "-out.log")
	if err != nil {
		log.Print(err)
	}
	errBytes, err := ioutil.ReadFile("logs/" + id + "-err.log")
	if err != nil {
		log.Print(err)
	}
	stdOut := strings.NewReader(string(outBytes))
	errOut := strings.NewReader(string(errBytes))

	return stdOut, errOut
}

func (store *FileLogStorage) Close(id string) {

}
