package views

type Building struct {
	mp *Map
	name string
	image string
	kind string
	pos Position
}

func NewBuilding(m *Map) *Building {
	return &Building{
		mp:m,
		name: "",
		image: "",
		kind: "",
		pos: Position{
			x:0,
			y:0},
	}
}
