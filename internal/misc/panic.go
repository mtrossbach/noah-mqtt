package misc

import (
	"log/slog"
	"time"
)

func Panic(err error) {
	slog.Error("PANIC ... waiting 15s to exit ...")
	<-time.After(15 * time.Second)
	panic(err)
}
