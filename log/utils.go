package log

import (
	"github.com/pritunl/mongo-go-driver/v2/bson"
	"github.com/pritunl/mongo-go-driver/v2/mongo/options"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/event"
)

func Get(db *database.Database, logId bson.ObjectID) (
	entry *Entry, err error) {

	coll := db.Logs()
	entry = &Entry{}

	err = coll.FindOneId(logId, entry)
	if err != nil {
		return
	}

	return
}

func GetAll(db *database.Database, query *bson.M, page, pageCount int64) (
	entries []*Entry, count int64, err error) {

	coll := db.Logs()
	entries = []*Entry{}

	if len(*query) == 0 {
		count, err = coll.EstimatedDocumentCount(db)
		if err != nil {
			err = database.ParseError(err)
			return
		}
	} else {
		count, err = coll.CountDocuments(db, query)
		if err != nil {
			err = database.ParseError(err)
			return
		}
	}

	opts := options.Find().
		SetSort(bson.D{{"$natural", -1}})

	if pageCount != 0 {
		if pageCount == 0 {
			pageCount = 20
		}
		maxPage := count / pageCount
		if count == pageCount {
			maxPage = 0
		}
		page = min(page, maxPage)
		skip := min(page*pageCount, count)
		opts.SetSkip(skip).SetLimit(pageCount)
	}

	cursor, err := coll.Find(
		db,
		query,
		opts,
	)
	if err != nil {
		err = database.ParseError(err)
		return
	}
	defer cursor.Close(db)

	for cursor.Next(db) {
		entry := &Entry{}
		err = cursor.Decode(entry)
		if err != nil {
			err = database.ParseError(err)
			return
		}

		entries = append(entries, entry)
	}

	err = cursor.Err()
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}

func Clear(db *database.Database) (err error) {
	coll := db.Logs()

	_, err = coll.DeleteMany(db, &bson.M{})
	if err != nil {
		err = database.ParseError(err)
		return
	}

	_ = event.PublishDispatch(db, "log.change")

	return
}
