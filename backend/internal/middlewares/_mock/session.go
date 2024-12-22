package middlewares_mock

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// This mock is basically for simulating the case of getting error upon session.Save()

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
	data      map[interface{}]interface{}
	saveErr   error
	id        string
	options   *sessions.Options
	dirtyData map[interface{}]interface{}
}

func NewMockSession(saveErr error) *MockSession {
	return &MockSession{
		data:      make(map[interface{}]interface{}),
		saveErr:   saveErr,
		id:        "mock-session-id",
		options:   &sessions.Options{},
		dirtyData: make(map[interface{}]interface{}),
	}
}

func (m *MockSession) Get(key interface{}) interface{} {
	return m.data[key]
}

func (m *MockSession) Set(key, value interface{}) {
	m.dirtyData[key] = value
}

func (m *MockSession) Delete(key interface{}) {
	delete(m.dirtyData, key)
}

func (m *MockSession) Clear() {
	m.dirtyData = make(map[interface{}]interface{})
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
	if m.saveErr != nil {
		m.dirtyData = m.data // simulate rollback upon save failure
		return m.saveErr
	}

	m.data = m.dirtyData
	return nil
}

func (m *MockSession) ID() string {
	return m.id
}

func (m *MockSession) Values() map[interface{}]interface{} {
	return m.data
}
