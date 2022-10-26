package client

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aserto-dev/go-directory-cli/counter"
	dsc "github.com/aserto-dev/go-directory/aserto/directory/common/v2"
	dsw "github.com/aserto-dev/go-directory/aserto/directory/writer/v2"
	"github.com/pkg/errors"
)

type Loader struct {
	Objects   []*dsc.Object   `json:"objects"`
	Relations []*dsc.Relation `json:"relations"`
}

func (c *Client) Import(ctx context.Context, files []string) error {
	var data []Loader

	ctr := counter.New()

	// read all files
	for _, file := range files {
		var loader Loader
		c.UI.Normal().Msgf("Reading file %s", file)
		b, err := os.ReadFile(file)
		if err != nil {
			return errors.Wrapf(err, "failed to read file: [%s]", file)
		}
		if err := json.Unmarshal(b, &loader); err != nil {
			return errors.Wrapf(err, "failed unmarshal file: [%s]", file)
		}

		data = append(data, loader)
	}

	// import all objects
	fmt.Fprint(c.UI.Output(), "Importing objects...\n")
	for _, d := range data {
		for _, object := range d.Objects {
			_, err := c.Writer.SetObject(ctx, &dsw.SetObjectRequest{Object: object})
			if err != nil {
				return err
			}
			ctr.Objects.Incr().Print(c.UI.Output())
		}
		fmt.Fprintln(c.UI.Output())
	}

	// import all relations
	fmt.Fprint(c.UI.Output(), "Importing relations...\n")
	for _, d := range data {
		for _, relation := range d.Relations {
			_, err := c.Writer.SetRelation(ctx, &dsw.SetRelationRequest{Relation: relation})
			if err != nil {
				return err
			}
			ctr.Relations.Incr().Print(c.UI.Output())
		}
		fmt.Fprintln(c.UI.Output())
	}

	ctr.Print(c.UI.Output())

	return nil
}
