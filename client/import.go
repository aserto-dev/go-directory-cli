package client

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/aserto-dev/go-directory-cli/counter"
	"github.com/aserto-dev/go-directory-cli/js"
	"github.com/pkg/errors"
)

func (c *Client) Import(ctx context.Context, files []string) error {

	ctr := counter.New()

	// read all files
	for _, file := range files {
		c.UI.Normal().Msgf("Reading file %s", file)

		var objectType string
		reader, err := js.NewReader(file)
		if err != nil {
			c.UI.Problem().Msgf("Skipping file [%s]: [%s]", file, err.Error())
		}

		if reader != nil {
			objectType = reader.GetObjectType()
		} else {
			basePath := filepath.Base(file)
			switch basePath {
			case ObjectsFileName:
				objectType = ObjectsStr

			case RelationsFileName:
				objectType = RelationsStr
			default:
				return errors.New("files objects.json|relations.json or json root key not found")
			}
			b, err := os.Open(file)
			if err != nil {
				return errors.Wrapf(err, "failed to open file: [%s]", file)
			}
			reader, err = js.NewArrayReader(b)
			if err != nil {
				c.UI.Problem().Msgf("Skipping file [%s]: [%s]", file, err.Error())
			}
		}

		switch objectType {
		case ObjectsStr:
			if err := c.loadObjects(ctx, reader, ctr.Objects); err != nil {
				return err
			}

		case RelationsStr:
			if err := c.loadRelations(ctx, reader, ctr.Relations); err != nil {
				return err
			}
		default:
			return errors.Errorf("invalid object type: [%s]", objectType)
		}
		fmt.Fprintln(c.UI.Output())
	}

	ctr.Print(c.UI.Output())

	return nil
}
