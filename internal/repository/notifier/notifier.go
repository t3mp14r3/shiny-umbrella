package notifier

import (
	"log"
	"strconv"
	"time"

	"github.com/lib/pq"
	"github.com/t3mp14r3/shiny-umbrella/internal/cron"
	"go.uber.org/zap"
)

// notifier encapsulates the state of the listener connection.
type Notifier struct {
	listener *pq.Listener
	failed   chan error
    logger   *zap.Logger
    cron     *cron.Cron
}

// newNotifier creates a new notifier for given PostgreSQL credentials.
func New(logger *zap.Logger, cron *cron.Cron, connString string) (*Notifier, error) {
    n := &Notifier{failed: make(chan error, 2), logger: logger, cron: cron}

	listener := pq.NewListener(
		connString,
		10*time.Second, time.Minute,
		n.logListener)

	if err := listener.Listen("inserts"); err != nil {
		listener.Close()
        log.Printf("Failed to start a notifier inserts listener: %v\n", err)
		return nil, err
	}
	if err := listener.Listen("updates"); err != nil {
		listener.Close()
        log.Printf("Failed to start a notifier inserts listener: %v\n", err)
		return nil, err
	}
	if err := listener.Listen("deletes"); err != nil {
		listener.Close()
        log.Printf("Failed to start a notifier inserts listener: %v\n", err)
		return nil, err
	}

	n.listener = listener
	return n, nil
}

// logListener is the state change callback for the listener.
func (n *Notifier) logListener(event pq.ListenerEventType, err error) {
	if err != nil {
		n.logger.Error("Notifier listening failure", zap.Error(err))
	}
	if event == pq.ListenerEventConnectionAttemptFailed {
		n.failed <- err
	}
}

// fetch is the main loop of the notifier to receive data from
// the database in JSON-FORMAT and send it down the send channel.
func (n *Notifier) Listen() error {
	for {
		select {
		case e := <-n.listener.Notify:
			if e == nil {
				continue
			}

            id, err := strconv.ParseInt(e.Extra, 10, 64)

            if err != nil {
                n.logger.Error("Failed to parse automatic record ID", zap.Error(err))
                continue
            }

            n.cron.Update(id, e.Channel)
		case err := <-n.failed:
			return err
		case <-time.After(time.Minute):
			go n.listener.Ping()
		}
	}
}

func (n *Notifier) Shutdown() error {
	close(n.failed)
	return n.listener.Close()
}
