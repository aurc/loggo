package db

import (
	"fmt"
	"os"
	"path"

	"github.com/google/uuid"
	"github.com/tidwall/buntdb"
)

type LogDB struct {
	db        *buntdb.DB
	fileName  string
	feeder    <-chan string
	isFeeding bool
}

func MakeLogDB(name string, feeder <-chan string) *LogDB {
	db := &LogDB{
		feeder: feeder,
	}
	if len(name) > 0 {
		db.fileName = name
	} else {
		tmpDir := os.TempDir()
		fileName := uuid.New().String() + ".db"
		db.fileName = path.Join(tmpDir, fileName)
	}
	var err error
	db.db, err = buntdb.Open(db.fileName)
	if err != nil {
		panic(err)
	}
	db.startFeedLoop()
	return db
}

func (l *LogDB) Close() {
	l.isFeeding = false
}

func (l *LogDB) StreamResults() <-chan string {
	//l.db.View(func(tx *buntdb.Tx) error {
	//	tx.AscendGreaterOrEqual()
	//})
	return nil
}

func (l *LogDB) startFeedLoop() {
	l.isFeeding = true
	go func() {
		for i := int64(1); l.isFeeding; i++ {
			row := <-l.feeder
			if len(row) > 0 {
				l.db.Update(func(tx *buntdb.Tx) error {
					key := fmt.Sprintf("%d", i)
					tx.Set(key, row, nil)
					return nil
				})
			}
		}
	}()
}
