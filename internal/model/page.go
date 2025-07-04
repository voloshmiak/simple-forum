package model

type Page struct {
	StringMap       map[string]string
	IntMap          map[string]int
	FloatMap        map[string]float64
	Data            map[string]any
	Flash           string
	Warning         string
	Error           string
	IsAuthenticated bool
	IsAuthor        bool
	IsAdmin         bool
	CSRFToken       string
}
