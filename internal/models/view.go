package models

type ViewData struct {
	StringMap       map[string]string
	IntMap          map[string]int
	FloatMap        map[string]float64
	Data            map[string]any
	IsAuthenticated bool
}
