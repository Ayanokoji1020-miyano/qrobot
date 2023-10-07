package qqRobot

import (
	"context"
	"log"
	"os"
	"os/signal"
	"qqRobot/service"
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
