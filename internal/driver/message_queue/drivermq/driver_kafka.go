package drivermq

import (
	"fmt"
	"github.com/segmentio/kafka-go"
	"github.com/shbov/hse-go_final/internal/driver/app"
	"github.com/shbov/hse-go_final/internal/driver/message_queue"
	"go.uber.org/zap"
)

var _ message_queue.MessageQueue = (*driverKafka)(nil)

type driverKafka struct {
	wc kafka.WriterConfig
	rc kafka.ReaderConfig
}

//
//func (r *driverKafka) conn(ctx context.Context) Conn {
//	if tx, ok := ctx.Value(repo.CtxKeyTx).(pgx.Tx); ok {
//		return tx
//	}
//
//	return r.pgxPool
//}

type logWrap struct {
	l *zap.Logger
}

func (wrap logWrap) logf(msg string, a ...interface{}) {
	wrap.l.Info(fmt.Sprintf(msg, a))
}

func New(conf *app.KafkaConfig, log *zap.Logger) (message_queue.MessageQueue, error) {
	w := kafka.WriterConfig{
		Brokers:  conf.Brokers,
		Topic:    conf.Topic,
		Balancer: &kafka.LeastBytes{},
		Logger:   kafka.LoggerFunc(logWrap{log}.logf),
	}

	r := kafka.ReaderConfig{
		Brokers: conf.Brokers,
		GroupID: conf.GroupID,
		Topic:   conf.Topic,
	}

	d := &driverKafka{
		wc: w,
		rc: r,
	}

	log.Info("message_queue successfully created")
	return d, nil
}
