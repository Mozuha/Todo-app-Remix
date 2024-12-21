package middlewares_mock

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type MockSessionStore struct {
	sessions.Store
	saveErr error
}

func (m *MockSessionStore) Get(ctx *gin.Context, name string) (sessions.Session, error) {
	session := &MockSession{
		data:    make(map[interface{}]interface{}),
		saveErr: m.saveErr,
	}
	return session, nil
}

type MockSession struct {
	data    map[interface{}]interface{}
	saveErr error
	id      string
	options *sessions.Options
}

func NewMockSession(saveErr error) *MockSession {
	return &MockSession{
		data:    make(map[interface{}]interface{}),
		saveErr: saveErr,
		id:      "mock-session-id",
		options: &sessions.Options{},
	}
}

func (m *MockSession) Get(key interface{}) interface{} {
	return m.data[key]
}

func (m *MockSession) Set(key, value interface{}) {
	m.data[key] = value
}

func (m *MockSession) Delete(key interface{}) {
	delete(m.data, key)
}

func (m *MockSession) Clear() {
	m.data = make(map[interface{}]interface{})
}

func (m *MockSession) AddFlash(value interface{}, vars ...string) {
}

func (m *MockSession) Flashes(vars ...string) []interface{} {
	return nil
}

func (m *MockSession) Options(options sessions.Options) {
	m.options = &options
}

func (m *MockSession) Save() error {
	return m.saveErr
}

func (m *MockSession) ID() string {
	return m.id
}

func (m *MockSession) Values() map[interface{}]interface{} {
	return m.data
}
