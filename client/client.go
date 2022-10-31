package client

import (
	"github.com/aserto-dev/clui"
	asertoClient "github.com/aserto-dev/go-aserto/client"
	dse "github.com/aserto-dev/go-directory/aserto/directory/exporter/v2"
	dsi "github.com/aserto-dev/go-directory/aserto/directory/importer/v2"
	dsr "github.com/aserto-dev/go-directory/aserto/directory/reader/v2"
	dsw "github.com/aserto-dev/go-directory/aserto/directory/writer/v2"
)

type Client struct {
	conn     *asertoClient.Connection
	Writer   dsw.WriterClient
	Exporter dse.ExporterClient
	Importer dsi.ImporterClient
	Reader   dsr.ReaderClient
	UI       *clui.UI
}

func New(conn *asertoClient.Connection, ui *clui.UI) (*Client, error) {
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
