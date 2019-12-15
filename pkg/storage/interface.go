package storage

import "strings"

type LogStorage interface {
	Create(id string)
	Insert(id string, out *strings.Builder, err *strings.Builder)
	Get(id string) (stdOut *strings.Reader, stdErr *strings.Reader)
	Close(id string)
}
