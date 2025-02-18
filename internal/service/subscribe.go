package service

import (
	"math/rand"
	"time"

	"github.com/gorilla/websocket"
	"github.com/itering/substrate-api-rpc/rpc"
	ws "github.com/itering/substrate-api-rpc/websocket"
	"golang.org/x/exp/slog"
)

const (
	runtimeVersion = iota + 1
	newHeader
	finalizeHeader
)

func logError(msg string, err error) {
	if err != nil {
		slog.Error(msg, "error", err)
	}
}

func (s *Service) Subscribe(conn ws.WsConn, stop chan struct{}) {
	var err error

	defer conn.Close()

	dead := make(chan struct{}, 1)
	reconnected := make(chan struct{}, 1)

	subscribeSrv := s.initSubscribeService(stop)
	go func() {
		defer close(dead)
		defer close(reconnected)
		waitForReconnect := false
		for {
			select {
			case <-stop:
				return
			default:
				if !conn.IsConnected() {
					time.Sleep(time.Second * 1)
					continue
				}
				if waitForReconnect {
					waitForReconnect = false
					<-reconnected
					time.Sleep(time.Second * 10)
				}
				if !conn.IsConnected() {
					time.Sleep(time.Second * 1)
					continue
				}
				_, message, err := conn.ReadMessage()
				if err != nil {
					logError("read failed", err)
					if len(dead) == 0 {
						dead <- struct{}{}
						waitForReconnect = true
					}
					continue
				}
				_ = subscribeSrv.parser(message)
			}
		}
	}()

	for !conn.IsConnected() {
		time.Sleep(time.Second)
	}

	if err = conn.WriteMessage(websocket.TextMessage, rpc.ChainGetRuntimeVersion(runtimeVersion)); err != nil {
		logError("write failed", err)
	}
	if err = conn.WriteMessage(websocket.TextMessage, rpc.ChainSubscribeNewHead(newHeader)); err != nil {
		logError("write failed", err)
	}
	if err = conn.WriteMessage(websocket.TextMessage, rpc.ChainSubscribeFinalizedHeads(finalizeHeader)); err != nil {
		logError("write failed", err)
	}

	ticker := time.NewTicker(time.Second * 3)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if !conn.IsConnected() {
				slog.Debug("connection is not connected")
				continue
			}
			if err := conn.WriteMessage(websocket.TextMessage, rpc.SystemHealth(rand.Intn(100)+finalizeHeader)); err != nil {
				logError("system health get failed", err)
			}
		case <-stop:
			err = conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				logError("write close failed", err)
				return
			}
			return
		case <-dead:
			slog.Warn("connection is dead, reconnecting...")
			conn.CloseAndReconnect()
			for !conn.IsConnected() {
				time.Sleep(time.Second)
			}
			if err = conn.WriteMessage(websocket.TextMessage, rpc.ChainGetRuntimeVersion(runtimeVersion)); err != nil {
				logError("write failed", err)
			}
			if err = conn.WriteMessage(websocket.TextMessage, rpc.ChainSubscribeNewHead(newHeader)); err != nil {
				logError("write failed", err)
			}
			if err = conn.WriteMessage(websocket.TextMessage, rpc.ChainSubscribeFinalizedHeads(finalizeHeader)); err != nil {
				logError("write failed", err)
			}
			reconnected <- struct{}{}
		}
	}
}
