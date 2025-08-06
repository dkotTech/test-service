package http

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"

	"test-service/events"
	"test-service/helpers"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func WsHandlers(log *slog.Logger, eventsSrv events.Service) http.Handler {
	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		eventKindsStr := r.URL.Query()["event_kind"]

		conn, err := makeSyncConnection(ctx, w, r)
		if err != nil {
			helpers.SendError(ctx, w, err)
			return
		}
		defer func() {
			conn.Close()
		}()

		eventKinds := make([]events.EventKind, 0, len(eventKindsStr))
		for _, kind := range eventKindsStr {
			eventKinds = append(eventKinds, events.EventKind(kind))
		}

		// run a ping worker to track a disconnected clients
		ctx = runPing(ctx, conn)

		eventClient := eventsSrv.RegisterConnection(ctx, eventKinds)
		defer func() { _ = eventClient.Shutdown() }()

		if err := conn.WriteJSON(map[string]interface{}{"ok": "connected", "events": eventClient.EventsSubscribed()}); err != nil {
			log.Error("WriteMessage error", "err", err)
			return
		}

		// send a events untile context done
	loop:
		for {
			select {
			case event := <-eventClient.Subscribe():
				if err := conn.WriteJSON(event); err != nil {
					log.Error("WriteMessage error", "err", err)
					return
				}
			case <-ctx.Done():
				conn.Close()
				log.Info("ws close")
				break loop
			}
		}
	})

	return r
}

func makeSyncConnection(ctx context.Context, w http.ResponseWriter, r *http.Request) (*helpers.SyncConn, error) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}

	return helpers.NewSyncConnection(conn), nil
}

const (
	pongWait   = 2 * time.Second
	pingPeriod = 1 * time.Second
)

func runPing(ctx context.Context, conn *helpers.SyncConn) context.Context {
	ctx, cancel := context.WithCancel(ctx)

	_ = conn.SetReadDeadline(time.Now().Add(pongWait))
	conn.SetPongHandler(func(string) error {
		_ = conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	go func() {
		ticker := time.NewTicker(pingPeriod)
		defer func() {
			ticker.Stop()
		}()

		for range ticker.C {
			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				// if we can not receive a pong msg, cancel a context
				cancel()
				return
			}
		}
	}()

	return ctx
}
