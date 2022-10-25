package client

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

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
		for i, object := range d.Objects {
			_, err := c.Writer.SetObject(ctx, &dsw.SetObjectRequest{Object: object})
			if err != nil {
				return err
			}
			fmt.Fprintf(c.UI.Output(), "\033[2K\r%15s %d", "objects:", i+1)
		}
		fmt.Fprintln(c.UI.Output())
	}

	// import all relations
	fmt.Fprint(c.UI.Output(), "Importing relations...\n")
	for i, d := range data {
		for _, relation := range d.Relations {
			_, err := c.Writer.SetRelation(ctx, &dsw.SetRelationRequest{Relation: relation})
			if err != nil {
				return err
			}
			fmt.Fprintf(c.UI.Output(), "\033[2K\r%15s %d", "relations:", i+1)
		}
		fmt.Fprintln(c.UI.Output())
	}

	return nil
}
