package bastion

import "github.com/pritunl/mongo-go-driver/v2/bson"

var (
	state = map[bson.ObjectID]*Bastion{}
)

func Get(authrId bson.ObjectID) (bast *Bastion) {
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

func New(authrId bson.ObjectID) (bast *Bastion) {
	bast = &Bastion{
		Authority: authrId,
	}

	state[authrId] = bast

	return
}
