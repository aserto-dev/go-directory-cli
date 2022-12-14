package client

import (
	"github.com/aserto-dev/clui"
	dse "github.com/aserto-dev/go-directory/aserto/directory/exporter/v2"
	dsi "github.com/aserto-dev/go-directory/aserto/directory/importer/v2"
	dsr "github.com/aserto-dev/go-directory/aserto/directory/reader/v2"
	dsw "github.com/aserto-dev/go-directory/aserto/directory/writer/v2"
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
	conn     grpc.ClientConnInterface
	Writer   dsw.WriterClient
	Exporter dse.ExporterClient
	Importer dsi.ImporterClient
	Reader   dsr.ReaderClient
	UI       *clui.UI
}

func New(conn grpc.ClientConnInterface, ui *clui.UI) (*Client, error) {
	c := Client{
		conn:     conn,
		Writer:   dsw.NewWriterClient(conn),
		Exporter: dse.NewExporterClient(conn),
		Importer: dsi.NewImporterClient(conn),
		Reader:   dsr.NewReaderClient(conn),
		UI:       ui,
	}
	return &c, nil
}
