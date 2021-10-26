package retranslator

import (
	"github.com/ozonmp/edu-solution-api/internal/mocks"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
)

func TestStart(t *testing.T) {

	ctrl := gomock.NewController(t)
	repo := mocks.NewMockEventRepo(ctrl)
	sender := mocks.NewMockEventSender(ctrl)

	repo.EXPECT().Lock(gomock.Any()).AnyTimes()

	cfg := Config{
		ChannelSize:   512,
		ConsumerCount: 12,
		ConsumeSize:   10,
		ConsumeTimeout: time.Second,
		ProducerCount: 200,
		WorkerCount:   20,
		Repo:          repo,
		Sender:        sender,
	}

	retranslator := NewRetranslator(cfg)
	retranslator.Start()
	time.Sleep(time.Second * 20)
	retranslator.Close()
}
