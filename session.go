// Copyright 2016 zxfonline@sina.com. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package netmsg

import (
	"errors"
	"sync"
	"time"
)

const (
	DEFAULT_ACCESS_TIMEOUT = 60 * time.Second
)

type PipeMsg struct {
	Data  interface{}
	Error error
}
type Handler struct {
	Pipe chan *PipeMsg
}

type Session struct {
	Id     int64
	Hander *Handler
}

var (
	idPool   int64
	sessions map[int64]*Session
	m        *sync.RWMutex
)

func init() {
	sessions = make(map[int64]*Session)
	m = new(sync.RWMutex)
}

func (s *Session) Write(data *PipeMsg) {
	s.Hander.Pipe <- data
}

func (s *Session) Read(timeout time.Duration) *PipeMsg {
	if timeout > 0 {
		blockChan := time.NewTimer(timeout)
		select {
		case msg := <-s.Hander.Pipe:
			blockChan.Stop()
			return msg
		case <-blockChan.C:
			return &PipeMsg{Error: errors.New("timeout")}
		}
	} else {
		return <-s.Hander.Pipe
	}
}

func NewSession() int64 {
	m.Lock()
	defer m.Unlock()
	h := Handler{
		Pipe: make(chan *PipeMsg),
	}
	idPool++
	s := Session{
		Id:     idPool,
		Hander: &h,
	}
	sessions[s.Id] = &s
	return s.Id
}

func GetSession(id int64) *Session {
	m.RLock()
	defer m.RUnlock()
	if s, ok := sessions[id]; ok {
		return s
	}
	return nil
}

func DelSession(id int64) {
	m.Lock()
	defer m.Unlock()
	delete(sessions, id)
}
