package models

type Page struct {
	StringMap       map[string]string
	IntMap          map[string]int
	FloatMap        map[string]float64
	Data            map[string]any
	IsAuthenticated bool
	IsAuthor        bool
	IsAdmin         bool
}
