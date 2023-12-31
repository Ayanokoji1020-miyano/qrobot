package qrobot

import (
	"context"
	"github.com/Ayanokoji1020-miyano/qrobot/service"
	"log"
	"os"
	"os/signal"
)

func Login(ctx context.Context, qq int64, protocol int) {
	err := login(qq, "", protocol)
	if err != nil {
		log.Println(err)
		return
	}

	go service.Run(ctx, qqClient)

	// 等待退出信号
	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt, os.Kill)
	<-ch
}
