package main

import (
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/teachain/stats/internal/config"
	"github.com/teachain/stats/internal/logic"
	"os"
	"os/signal"
	"syscall"
)

var configFile *string = flag.String("config", "./etc/config.yaml", "Path to config file")

func main() {
	flag.Parse()
	c, err := config.MustLoadConfig(*configFile)
	if err != nil {
		fmt.Printf("Load file %s:%s\n", *configFile, err.Error())
		return
	}
	builder, err := logic.NewBuilder(c)
	if err != nil {
		fmt.Printf("NewBuilder %s\n", err.Error())
		return
	}
	err = builder.Start()
	if err != nil {
		fmt.Printf("Builder Start%s\n", err.Error())
		return
	}
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	builder.Stop()
}
