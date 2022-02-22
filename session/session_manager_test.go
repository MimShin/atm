package session

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	hashLen = 64
	userID  = "user-0001"
	pin     = "1234"
)

func TestCRUD(t *testing.T) {
	sm := NewSessionManager()
	require.IsType(t, sm, &SessionManager{})
	require.NotNil(t, sm.sessions)

	// fetch non-existing session
	sessionID, ok := sm.Get(userID)
	assert.False(t, ok)
	assert.Empty(t, sessionID)

	// load-or-store new session
	origSessID := sm.LoadOrStore(userID, pin)
	require.Len(t, origSessID, hashLen)

	// load-or-store existing session -> shouldn't overwrite the existing session-id
	sessID := sm.LoadOrStore(userID, "a-new-pin")
	assert.Equal(t, origSessID, sessID)

	// fetch existing session
	sessID, ok = sm.Get(userID)
	assert.True(t, ok)
	assert.Equal(t, origSessID, sessID)

	// delete session
	sm.Delete(userID)
	sessID, ok = sm.Get(userID)
	assert.False(t, ok)
	assert.Empty(t, sessID)
}

func TestCheck(t *testing.T) {
	sm := NewSessionManager()
	require.IsType(t, sm, &SessionManager{})
	require.NotNil(t, sm.sessions)

	// check non-existing session
	ok := sm.Check("non-existing-userID", "anything")
	assert.False(t, ok)

	// store new session
	sessID := sm.LoadOrStore(userID, pin)
	require.Len(t, sessID, hashLen)

	// check wrong sessionID
	ok = sm.Check(userID, "wrong")
	assert.False(t, ok)

	// check existing session
	ok = sm.Check(userID, sessID)
	assert.True(t, ok)

	// delete and check
	sm.Delete(userID)
	ok = sm.Check(userID, sessID)
	assert.False(t, ok)
}
