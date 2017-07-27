package log

import (
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/utils"
	"gopkg.in/mgo.v2/bson"
)

func Get(db *database.Database, logId bson.ObjectId) (
	entry *Entry, err error) {

	coll := db.Logs()
	entry = &Entry{}

	err = coll.FindOneId(logId, entry)
	if err != nil {
		return
	}

	return
}

func GetAll(db *database.Database, query *bson.M, page, pageCount int) (
	entries []*Entry, count int, err error) {

	coll := db.Logs()
	entries = []*Entry{}

	qury := coll.Find(query)

	count, err = qury.Count()
	if err != nil {
		err = database.ParseError(err)
		return
	}

	skip := utils.Min(page*pageCount, utils.Max(0, count-pageCount))

	cursor := qury.Sort("-$natural").Skip(skip).Limit(pageCount).Iter()

	entry := &Entry{}
	for cursor.Next(entry) {
		entries = append(entries, entry)
		entry = &Entry{}
	}

	err = cursor.Close()
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}
