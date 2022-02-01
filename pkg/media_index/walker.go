package media_index

import (
	"crypto/md5"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/barasher/go-exiftool"
	"github.com/h2non/filetype"
)

type AudioWalker struct {
	File     chan AudioFile
	Progress <-chan int
	Finished chan bool
	exif     *exiftool.Exiftool
}

func NewAudioWalker(paths []string) (*AudioWalker, error) {
	w := AudioWalker{
		File:     make(chan AudioFile, 1),
		Progress: make(chan int, 1),
		Finished: make(chan bool, 1),
	}
	exif, err := exiftool.NewExiftool()
	if err != nil {
		return nil, err
	}
	w.exif = exif
	go func() {
		defer close(w.File)
		for _, fpath := range paths {
			err := filepath.WalkDir(fpath, walk(w, fpath))
			if err != nil {
				log.Println("Error walking for:", fpath)
			}

		}
		w.exif.Close()
		w.Finished <- true
	}()

	return &w, nil
}

func walk(w AudioWalker, fpath string) func(lpath string, d fs.DirEntry, err error) error {
	return func(lpath string, d fs.DirEntry, err error) error {
		if !d.IsDir() {
			file, err := os.Open(lpath)
			if err != nil {
				log.Println("Skip indexing, error reading:", lpath)
				return nil
			}
			header := make([]byte, 261)
			n, err := file.Read(header)
			if err != nil || n != 261 {
				log.Println("Skip indexing, error reading header for:", lpath)
				return nil
			}
			file.Close()
			if filetype.IsAudio(header) {

				hs, err := createHash(lpath)
				if err != nil {
					return nil
				}

				af := AudioFile{
					ID:   hs,
					Path: lpath,
					Name: d.Name(),
					Root: fpath,
				}

				fileInfos := w.exif.ExtractMetadata(lpath)
				for _, fileInfo := range fileInfos {
					if fileInfo.Err != nil {
						continue
					}
					album, err := fileInfo.GetString("Album")
					if err == nil {
						af.Album = album
					}
					artist, err := fileInfo.GetString("Artist")
					if err == nil {
						af.Artist = artist
					}
					length, err := fileInfo.GetString("Duration")
					if err == nil {
						af.Duration = makeNiceDuration(length)
					}
					//for k, v := range fileInfo.Fields {
					//	log.Printf("[%v]: %v\n", k, v)
					//}
				}

				ldir, _ := filepath.Split(lpath)
				rel, err := filepath.Rel(fpath, ldir)
				if err != nil {
					return nil
				}
				//_, folder := filepath.Split(rel)
				af.Folder = rel
				af.Tags = strings.Split(rel, string(filepath.Separator))

				w.File <- af

			}
			return nil
		}
		return nil
	}
}

func createHash(filename string) (string, error) {
	//f := strings.NewReader(filename)
	//hs := sha256.New()
	//if _, err := io.Copy(hs, f); err != nil {
	//	return "", err
	//}
	hs := md5.New()
	hs.Write([]byte(filename))
	return fmt.Sprintf("%x", hs.Sum(nil)), nil
}

func parseDuration(d string) time.Duration {
	dd := strings.ReplaceAll(d, "(approx)", "")
	dd = strings.TrimSpace(dd)
	dd = strings.Replace(dd, ":", "h", 1)
	dd = strings.Replace(dd, ":", "m", 1)
	dd = fmt.Sprintf("%ss", dd)

	ddd, err := time.ParseDuration(dd)
	if err != nil {
		return time.Second * 0
	}
	return ddd
}

func makeNiceDuration(d string) string {
	dd := strings.ReplaceAll(d, "(approx)", "")
	dd = strings.TrimSpace(dd)
	return dd
}
