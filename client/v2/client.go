package v2

import (
	"github.com/aserto-dev/clui"

	dse2 "github.com/aserto-dev/go-directory/aserto/directory/exporter/v2"
	dsi2 "github.com/aserto-dev/go-directory/aserto/directory/importer/v2"
	dsr2 "github.com/aserto-dev/go-directory/aserto/directory/reader/v2"
	dsw2 "github.com/aserto-dev/go-directory/aserto/directory/writer/v2"

	"google.golang.org/grpc"
)

type Client struct {
	conn     grpc.ClientConnInterface
	Writer   dsw2.WriterClient
	Exporter dse2.ExporterClient
	Importer dsi2.ImporterClient
	Reader   dsr2.ReaderClient
	UI       *clui.UI
}

func New(conn grpc.ClientConnInterface, ui *clui.UI) (*Client, error) {
	c := Client{
		conn:     conn,
		Writer:   dsw2.NewWriterClient(conn),
		Exporter: dse2.NewExporterClient(conn),
		Importer: dsi2.NewImporterClient(conn),
		Reader:   dsr2.NewReaderClient(conn),
		UI:       ui,
	}
	return &c, nil
}
