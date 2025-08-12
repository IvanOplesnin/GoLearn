package main

import (
	"flag"
	"log"
	"os"
	"time"
)

type Config struct {
	Dir       string
	Debounce  time.Duration
	Recursive bool
	Mode      string
}

const (
	ModeFSnotify string = "fsnotify"
	ModePolling  string = "polling"
)

func ParseConfig() Config {
	cfg := Config{}

	flag.StringVar(&cfg.Dir, "dir", ".", "Папка, которую отслеживают")
	debounce := flag.Int("deb", 300, "Интервал разгрузки сообщений(миллисикунды)")
	flag.BoolVar(&cfg.Recursive, "recursive", false, "Смотреть за внутреннними директориями рекурсивно")
	flag.StringVar(&cfg.Mode, "mode", ModeFSnotify, "Режим работы")

	flag.Parse()

	cfg.Debounce = time.Duration(*debounce) * time.Millisecond

	if cfg.Mode != ModeFSnotify && cfg.Mode != ModePolling {
		flag.Usage()
		os.Exit(2)
	}
	if st, err := os.Stat(cfg.Dir); err != nil || !st.IsDir() {
		log.Printf("Не правильно прописан путь к папке, %s", cfg.Dir)
		flag.Usage()
		os.Exit(2)
	}

	return cfg
}
