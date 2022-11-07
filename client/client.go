package client

import (
	asertogoClient "github.com/aserto-dev/aserto-go/client"
	"github.com/aserto-dev/clui"
	dse "github.com/aserto-dev/go-directory/aserto/directory/exporter/v2"
	dsi "github.com/aserto-dev/go-directory/aserto/directory/importer/v2"
	dsr "github.com/aserto-dev/go-directory/aserto/directory/reader/v2"
	dsw "github.com/aserto-dev/go-directory/aserto/directory/writer/v2"
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
	conn     *asertogoClient.Connection
	Writer   dsw.WriterClient
	Exporter dse.ExporterClient
	Importer dsi.ImporterClient
	Reader   dsr.ReaderClient
	UI       *clui.UI
}

func New(conn *asertogoClient.Connection, ui *clui.UI) (*Client, error) {
	c := Client{
		conn:     conn,
		Writer:   dsw.NewWriterClient(conn.Conn),
		Exporter: dse.NewExporterClient(conn.Conn),
		Importer: dsi.NewImporterClient(conn.Conn),
		Reader:   dsr.NewReaderClient(conn.Conn),
		UI:       ui,
	}
	return &c, nil
}
