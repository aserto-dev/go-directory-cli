package client

import (
	"github.com/aserto-dev/clui"

	dse2 "github.com/aserto-dev/go-directory/aserto/directory/exporter/v2"
	dsi2 "github.com/aserto-dev/go-directory/aserto/directory/importer/v2"
	dsr2 "github.com/aserto-dev/go-directory/aserto/directory/reader/v2"
	dsw2 "github.com/aserto-dev/go-directory/aserto/directory/writer/v2"

	dse3 "github.com/aserto-dev/go-directory/aserto/directory/exporter/v3"
	dsi3 "github.com/aserto-dev/go-directory/aserto/directory/importer/v3"
	dsm3 "github.com/aserto-dev/go-directory/aserto/directory/model/v3"
	dsr3 "github.com/aserto-dev/go-directory/aserto/directory/reader/v3"
	dsw3 "github.com/aserto-dev/go-directory/aserto/directory/writer/v3"

	"google.golang.org/grpc"
)

const (
	ObjectsStr            = "objects"
	ObjectsFileName       = "objects.json"
	RelationsStr          = "relations"
	RelationsFileName     = "relations.json"
	ObjectTypesStr        = "object_types"
	ObjectTypesFileName   = "object_types.json"
	PermissionsStr        = "permissions"
	PermissionsFileName   = "permissions.json"
	RelationTypesStr      = "relation_types"
	RelationTypesFileName = "relation_types.json"
)

type Client struct {
	conn      grpc.ClientConnInterface
	Writer    dsw2.WriterClient
	Exporter  dse2.ExporterClient
	Importer  dsi2.ImporterClient
	Reader    dsr2.ReaderClient
	Model3    dsm3.ModelClient
	Reader3   dsr3.ReaderClient
	Writer3   dsw3.WriterClient
	Importer3 dsi3.ImporterClient
	Exporter3 dse3.ExporterClient
	UI        *clui.UI
}

func New(conn grpc.ClientConnInterface, ui *clui.UI) (*Client, error) {
	c := Client{
		conn:      conn,
		Writer:    dsw2.NewWriterClient(conn),
		Exporter:  dse2.NewExporterClient(conn),
		Importer:  dsi2.NewImporterClient(conn),
		Reader:    dsr2.NewReaderClient(conn),
		Model3:    dsm3.NewModelClient(conn),
		Reader3:   dsr3.NewReaderClient(conn),
		Writer3:   dsw3.NewWriterClient(conn),
		Importer3: dsi3.NewImporterClient(conn),
		Exporter3: dse3.NewExporterClient(conn),
		UI:        ui,
	}
	return &c, nil
}
