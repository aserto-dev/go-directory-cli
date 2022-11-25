package client

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"io"
	"os"
	"path"

	"github.com/aserto-dev/go-directory-cli/counter"
	"github.com/aserto-dev/go-directory-cli/js"
	dse "github.com/aserto-dev/go-directory/aserto/directory/exporter/v2"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (c *Client) Backup(ctx context.Context, file string) error {

	stream, err := c.Exporter.Export(ctx, &dse.ExportRequest{
		Options:   uint32(dse.Option_OPTION_ALL),
		StartFrom: &timestamppb.Timestamp{},
	})
	if err != nil {
		return err
	}

	tmpDir, err := os.MkdirTemp("", "*")
	if err != nil {
		return err
	}
	defer func() {
		_ = os.RemoveAll(tmpDir)
	}()

	dirPath := path.Join(tmpDir, "backup")
	if err := os.MkdirAll(dirPath, 0700); err != nil {
		return err
	}

	if err := c.createBackupFiles(stream, dirPath); err != nil {
		return err
	}

	tf, err := os.Create(file)
	if err != nil {
		return nil
	}
	defer func() {
		tf.Close()
	}()

	gw, err := gzip.NewWriterLevel(tf, gzip.BestCompression)
	if err != nil {
		return nil
	}
	defer func() {
		gw.Close()
	}()

	tw := tar.NewWriter(gw)
	defer func() {
		tw.Close()
	}()

	_ = addToArchive(tw, path.Join(dirPath, "object_types.json"))
	_ = addToArchive(tw, path.Join(dirPath, "permissions.json"))
	_ = addToArchive(tw, path.Join(dirPath, "relation_types.json"))
	_ = addToArchive(tw, path.Join(dirPath, "objects.json"))
	_ = addToArchive(tw, path.Join(dirPath, "relations.json"))

	return nil
}

func addToArchive(tw *tar.Writer, filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return err
	}

	header, err := tar.FileInfoHeader(info, info.Name())
	if err != nil {
		return err
	}

	if err := tw.WriteHeader(header); err != nil {
		return err
	}

	if _, err := io.Copy(tw, file); err != nil {
		return err
	}

	return nil
}

func (c *Client) createBackupFiles(stream dse.Exporter_ExportClient, dirPath string) error {
	objTypes, err := js.NewWriter(path.Join(dirPath, ObjectTypesFileName), ObjectTypesStr)
	if err != nil {
		return err
	}
	defer objTypes.Close()

	permissions, err := js.NewWriter(path.Join(dirPath, PermissionsFileName), PermissionsStr)
	if err != nil {
		return err
	}
	defer permissions.Close()

	relTypes, err := js.NewWriter(path.Join(dirPath, RelationTypesFileName), RelationTypesStr)
	if err != nil {
		return err
	}
	defer relTypes.Close()

	objects, err := js.NewWriter(path.Join(dirPath, ObjectsFileName), ObjectsStr)
	if err != nil {
		return err
	}
	defer objects.Close()

	relations, err := js.NewWriter(path.Join(dirPath, RelationsFileName), RelationsStr)
	if err != nil {
		return err
	}
	defer relations.Close()

	ctr := counter.New()
	objectTypesCounter := ctr.ObjectTypes()
	permissionsCounter := ctr.Permissions()
	relationTypesCounter := ctr.RelationTypes()
	objectsCounter := ctr.Objects()
	relationsCounter := ctr.Relations()

	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		switch m := msg.Msg.(type) {
		case *dse.ExportResponse_ObjectType:
			err = objTypes.Write(m.ObjectType)
			objectTypesCounter.Incr().Print(c.UI.Output())

		case *dse.ExportResponse_Permission:
			err = permissions.Write(m.Permission)
			permissionsCounter.Incr().Print(c.UI.Output())

		case *dse.ExportResponse_RelationType:
			err = relTypes.Write(m.RelationType)
			relationTypesCounter.Incr().Print(c.UI.Output())

		case *dse.ExportResponse_Object:
			err = objects.Write(m.Object)
			objectsCounter.Incr().Print(c.UI.Output())

		case *dse.ExportResponse_Relation:
			err = relations.Write(m.Relation)
			relationsCounter.Incr().Print(c.UI.Output())

		default:
			c.UI.Exclamation().Msg("Unknown message type")
		}

		if err != nil {
			c.UI.Problem().Msgf("Error: %v", err)
		}
	}

	ctr.Print(c.UI.Output())

	return nil
}
