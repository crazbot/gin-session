package ginsession_test

import (
	"encoding/json"
	"testing"
	"time"

	ginsession "github.com/eyediy/gin-session"
)

var (
	sessionManager *ginsession.SessionManager
)

func init() {
	sessionManager = ginsession.NewSessionManager("session.properties")
}

type AAA struct {
	ID   int
	Name string
	a    int
	b    int
}

func jsonInterface(bytes []byte, ss *map[string]interface{}) error {
	return json.Unmarshal(bytes, ss)
}

func Test_json(t *testing.T) {
	Data := make(map[string]interface{})
	t.Log(Data)
	Data["abc"] = 1111
	Data["ccb"] = AAA{22, "aaa", 22, 33}
	Data["kele"] = "ssss"
	bytes, err := json.Marshal(Data)
	if err != nil {
		t.Error(err)
	}
	t.Log(bytes)

	var dd interface{}
	dd = &AAA{22, "aaa", 22, 33}
	t.Log(dd)
	bytes, err = json.Marshal(dd)
	if err != nil {
		t.Error(err)
	}
	t.Log(bytes)

	var ss map[string]interface{}
	t.Log(ss)
	//var ss interface{}
	err = jsonInterface(bytes, &ss)
	//err = json.Unmarshal(bytes, &ss)
	if err != nil {
		t.Error(err)
	}
	t.Log(ss)
}

func Test_getSession(t *testing.T) {
	t.Log(sessionManager)
	if sessionManager == nil {
		t.Error("sessionManager failed")
	}
	session := sessionManager.GetSession("jkkaa")
	t.Log(session.ID)
}

func Test_saveSession(t *testing.T) {
	session := sessionManager.GetSession("aaa")
	t.Log(session.ID)
	err := session.Save()
	if err != nil {
		t.Error(err)
	}
}

func Test_completeFlow(t *testing.T) {
	session := sessionManager.GetSession("aaa")
	t.Log(session.ID)
	err := session.Save()
	if err != nil {
		t.Error(err)
	}
	// retrieve the session
	newSession := sessionManager.GetSession(session.ID)
	t.Log(newSession)
	// bind some data
	newSession.Data.Value["signed"] = true
	newSession.Data.Value["userId"] = 558
	err = newSession.Save()
	if err != nil {
		t.Error(err)
	}

	time.Sleep(10 * time.Second)

	// retrieve again
	newSession = sessionManager.GetSession(session.ID)
	if newSession.ID != session.ID {
		t.Error("session expired")
	}
	oldSession := &ginsession.Session{}
	oldSession.Update()

	time.Sleep(1 * time.Second)

	newSession = sessionManager.GetSession(session.ID)
	t.Log(newSession)

}

func Test_copySession(t *testing.T) {
	session := sessionManager.GetSession("aaa")
	t.Log(session.ID)
	oldSession := session.Copy()
	if oldSession == nil {
		t.Error("session.Copy() failed")
	}
	session.Data.Value["test"] = "OK"
	session.Update()
	t.Log(oldSession)
	t.Log(session)
}

func Test_shouldUpdate(t *testing.T) {
	session := sessionManager.GetSession("aaa")
	t.Log(session.ID)
	oldSession := session.Copy()
	if oldSession == nil {
		t.Error("session.Copy() failed")
	}
	session.Data.Value["test"] = "OK"
	if session.ShouldUpdate(oldSession) {
		t.Error("should not update")
	}
	session.Update()
	if !session.ShouldUpdate(oldSession) {
		t.Error("should update")
	}
}

func Test_shouldKeepAlive(t *testing.T) {
	session := sessionManager.GetSession("aaa")
	t.Log(session.ID)
	oldSession := session.Copy()
	if oldSession == nil {
		t.Error("session.Copy() failed")
	}
	session.Data.Value["test"] = "OK"
	if session.ShouldUpdate(oldSession) {
		t.Error("should not update")
	}
	time.Sleep(10 * time.Second) // 60 = 配置里的maxAge
	if !session.ShouldUpdate(oldSession) {
		t.Error("should update")
	}
}
