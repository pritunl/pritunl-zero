package bastion

import "github.com/pritunl/mongo-go-driver/bson/primitive"

var (
	state = map[primitive.ObjectID]*Bastion{}
)

func Get(authrId primitive.ObjectID) (bast *Bastion) {
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

func New(authrId primitive.ObjectID) (bast *Bastion) {
	bast = &Bastion{
		Authority: authrId,
	}

	state[authrId] = bast

	return
}
