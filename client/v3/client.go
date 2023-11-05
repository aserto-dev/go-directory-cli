package v3

import (
	"github.com/aserto-dev/clui"

	dse3 "github.com/aserto-dev/go-directory/aserto/directory/exporter/v3"
	dsi3 "github.com/aserto-dev/go-directory/aserto/directory/importer/v3"
	dsm3 "github.com/aserto-dev/go-directory/aserto/directory/model/v3"
	dsr3 "github.com/aserto-dev/go-directory/aserto/directory/reader/v3"
	dsw3 "github.com/aserto-dev/go-directory/aserto/directory/writer/v3"

	"google.golang.org/grpc"
)

type Client struct {
	conn     grpc.ClientConnInterface
	Model    dsm3.ModelClient
	Reader   dsr3.ReaderClient
	Writer   dsw3.WriterClient
	Importer dsi3.ImporterClient
	Exporter dse3.ExporterClient
	UI       *clui.UI
}

func New(conn grpc.ClientConnInterface, ui *clui.UI) (*Client, error) {
	c := Client{
		conn:     conn,
		Model:    dsm3.NewModelClient(conn),
		Reader:   dsr3.NewReaderClient(conn),
		Writer:   dsw3.NewWriterClient(conn),
		Importer: dsi3.NewImporterClient(conn),
		Exporter: dse3.NewExporterClient(conn),
		UI:       ui,
	}
	return &c, nil
}
