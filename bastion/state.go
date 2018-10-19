package bastion

import "gopkg.in/mgo.v2/bson"

var (
	state = map[bson.ObjectId]*Bastion{}
)

func Get(authrId bson.ObjectId) (bast *Bastion) {
	bast = state[authrId]

	if bast != nil && !bast.State() {
		delete(state, authrId)
		bast = nil
	}

	return
}

func GetAll() (basts []*Bastion) {
	basts = []*Bastion{}

	for authrId, bast := range state {
		if !bast.State() {
			delete(state, authrId)
			continue
		}

		basts = append(basts, bast)
	}

	return
}

func New(authrId bson.ObjectId) (bast *Bastion) {
	bast = &Bastion{
		Authority: authrId,
	}

	state[authrId] = bast

	return
}
