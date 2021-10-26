package edu_solution_api

import (
	"context"
	"github.com/ozonmp/edu-solution-api/internal/retranslator"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	ctx, cancelFunc := context.WithCancel(context.Background())
	sigs := make(chan os.Signal, 1)

	cfg := retranslator.Config{
		ChannelSize:   512,
		ConsumerCount: 2,
		ConsumeSize:   10,
		ProducerCount: 28,
		WorkerCount:   2,
	}

	retranslator := retranslator.NewRetranslator(cfg)
	retranslator.Start(ctx)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	<-sigs
	cancelFunc()

}

