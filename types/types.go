package types

import "fmt"

type NameT string

type IdT string

type IdMapT = map[NameT]IdT

type GFileT struct {
	Id   IdT
	Name NameT
}

func (g GFileT) String() (_ string) {
	return fmt.Sprintf("%s (%s)", g.Name, g.Id)
}

type LocalGFileT struct {
	GFileT
	Path string
}

func (l LocalGFileT) String() (_ string) {
	return fmt.Sprintf("%v: %s", l.GFileT, l.Path)
}
