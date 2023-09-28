package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"qqRobot/service"
)

func main() {
	err := login(2918752369, "", ProtocolWatch)
	if err != nil {
		log.Println(err)
		return
	}

	//todo ctx 改为子上下文
	ctx := context.Background()
	go service.Run(ctx, qqClient)

	// 等待退出信号
	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt, os.Kill)
	<-ch
}
