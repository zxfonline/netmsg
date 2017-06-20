// Copyright 2016 zxfonline@sina.com. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package netmsg

import (
	"fmt"
	"time"

	. "github.com/zxfonline/taskexcutor"
)

type Msg struct {
	SessionId int64
	Data      interface{}
	CallBack  CallBackMsg
}

//事件回调
type CallBackMsg func(interface{}) interface{}

func SyncSendMsg(excutor Excutor, sessionid int64, data interface{}, callback CallBackMsg) error {
	return excutor.Excute(NewTaskService(func(params ...interface{}) {
		msg := (params[0]).(Msg)
		if s := GetSession(msg.SessionId); s != nil {
			defer func() {
				if err := recover(); err != nil {
					s.Write(&PipeMsg{Error: fmt.Errorf("%v", err)})
				}
			}()
			s.Write(&PipeMsg{Data: msg.CallBack(msg.Data)})
		}
	}, Msg{SessionId: sessionid, Data: data, CallBack: callback}))
}

func AsyncSendMsg(excutor Excutor, data interface{}, callback CallBackMsg) error {
	return excutor.Excute(NewTaskService(func(params ...interface{}) {
		msg := (params[0]).(Msg)
		msg.CallBack(msg.Data)
	}, Msg{Data: data, CallBack: callback}))
}

func RecMsg(sId int64) interface{} {
	if s := GetSession(sId); s != nil {
		rt := s.Read(0)
		DelSession(sId)
		return rt
	}
	return nil
}

func RecMsgWithTime(sId int64, timeout time.Duration) interface{} {
	if s := GetSession(sId); s != nil {
		rt := s.Read(timeout)
		DelSession(sId)
		return rt
	}
	return nil
}
