package counter

import (
	"fmt"
	"io"
	"os"
	"sync/atomic"

	"github.com/mattn/go-isatty"
)

const (
	objectTypes   = "object types"
	permissions   = "permissions"
	relationTypes = "relation types"
	objects       = "objects"
	relations     = "relations"
	skipped       = "skipped"
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
	Name    string
	value   int64
	skipped int64
}

func (c *Item) Incr() *Item {
	atomic.AddInt64(&c.value, 1)
	return c
}

func (c *Item) Skip() *Item {
	atomic.AddInt64(&c.skipped, 1)
	return c
}

func (c *Item) Print(w io.Writer) {
	if isatty.IsTerminal(os.Stdout.Fd()) {
		fmt.Fprintf(w, "\033[2K\r%15s: %d", c.Name, c.value)
	}
}

func (c *Counter) ObjectTypes() *Item {
	if c.objectTypes == nil {
		c.objectTypes = &Item{Name: objectTypes}
	}
	return c.objectTypes
}

func (c *Counter) Permissions() *Item {
	if c.permissions == nil {
		c.permissions = &Item{Name: permissions}
	}
	return c.permissions
}

func (c *Counter) RelationTypes() *Item {
	if c.relationTypes == nil {
		c.relationTypes = &Item{Name: relationTypes}
	}
	return c.relationTypes
}

func (c *Counter) Objects() *Item {
	if c.objects == nil {
		c.objects = &Item{Name: objects}
	}
	return c.objects
}

func (c *Counter) Relations() *Item {
	if c.relations == nil {
		c.relations = &Item{Name: relations}
	}
	return c.relations
}

func (c *Counter) Print(w io.Writer) {
	if isatty.IsTerminal(os.Stdout.Fd()) {
		fmt.Fprintf(w, "\033[2K\r")
	}

	if c.objectTypes != nil {
		fmt.Fprintf(w, "%15s %d\n", objectTypes, c.objectTypes.value)
	} else {
		fmt.Fprintf(w, "%15s %s\n", objectTypes, skipped)
	}

	if c.permissions != nil {
		fmt.Fprintf(w, "%15s %d\n", permissions, c.permissions.value)
	} else {
		fmt.Fprintf(w, "%15s %s\n", permissions, skipped)
	}

	if c.relationTypes != nil {
		fmt.Fprintf(w, "%15s %d\n", relationTypes, c.relationTypes.value)
	} else {
		fmt.Fprintf(w, "%15s %s\n", relationTypes, skipped)
	}

	if c.objects != nil {
		msg := ""
		if c.objects.skipped > 0 {
			msg = " WARNING data contained unknown fields"
		}
		fmt.Fprintf(w, "%15s %d%s\n", objects, c.objects.value, msg)
	} else {
		fmt.Fprintf(w, "%15s %s\n", objects, skipped)
	}

	if c.relations != nil {
		msg := ""
		if c.objects.skipped > 0 {
			msg = " WARNING data contained unknown fields"
		}
		fmt.Fprintf(w, "%15s %d%s\n", relations, c.relations.value, msg)
	} else {
		fmt.Fprintf(w, "%15s %s\n", relations, skipped)
	}
}
