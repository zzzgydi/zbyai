package thread

import (
	"errors"
	"sync"
)

var threadMutex sync.Mutex
var threadWorkMap = make(map[string]bool)

func LockThread(threadId string) error {
	threadMutex.Lock()
	defer threadMutex.Unlock()

	if _, ok := threadWorkMap[threadId]; ok {
		return errors.New("thread is working")
	}

	threadWorkMap[threadId] = true
	return nil
}

func UnlockThread(threadId string) {
	threadMutex.Lock()
	defer threadMutex.Unlock()

	delete(threadWorkMap, threadId)
}
