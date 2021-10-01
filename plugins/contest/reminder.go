package contest

import (
	"bytes"
	"encoding/gob"
	"sync"
)

type reminder struct {
	Contest Contest
	Group   int64
	User    []int64
}

var reminderMutex sync.Mutex
var reminderList []*reminder

func saveReminder() {
	var w bytes.Buffer
	_ = gob.NewEncoder(&w).Encode(reminderList)
	_ = db.Put(reminderDBKey, w.Bytes())
}
