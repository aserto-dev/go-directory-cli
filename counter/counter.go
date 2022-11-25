package counter

import (
	"fmt"
	"io"
	"os"
	"sync/atomic"

	"github.com/mattn/go-isatty"
)

type Counter struct {
	objectTypes   *Item
	permissions   *Item
	relationTypes *Item
	objects       *Item
	relations     *Item
}

func New() *Counter {
	return &Counter{}
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
	if isatty.IsTerminal(os.Stdout.Fd()) {
		fmt.Fprintf(w, "\033[2K\r%15s: %d", c.Name, c.value)
	}
}

func (c *Counter) ObjectTypes() *Item {
	if c.objectTypes == nil {
		c.objectTypes = &Item{Name: "object types"}
	}
	return c.objectTypes
}

func (c *Counter) Permissions() *Item {
	if c.permissions == nil {
		c.permissions = &Item{Name: "permissions"}
	}
	return c.permissions
}

func (c *Counter) RelationTypes() *Item {
	if c.relations == nil {
		c.relations = &Item{Name: "relations"}
	}
	return c.relations
}

func (c *Counter) Objects() *Item {
	if c.objects == nil {
		c.objects = &Item{Name: "objects"}
	}
	return c.objects
}

func (c *Counter) Relations() *Item {
	if c.relations == nil {
		c.relations = &Item{Name: "relations"}
	}
	return c.relations
}

func (c *Counter) Print(w io.Writer) {
	if isatty.IsTerminal(os.Stdout.Fd()) {
		fmt.Fprintf(w, "\033[2K\r")
	}

	if c.objectTypes != nil {
		fmt.Fprintf(w, "%15s %d\n", "object types:", c.objectTypes.value)
	} else {
		fmt.Fprintf(w, "%15s %s\n", "object types:", "skipped")
	}

	if c.permissions != nil {
		fmt.Fprintf(w, "%15s %d\n", "permissions:", c.permissions.value)
	} else {
		fmt.Fprintf(w, "%15s %s\n", "permissions:", "skipped")
	}

	if c.relationTypes != nil {
		fmt.Fprintf(w, "%15s %d\n", "relation types:", c.relationTypes.value)
	} else {
		fmt.Fprintf(w, "%15s %s\n", "relation types:", "skipped")
	}

	if c.objects != nil {
		fmt.Fprintf(w, "%15s %d\n", "objects:", c.objects.value)
	} else {
		fmt.Fprintf(w, "%15s %s\n", "objects:", "skipped")
	}

	if c.relations != nil {
		fmt.Fprintf(w, "%15s %d\n", "relations:", c.relations.value)
	} else {
		fmt.Fprintf(w, "%15s %s\n", "relations:", "skipped")
	}
}
