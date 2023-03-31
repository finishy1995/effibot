package scenario

import (
	"fmt"
	"github.com/finishy1995/effibot/library/id"
	"sync"
)

// TODO: move storage from local memory to database

var (
	sessionMutex   sync.RWMutex
	paragraphMutex sync.RWMutex
	sessionMap     = map[id.ID]*session{}
	paragraphMap   = map[id.ID]*paragraph{}
)

func getSession(sessionId id.ID) *session {
	sessionMutex.RLock()
	defer sessionMutex.RUnlock()
	return sessionMap[sessionId]
}

func deleteSession(sessionId id.ID) {
	sessionMutex.Lock()
	defer sessionMutex.Unlock()
	delete(sessionMap, sessionId)
}

func addSession(sessionId id.ID, instance *session) error {
	sessionMutex.Lock()
	defer sessionMutex.Unlock()
	if _, ok := sessionMap[sessionId]; ok {
		return fmt.Errorf("session duplicated, id: %s", sessionId)
	}
	sessionMap[sessionId] = instance
	return nil
}

func getParagraph(paragraphId id.ID) *paragraph {
	paragraphMutex.RLock()
	defer paragraphMutex.RUnlock()
	return paragraphMap[paragraphId]
}

func deleteParagraph(paragraphId id.ID) {
	paragraphMutex.Lock()
	defer paragraphMutex.Unlock()
	delete(paragraphMap, paragraphId)
}

func addParagraph(paragraphId id.ID, instance *paragraph) error {
	paragraphMutex.Lock()
	defer paragraphMutex.Unlock()
	if _, ok := paragraphMap[paragraphId]; ok {
		return fmt.Errorf("paragraph duplicated, id: %s", paragraphId)
	}
	paragraphMap[paragraphId] = instance
	return nil
}
