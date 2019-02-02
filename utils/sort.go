package utils

import (
	"sort"

	"github.com/pritunl/mongo-go-driver/bson/primitive"
)

type ObjectIdSlice []primitive.ObjectID

func (o ObjectIdSlice) Len() int {
	return len(o)
}

func (o ObjectIdSlice) Less(i, j int) bool {
	return o[i].Hex() < o[j].Hex()
}

func (o ObjectIdSlice) Swap(i, j int) {
	o[i], o[j] = o[j], o[i]
}

func SortObjectIds(x []primitive.ObjectID) {
	sort.Sort(ObjectIdSlice(x))
}
