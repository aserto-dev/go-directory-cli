package counter

import (
	"fmt"
	"io"
	"sync/atomic"
)

type Counter struct {
	ObjectTypes   *Item
	Permissions   *Item
	RelationTypes *Item
	Objects       *Item
	Relations     *Item
}

func New() *Counter {
	return &Counter{
		ObjectTypes:   &Item{Name: "object types"},
		Permissions:   &Item{Name: "permissions"},
		RelationTypes: &Item{Name: "relations"},
		Objects:       &Item{Name: "objects"},
		Relations:     &Item{Name: "relations"},
	}
}

type Item struct {
	Name  string
	value int64
}

func (c *Item) Incr() *Item {
	atomic.AddInt64(&c.value, 1)
	return c
}

func (c *Item) Print(w io.Writer) {
	fmt.Fprintf(w, "\033[2K\r%15s: %d", c.Name, c.value)
}

func (c *Counter) Print(w io.Writer) {
	fmt.Fprintf(w, "\033[2K\r")
	fmt.Fprintf(w, "%15s %d\n", "object types:", c.ObjectTypes.value)
	fmt.Fprintf(w, "%15s %d\n", "permissions:", c.Permissions.value)
	fmt.Fprintf(w, "%15s %d\n", "relation types:", c.RelationTypes.value)
	fmt.Fprintf(w, "%15s %d\n", "objects:", c.Objects.value)
	fmt.Fprintf(w, "%15s %d\n", "relations:", c.Relations.value)
}
