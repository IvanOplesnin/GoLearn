package fsnotify

import (
	"context"
	"time"
)

type Op string

const (
	CreateOp Op = "Create"
	WriteOp  Op = "Write"
	DeletOp  Op = "Delete"
	Rename   Op = "Rename"
)

type Config struct {
	Root          string        // корневая папка для наблюдения
	Recursive     bool          // следить за подпапками
	BufferBytes   int           // размер буфера на запрос
	QueueCapacity int           // глубина очереди сырых событий
	Debounce      time.Duration // окно коалесценции (делается вне источника — при желании)
}

type Event struct {
	Path    string
	Op      Op
	IsDir   bool
	OldPath string
	When    time.Time
}

type Watcher interface {
	Events() <-chan Event            // нормализованные события
	Errors() <-chan error            // ошибки
	Start(ctx context.Context) error // запускает горутины и подписки
	Close() error                    // явное закрытие (равносильно отмене контекста)
}
