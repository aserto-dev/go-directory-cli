package v3

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"errors"
	"io"
	"os"
	"path"
	"strings"

	"github.com/aserto-dev/go-directory-cli/client/x"
	"github.com/aserto-dev/go-directory-cli/counter"
	"github.com/aserto-dev/go-directory-cli/js"
	dsc3 "github.com/aserto-dev/go-directory/aserto/directory/common/v3"
	dsi3 "github.com/aserto-dev/go-directory/aserto/directory/importer/v3"
	"golang.org/x/sync/errgroup"
)

func (c *Client) Restore(ctx context.Context, file string) error {
	tf, err := os.Open(file)
	if err != nil {
		return err
	}
	defer tf.Close()

	gz, err := gzip.NewReader(tf)
	if err != nil {
		return err
	}
	defer gz.Close()

	tr := tar.NewReader(gz)

	ctr := counter.New()
	defer ctr.Print(c.Out())

	g, iCtx := errgroup.WithContext(context.Background())
	stream, err := c.Importer.Import(iCtx)
	if err != nil {
		return err
	}

	g.Go(c.receiver(stream))

	g.Go(c.restoreHandler(stream, tr, ctr))

	return g.Wait()
}

func (c *Client) receiver(stream dsi3.Importer_ImportClient) func() error {
	return func() error {
		for {
			_, err := stream.Recv()
			if err == io.EOF {
				return nil
			}

			if err != nil {
				return err
			}
		}
	}
}

func (c *Client) restoreHandler(stream dsi3.Importer_ImportClient, tr *tar.Reader, ctr *counter.Counter) func() error {
	objectsCounter := ctr.Objects()
	relationsCounter := ctr.Relations()

	return func() error {
		for {
			header, err := tr.Next()
			if errors.Is(err, io.EOF) {
				break
			}

			if err != nil {
				return err
			}

			if header == nil || header.Typeflag != tar.TypeReg {
				continue
			}

			r, err := js.NewReader(tr)
			if err != nil {
				return err
			}

			name := path.Clean(header.Name)
			switch name {
			case x.ObjectsFileName:
				if err := c.loadObjects(stream, r, objectsCounter); err != nil {
					return err
				}

			case x.RelationsFileName:
				if err := c.loadRelations(stream, r, relationsCounter); err != nil {
					return err
				}
			}
		}

		if err := stream.CloseSend(); err != nil {
			return err
		}

		return nil
	}
}

func (c *Client) loadObjects(stream dsi3.Importer_ImportClient, objects *js.Reader, ctr *counter.Item) error {
	defer objects.Close()

	var m dsc3.Object

	for {
		err := objects.Read(&m)
		if err == io.EOF {
			break
		}

		if err != nil {
			if strings.Contains(err.Error(), "unknown field") {
				ctr.Skip()
				continue
			}
			return err
		}

		if err := stream.Send(&dsi3.ImportRequest{
			OpCode: dsi3.Opcode_OPCODE_SET,
			Msg: &dsi3.ImportRequest_Object{
				Object: &m,
			},
		}); err != nil {
			return err
		}
		ctr.Incr().Print(c.Out())
	}

	return nil
}

func (c *Client) loadRelations(stream dsi3.Importer_ImportClient, relations *js.Reader, ctr *counter.Item) error {
	defer relations.Close()

	var m dsc3.Relation

	for {
		err := relations.Read(&m)
		if err == io.EOF {
			break
		}
		if err != nil {
			if strings.Contains(err.Error(), "unknown field") {
				ctr.Skip()
				continue
			}
			return err
		}

		if err := stream.Send(&dsi3.ImportRequest{
			OpCode: dsi3.Opcode_OPCODE_SET,
			Msg: &dsi3.ImportRequest_Relation{
				Relation: &m,
			},
		}); err != nil {
			return err
		}

		ctr.Incr().Print(c.Out())
	}

	return nil
}
