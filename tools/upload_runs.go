package tools

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
)

type UploadRunsCmd struct {
	flags *flag.FlagSet
	// URL to upload to
	Url string
	// Directory to search for .run files
	RunsDir string
	// Path to either .run file or .tar.gz file with runs in it
	RunsFile string
	// Number of workers
	Workers int
	// Destination URL
	destUrl *url.URL
	// Work pool
	wp *WorkerPool[uploadItem]
}

type uploadItem struct {
	Name string
	Data io.ReadCloser
}

func NewUploadRunsCmd() *UploadRunsCmd {
	cmd := new(UploadRunsCmd)
	fg := flag.NewFlagSet("upload-runs", flag.ExitOnError)
	fg.StringVar(&cmd.Url, "url", "", "URL to upload to")
	fg.StringVar(&cmd.RunsDir, "dir", "", "Search this directory recursively for .run files and upload them")
	fg.StringVar(&cmd.RunsFile, "file", "", "Either a .run file, or a .tar.gz file made with 'json-export'")
	cmd.Workers = 4
	cmd.flags = fg
	return cmd
}

func (cmd *UploadRunsCmd) Flags() *flag.FlagSet {
	return cmd.flags
}

func (cmd *UploadRunsCmd) Description() string {
	return `upload runs from a file, directory, or archive`
}

func (cmd *UploadRunsCmd) Run() error {
	// Check url
	if cmd.Url == "" {
		return fmt.Errorf("must provide -url")
	}
	// If both are either empty, or both are given, error
	if (cmd.RunsDir == "") == (cmd.RunsFile == "") {
		return fmt.Errorf("must provide either -dir or -file, not both")
	}

	dstUrl, err := url.Parse(cmd.Url)
	if err != nil {
		return err
	}
	cmd.destUrl = dstUrl

	cmd.wp = NewWorkerPool(cmd.Workers, cmd.postWorker)
	defer cmd.wp.Close()
	// Handle directory
	if cmd.RunsDir != "" {
		return cmd.uploadDir(cmd.RunsDir)
	} else if cmd.RunsFile != "" {
		if strings.HasSuffix(cmd.RunsFile, ".tar.gz") {
			return cmd.UploadTar(cmd.RunsFile)
		} else {
			fd, err := os.Open(cmd.RunsFile)
			if err != nil {
				return err
			}
			cmd.putFile(cmd.RunsFile, fd)
		}
	}
	return nil
}

func (cmd *UploadRunsCmd) UploadTar(tarPath string) error {
	srcFd, err := os.Open(tarPath)
	if err != nil {
		return err
	}
	defer srcFd.Close()
	srcGz, err := gzip.NewReader(srcFd)
	if err != nil {
		return err
	}
	defer srcGz.Close()
	tarRd := tar.NewReader(srcGz)
	for {
		hdr, err := tarRd.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		buf := new(bytes.Buffer)
		if _, err := io.Copy(buf, tarRd); err != nil {
			return err
		}
		cmd.putFile(hdr.Name, io.NopCloser(buf))
	}
	return nil
}

func (cmd *UploadRunsCmd) uploadDir(dirPath string) error {
	rootFs := os.DirFS(cmd.RunsDir)
	return fs.WalkDir(rootFs, ".", func(fpath string, d fs.DirEntry, err error) error {
		if d.Type().IsRegular() && path.Ext(fpath) == ".run" {
			fd, err := rootFs.Open(fpath)
			if err != nil {
				return err
			}
			cmd.putFile(fpath, fd)
		}
		return nil
	})
}

func (cmd *UploadRunsCmd) postWorker(item uploadItem) {
	const LOG_FMT = "%s %s\n"
	defer item.Data.Close()
	res, err := http.Post(cmd.destUrl.String(), "application/json", item.Data)
	if err != nil {
		fmt.Printf(LOG_FMT, item.Name, err.Error())
		return
	}
	defer res.Body.Close()
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf(LOG_FMT, item.Name, err.Error())
		return
	}
	fmt.Printf(LOG_FMT, item.Name, string(resBody))
}

func (cmd *UploadRunsCmd) putFile(name string, src io.ReadCloser) {
	cmd.wp.Submit(uploadItem{Name: name, Data: src})
}
