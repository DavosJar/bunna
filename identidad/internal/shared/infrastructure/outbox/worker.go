package outbox

import (
	"context"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

// OutboxWorker es un worker background que lee eventos pendientes del outbox
// y los publica en Kafka de forma asíncrona.
type OutboxWorker struct {
	outboxRepo OutboxRepositorio
	kafkaWriter *kafka.Writer
	interval    time.Duration
	maxRetries  int
	batchSize   int
}

func NewOutboxWorker(outboxRepo OutboxRepositorio, kafkaWriter *kafka.Writer, interval time.Duration) *OutboxWorker {
	return &OutboxWorker{
		outboxRepo:  outboxRepo,
		kafkaWriter: kafkaWriter,
		interval:    interval,
		maxRetries:  3,
		batchSize:   50,
	}
}

// Start inicia el worker en una goroutine. Se detiene cuando ctx se cancela.
func (w *OutboxWorker) Start(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(w.interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				w.flush(ctx)
				return
			case <-ticker.C:
				w.processBatch(ctx)
			}
		}
	}()
}

func (w *OutboxWorker) processBatch(ctx context.Context) {
	eventos, err := w.outboxRepo.GetPending(ctx, w.batchSize)
	if err != nil || len(eventos) == 0 {
		return
	}

	for _, evento := range eventos {
		msg := kafka.Message{
			Key:   []byte(evento.AggregateID),
			Value: []byte(evento.Payload),
		}

		// Intentar publicar con timeout para no colgar el worker
		publishCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
		err := w.kafkaWriter.WriteMessages(publishCtx, msg)
		cancel()

		if err != nil {
			log.Printf("[OutboxWorker] Error publicando evento %s: %v", evento.ID, err)
			_ = w.outboxRepo.MarkFailed(ctx, evento.ID, err.Error())
		} else {
			_ = w.outboxRepo.MarkPublished(ctx, evento.ID)
		}
	}
}

// flush publica todos los eventos pendientes antes de cerrar el worker.
func (w *OutboxWorker) flush(ctx context.Context) {
	for {
		eventos, err := w.outboxRepo.GetPending(ctx, w.batchSize)
		if err != nil || len(eventos) == 0 {
			return
		}
		for _, evento := range eventos {
			publishCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
			err := w.kafkaWriter.WriteMessages(publishCtx, kafka.Message{
				Key:   []byte(evento.AggregateID),
				Value: []byte(evento.Payload),
			})
			cancel()
			if err != nil {
				log.Printf("[OutboxWorker] Error en flush: %v", err)
				return
			}
			_ = w.outboxRepo.MarkPublished(ctx, evento.ID)
		}
	}
}
