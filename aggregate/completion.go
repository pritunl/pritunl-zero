package aggregate

import (
	"github.com/pritunl/mongo-go-driver/v2/bson"
	"github.com/pritunl/mongo-go-driver/v2/mongo/options"
	"github.com/pritunl/pritunl-zero/certificate"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/node"
	"github.com/pritunl/pritunl-zero/secret"
)

type Completion struct {
	Nodes        []*node.Node               `json:"nodes"`
	Certificates []*certificate.Certificate `json:"certificates"`
	Secrets      []*secret.Secret           `json:"secrets"`
}

func get(db *database.Database, coll *database.Collection,
	query, projection *bson.M, new func() any,
	add func(any)) (err error) {

	cursor, err := coll.Find(db, query, options.Find().
		SetProjection(projection))
	if err != nil {
		err = database.ParseError(err)
		return
	}
	defer cursor.Close(db)

	for cursor.Next(db) {
		item := new()
		err = cursor.Decode(item)
		if err != nil {
			err = database.ParseError(err)
			return
		}

		add(item)
	}

	return
}

func GetCompletion(db *database.Database, orgId bson.ObjectID) (
	cmpl *Completion, err error) {

	cmpl = &Completion{}
	query := &bson.M{}

	err = get(
		db,
		db.Nodes(),
		&bson.M{},
		&bson.M{
			"_id":              1,
			"name":             1,
			"zone":             1,
			"types":            1,
			"timestamp":        1,
			"cpu_units":        1,
			"memory_units":     1,
			"cpu_units_res":    1,
			"memory_units_res": 1,
		},
		func() any {
			return &node.Node{}
		},
		func(item any) {
			nde := item.(*node.Node)

			cmpl.Nodes = append(
				cmpl.Nodes,
				nde,
			)
		},
	)
	if err != nil {
		return
	}

	err = get(
		db,
		db.Certificates(),
		query,
		&bson.M{
			"_id":          1,
			"name":         1,
			"organization": 1,
			"type":         1,
		},
		func() any {
			return &certificate.Certificate{}
		},
		func(item any) {
			cmpl.Certificates = append(
				cmpl.Certificates,
				item.(*certificate.Certificate),
			)
		},
	)
	if err != nil {
		return
	}

	err = get(
		db,
		db.Secrets(),
		query,
		&bson.M{
			"_id":          1,
			"name":         1,
			"organization": 1,
			"type":         1,
		},
		func() any {
			return &secret.Secret{}
		},
		func(item any) {
			cmpl.Secrets = append(
				cmpl.Secrets,
				item.(*secret.Secret),
			)
		},
	)
	if err != nil {
		return
	}

	return
}
