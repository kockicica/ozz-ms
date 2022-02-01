package media_index

import (
	"fmt"

	"github.com/gosuri/uitable"
)

type AudioFile struct {
	ID       string
	Path     string
	Folder   string
	Root     string
	Name     string
	Artist   string
	Album    string
	Duration string
	Tags     []string
}

type AudioFiles []AudioFile

func (f AudioFiles) WriteOut() {
	table := uitable.New()
	table.MaxColWidth = 120
	table.Wrap = true
	for _, item := range f {
		table.AddRow("ID:", item.ID)
		table.AddRow("Name:", item.Name)
		table.AddRow("Path:", item.Path)
		table.AddRow("Folder:", item.Folder)
		table.AddRow("Artist:", item.Artist)
		table.AddRow("Album:", item.Album)
		table.AddRow("Duration:", item.Duration)
		table.AddRow("")
	}
	fmt.Println(table)

}
