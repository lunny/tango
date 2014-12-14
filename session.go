package tango

import (
	"time"

	"github.com/go-xweb/httpsession"
)

type SessionInterface interface {
	SetSession(*httpsession.Session)
}

type Sessions struct {
	*httpsession.Manager
}

func NewSessions(sessionTimeout time.Duration) *Sessions {
	sessionMgr := httpsession.Default()
	if sessionTimeout > time.Second {
		sessionMgr.SetMaxAge(sessionTimeout)
	}
	sessionMgr.Run()

	return &Sessions{Manager: sessionMgr}
}

func (itor *Sessions) Handle(ctx *Context) {
	if action := ctx.Action(); ctx != nil {
		if s, ok := action.(SessionInterface); ok {
			session := itor.Session(ctx.Req(), ctx.ResponseWriter)
			s.SetSession(session)
		}
	}

	ctx.Next()
}
