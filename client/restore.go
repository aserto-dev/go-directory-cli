package client

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"errors"
	"io"
	"os"
	"path"

	"github.com/aserto-dev/go-directory-cli/counter"
	"github.com/aserto-dev/go-directory-cli/js"
	dsc "github.com/aserto-dev/go-directory/aserto/directory/common/v2"
	dsi "github.com/aserto-dev/go-directory/aserto/directory/importer/v2"
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
	defer ctr.Print(c.UI.Output())

	g, iCtx := errgroup.WithContext(context.Background())
	stream, err := c.Importer.Import(iCtx)
	if err != nil {
		return err
	}

	g.Go(c.receiver(stream))

	g.Go(c.restoreHandler(stream, tr, ctr))

	return g.Wait()
}

func (c *Client) receiver(stream dsi.Importer_ImportClient) func() error {
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

func (c *Client) restoreHandler(stream dsi.Importer_ImportClient, tr *tar.Reader, ctr *counter.Counter) func() error {
	objectTypesCounter := ctr.ObjectTypes()
	permissionsCounter := ctr.Permissions()
	relationTypesCounter := ctr.RelationTypes()
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

			r, err := js.NewReader(tr, c.UI)
			if err != nil {
				return err
			}

			name := path.Clean(header.Name)
			switch name {
			case ObjectTypesFileName:
				if err := c.loadObjectTypes(stream, r, objectTypesCounter); err != nil {
					return err
				}

			case PermissionsFileName:
				if err := c.loadPermissions(stream, r, permissionsCounter); err != nil {
					return err
				}

			case RelationTypesFileName:
				if err := c.loadRelationTypes(stream, r, relationTypesCounter); err != nil {
					return err
				}

			case ObjectsFileName:
				if err := c.loadObjects(stream, r, objectsCounter); err != nil {
					return err
				}

			case RelationsFileName:
				if err := c.loadRelations(stream, r, relationsCounter); err != nil {
					return err
				}

			default:
				break
			}
		}

		if err := stream.CloseSend(); err != nil {
			return err
		}

		return nil
	}
}

func (c *Client) loadObjectTypes(stream dsi.Importer_ImportClient, objTypes *js.Reader, ctr *counter.Item) error {
	defer objTypes.Close()

	var m dsc.ObjectType

	for {
		err := objTypes.Read(&m)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if err := stream.Send(&dsi.ImportRequest{
			OpCode: dsi.Opcode_OPCODE_SET,
			Msg: &dsi.ImportRequest_ObjectType{
				ObjectType: &m,
			},
		}); err != nil {
			return err
		}

		ctr.Incr().Print(c.UI.Output())
	}

	return nil
}

func (c *Client) loadPermissions(stream dsi.Importer_ImportClient, permissions *js.Reader, ctr *counter.Item) error {
	defer permissions.Close()

	var m dsc.Permission

	for {
		err := permissions.Read(&m)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if err := stream.Send(&dsi.ImportRequest{
			OpCode: dsi.Opcode_OPCODE_SET,
			Msg: &dsi.ImportRequest_Permission{
				Permission: &m,
			},
		}); err != nil {
			return err
		}

		ctr.Incr().Print(c.UI.Output())
	}

	return nil
}

func (c *Client) loadRelationTypes(stream dsi.Importer_ImportClient, relTypes *js.Reader, ctr *counter.Item) error {
	defer relTypes.Close()

	var m dsc.RelationType

	for {
		err := relTypes.Read(&m)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if err := stream.Send(&dsi.ImportRequest{
			OpCode: dsi.Opcode_OPCODE_SET,
			Msg: &dsi.ImportRequest_RelationType{
				RelationType: &m,
			},
		}); err != nil {
			return err
		}

		ctr.Incr().Print(c.UI.Output())
	}

	return nil
}

func (c *Client) loadObjects(stream dsi.Importer_ImportClient, objects *js.Reader, ctr *counter.Item) error {
	defer objects.Close()

	var m dsc.Object

	for {
		err := objects.Read(&m)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if err := stream.Send(&dsi.ImportRequest{
			OpCode: dsi.Opcode_OPCODE_SET,
			Msg: &dsi.ImportRequest_Object{
				Object: &m,
			},
		}); err != nil {
			return err
		}
		ctr.Incr().Print(c.UI.Output())
	}

	return nil
}

func (c *Client) loadRelations(stream dsi.Importer_ImportClient, relations *js.Reader, ctr *counter.Item) error {
	defer relations.Close()

	var m dsc.Relation

	for {
		err := relations.Read(&m)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if err := stream.Send(&dsi.ImportRequest{
			OpCode: dsi.Opcode_OPCODE_SET,
			Msg: &dsi.ImportRequest_Relation{
				Relation: &m,
			},
		}); err != nil {
			return err
		}

		ctr.Incr().Print(c.UI.Output())
	}

	return nil
}
