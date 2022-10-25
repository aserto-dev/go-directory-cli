package client

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"io"
	"os"
	"path"

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
		c.UI.Normal().Msgf("Deleting temporary dir [%s]", tmpDir)
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
		c.UI.Normal().Msgf("Closing tar file [%s]", file)
		tf.Close()
	}()

	gw, err := gzip.NewWriterLevel(tf, gzip.BestCompression)
	if err != nil {
		return nil
	}
	defer func() {
		c.UI.Normal().Msg("Closing gzip writer")
		gw.Close()
	}()

	tw := tar.NewWriter(gw)
	defer func() {
		c.UI.Normal().Msg("Closing tar writer")
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
	objTypes, _ := js.NewArrayWriter(path.Join(dirPath, "object_types.json"))
	defer objTypes.Close()

	permissions, _ := js.NewArrayWriter(path.Join(dirPath, "permissions.json"))
	defer permissions.Close()

	relTypes, _ := js.NewArrayWriter(path.Join(dirPath, "relation_types.json"))
	defer relTypes.Close()

	objects, _ := js.NewArrayWriter(path.Join(dirPath, "objects.json"))
	defer objects.Close()

	relations, _ := js.NewArrayWriter(path.Join(dirPath, "relations.json"))
	defer relations.Close()

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

		case *dse.ExportResponse_Permission:
			err = permissions.Write(m.Permission)

		case *dse.ExportResponse_RelationType:
			err = relTypes.Write(m.RelationType)

		case *dse.ExportResponse_Object:
			err = objects.Write(m.Object)

		case *dse.ExportResponse_Relation:
			err = relations.Write(m.Relation)

		default:
			c.UI.Exclamation().Msg("Unknown message type")
		}

		if err != nil {
			c.UI.Problem().Msgf("Error: %v", err)
		}
	}

	return nil
}
