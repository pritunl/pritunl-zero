package utils

import (
	"gopkg.in/mgo.v2/bson"
	"sort"
)

type ObjectIdSlice []bson.ObjectId

func (o ObjectIdSlice) Len() int {
	return len(o)
}

func (o ObjectIdSlice) Less(i, j int) bool {
	return o[i] < o[j]
}

func (o ObjectIdSlice) Swap(i, j int) {
	o[i], o[j] = o[j], o[i]
}

func SortObjectIds(x []bson.ObjectId) {
	sort.Sort(ObjectIdSlice(x))
}
