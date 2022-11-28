package client

import (
	"context"
	"fmt"
	"os"

	"github.com/aserto-dev/go-directory-cli/counter"
	"github.com/aserto-dev/go-directory-cli/js"
	"github.com/pkg/errors"
)

func (c *Client) Import(ctx context.Context, files []string) error {

	ctr := counter.New()

	// read all files
	for _, file := range files {
		c.UI.Normal().Msgf("Reading file %s", file)
		err := c.importFile(ctx, ctr, file)
		if err != nil {
			return err
		}
	}

	ctr.Print(c.UI.Output())

	return nil
}

func (c *Client) importFile(ctx context.Context, ctr *counter.Counter, file string) error {
	r, err := os.Open(file)
	if err != nil {
		return errors.Wrapf(err, "failed to open file: [%s]", file)
	}
	defer r.Close()

	var objectType string
	reader, err := js.NewReader(r, c.UI)
	if err != nil {
		c.UI.Problem().Msgf("Skipping file [%s]: [%s]", file, err.Error())
		return nil
	}

	if reader != nil {
		objectType = reader.GetObjectType()
	} else {
		c.UI.Problem().Msgf("Skipping file [%s]: invalid json format", file)
		return nil
	}

	stream, err := c.Importer.Import(ctx)
	if err != nil {
		return err
	}

	switch objectType {
	case ObjectsStr:
		if err := c.loadObjects(stream, reader, ctr.Objects); err != nil {
			return err
		}

	case RelationsStr:
		if err := c.loadRelations(stream, reader, ctr.Relations); err != nil {
			return err
		}
	default:
		return errors.Errorf("invalid object type: [%s]", objectType)
	}
	fmt.Fprintln(c.UI.Output())
	return nil
}
