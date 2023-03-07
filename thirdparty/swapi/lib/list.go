package lib

type Entity interface {
	Film | Person
}

type List[E Entity] struct {
	Count   int
	Next    *string
	Results []E
}
