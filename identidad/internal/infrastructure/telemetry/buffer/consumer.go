package buffer

import (
	"context"
	"log"
	"time"
)

// StartConsumer launches a background goroutine that drains the ring buffer
// and publishes events via the KafkaProducer.
func StartConsumer(ctx context.Context, buf *RingBuffer, producer *KafkaProducer, cfg Config) {
	consumerAliveGauge.Set(1)
	go func() {
		defer consumerAliveGauge.Set(0)
		flushInterval := time.Duration(cfg.FlushIntervalSeconds) * time.Second
		if flushInterval <= 0 {
			flushInterval = 500 * time.Millisecond
		}
		ticker := time.NewTicker(flushInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				drainRemaining(ctx, buf, producer, cfg)
				producer.Close()
				return
			case <-ticker.C:
			}

			fill := buf.FillRatio()
			// Adaptive batch size based on fill ratio
			batchSize := cfg.BatchSize
			if fill > 0.95 {
				batchSize = 1
			} else if fill > 0.85 {
				batchSize = batchSize / 4
				if batchSize < 1 {
					batchSize = 1
				}
			}

			batch := buf.ReadBatch(batchSize)
			if batch == nil {
				continue
			}

			// Publish with retries and exponential backoff
			var err error
			for attempt := 0; attempt <= cfg.MaxRetries || cfg.MaxRetries == 0; attempt++ {
				if attempt > 0 {
					// Exponential backoff: base^attempt, capped at max
					backoff := cfg.BackoffBase
					for i := 0; i < attempt && backoff < cfg.BackoffMax; i++ {
						backoff *= 2
					}
					if backoff > cfg.BackoffMax {
						backoff = cfg.BackoffMax
					}
					time.Sleep(backoff)
				}
				err = producer.Publish(batch)
				if err == nil {
					break
				}
				publishErrorsCounter.Inc()
			}
			if err != nil {
				log.Printf("telemetry: discarding batch after %d retries: %v", cfg.MaxRetries, err)
				eventsDroppedCounter.Add(float64(len(batch)))
			}
		}
	}()
}

func drainRemaining(ctx context.Context, buf *RingBuffer, producer *KafkaProducer, cfg Config) {
	drainCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	for {
		select {
		case <-drainCtx.Done():
			return
		default:
		}
		batch := buf.ReadBatch(cfg.BatchSize)
		if batch == nil {
			return
		}
		if err := producer.Publish(batch); err != nil {
			log.Printf("telemetry: error draining batch: %v", err)
			eventsDroppedCounter.Add(float64(len(batch)))
		}
	}
}
