package media_index

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/analysis/analyzer/keyword"
	"github.com/blevesearch/bleve/v2/analysis/lang/en"
	"github.com/blevesearch/bleve/v2/mapping"
	"github.com/blevesearch/bleve/v2/search/query"
	"github.com/blevesearch/bleve_index_api"
)

var BatchSize = 1000
var DefaultIndexName = "audio.bl"

func createIndexMapping() (mapping.IndexMapping, error) {

	engTextFieldMapping := bleve.NewTextFieldMapping()
	engTextFieldMapping.Analyzer = en.AnalyzerName

	keywordFieldMapping := bleve.NewTextFieldMapping()
	keywordFieldMapping.Analyzer = keyword.Name

	audioMapping := bleve.NewDocumentMapping()
	audioMapping.AddFieldMappingsAt("path", keywordFieldMapping)
	audioMapping.AddFieldMappingsAt("folder", keywordFieldMapping)
	audioMapping.AddFieldMappingsAt("name", engTextFieldMapping)
	audioMapping.AddFieldMappingsAt("artist", engTextFieldMapping)
	audioMapping.AddFieldMappingsAt("album", engTextFieldMapping)
	audioMapping.AddFieldMappingsAt("tags", keywordFieldMapping)
	audioMapping.AddFieldMappingsAt("duration", keywordFieldMapping)
	audioMapping.AddFieldMappingsAt("root", keywordFieldMapping)
	indexMapping := bleve.NewIndexMapping()
	indexMapping.AddDocumentMapping("audio", audioMapping)
	indexMapping.DefaultAnalyzer = "en"

	return indexMapping, nil

}

type MediaIndex struct {
	indexName    string
	index        bleve.Index
	batch        *bleve.Batch
	batchCount   int
	totalCount   int
	BatchWritten chan int
	Verbose      bool
}

func (i *MediaIndex) Create() error {
	indexMapping, err := createIndexMapping()
	if err != nil {
		return err
	}
	index, err := bleve.New(i.indexName, indexMapping)
	if err != nil {
		return err
	}
	i.index = index
	i.batch = i.index.NewBatch()
	return nil
}

func (i *MediaIndex) Open() error {
	idx, err := bleve.Open(i.indexName)
	if err != nil {
		absPath, err1 := filepath.Abs(i.indexName)
		if err1 != nil {
			return err
		}
		return errors.New(fmt.Sprintf("cannot open media index, path %s does not exist", absPath))
	}
	i.index = idx
	i.batch = i.index.NewBatch()
	return nil
}

func (i *MediaIndex) AddItem(item AudioFile) error {
	err := i.batch.Index(item.ID, item)
	if err != nil {
		return err
	}
	i.batchCount++
	i.totalCount++
	if i.batchCount == BatchSize {
		err = i.index.Batch(i.batch)
		if err != nil {
			return err
		}
		i.sendBatchWritten(i.totalCount)
		i.batch = i.index.NewBatch()
		i.batchCount = 0
	}
	return nil
}

func (i *MediaIndex) Flush() error {
	if i.batchCount > 0 {
		err := i.index.Batch(i.batch)
		if err != nil {
			return err
		}
		i.batchCount = 0
		i.sendBatchWritten(i.totalCount)
	}
	return nil
}

func (i *MediaIndex) Close() error {
	return i.index.Close()
}

func (i *MediaIndex) sendBatchWritten(count int) {
	select {
	case i.BatchWritten <- count:
	default:

	}
}

func (i *MediaIndex) Query(term string) (AudioFiles, error) {
	qr := bleve.NewQueryStringQuery(term)
	searchReq := bleve.NewSearchRequest(qr)
	searchReq.Fields = []string{"ID", "Path", "Name", "Artist", "Duration", "Root"}
	searchReq.From = 0
	searchReq.Size = 10000
	//data, err := json.Marshal(searchReq)
	//if err != nil {
	//	return nil, err
	//}
	//fmt.Println(string(data))
	if i.Verbose {
		qs, err := query.DumpQuery(i.index.Mapping(), qr)
		if err != nil {
			return nil, err
		}
		fmt.Println(qs)
	}
	res, err := i.index.Search(searchReq)
	if err != nil {
		//log.Println(err)
		return nil, err
	}
	var ret []AudioFile

	for _, hit := range res.Hits {
		doc, err := i.index.Document(hit.ID)
		if err != nil {
			//log.Println(err)
			return nil, err
		}
		af := AudioFile{
			ID: hit.ID,
		}
		doc.VisitFields(func(field index.Field) {
			switch field.Name() {
			case "Path":
				af.Path = string(field.Value())
			case "Folder":
				af.Folder = string(field.Value())
			case "Name":
				af.Name = string(field.Value())
			case "Artist":
				af.Artist = string(field.Value())
			case "Album":
				af.Album = string(field.Value())
			case "Tags":
				af.Tags = []string{}
			case "Duration":
				//dur, err := time.ParseDuration(string(field.Value()))
				//if err == nil {
				//	af.Duration = dur
				//}
				af.Duration = string(field.Value())
			case "Root":
				af.Root = string(field.Value())
			}
		})
		ret = append(ret, af)
	}
	return ret, nil
}

func (i *MediaIndex) GetPath(id string) string {
	doc, err := i.index.Document(id)
	if err != nil {
		return ""
	}
	path := ""
	doc.VisitFields(func(field index.Field) {
		if field.Name() == "Path" {
			path = string(field.Value())
			return
		}
	})
	return path
}

func NewIndex(indexName string) *MediaIndex {
	return &MediaIndex{
		indexName:    indexName,
		batchCount:   0,
		totalCount:   0,
		BatchWritten: make(chan int),
	}
}
