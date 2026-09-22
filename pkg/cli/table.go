package cli

import (
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

// NewTable 创建统一命令行表格
func NewTable() table.Writer {
	writer := table.NewWriter()
	writer.SetStyle(table.StyleLight)

	style := writer.Style()
	style.Format.Header = text.FormatDefault
	style.Color.Border = text.Colors{text.FgYellow}
	style.Color.Separator = text.Colors{text.FgYellow}
	style.Color.Header = text.Colors{text.Bold, text.FgHiWhite}

	return writer
}
