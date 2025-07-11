package notifier

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/lib/pq"
	"github.com/t3mp14r3/shiny-umbrella/internal/cron"
	"go.uber.org/zap"
)

type Notifier struct {
	listener *pq.Listener
	failed   chan error
    logger   *zap.Logger
    cron     *cron.Cron
}

func New(logger *zap.Logger, cron *cron.Cron, connString string) *Notifier {
    n := &Notifier{failed: make(chan error, 2), logger: logger, cron: cron}

	listener := pq.NewListener(
		connString,
		10*time.Second, time.Minute,
		n.logListener)

	if err := listener.Listen("inserts"); err != nil {
		listener.Close()
        log.Fatalf("Failed to start a notifier inserts listener: %v\n", err)
		return nil
	}
	if err := listener.Listen("updates"); err != nil {
		listener.Close()
        log.Fatalf("Failed to start a notifier updates listener: %v\n", err)
		return nil
	}
	if err := listener.Listen("deletes"); err != nil {
		listener.Close()
        log.Fatalf("Failed to start a notifier deletes listener: %v\n", err)
		return nil
	}
	if err := listener.Listen("new"); err != nil {
		listener.Close()
        log.Fatalf("Failed to start a notifier new listener: %v\n", err)
		return nil
	}

	n.listener = listener
	return nil
}

func (n *Notifier) logListener(event pq.ListenerEventType, err error) {
	if err != nil {
		n.logger.Error("Notifier listening failure", zap.Error(err))
	}
	if event == pq.ListenerEventConnectionAttemptFailed {
		n.failed <- err
	}
}

func (n *Notifier) Listen(ctx context.Context) error {
	for {
		select {
		case e := <-n.listener.Notify:
			if e == nil {
				continue
			}

            id, err := strconv.ParseInt(e.Extra, 10, 64)

            if err != nil {
                n.logger.Error("Failed to parse record ID", zap.Error(err))
                continue
            }

            if e.Channel == "new" {
                n.cron.Regular(id)
            } else {
                n.cron.Update(id, e.Channel)
            }
		case err := <-n.failed:
			return err
		case <-time.After(time.Minute):
			go n.listener.Ping()
        case <-ctx.Done():
            return nil
		}
	}
}

func (n *Notifier) Close() error {
	close(n.failed)
	return n.listener.Close()
}
