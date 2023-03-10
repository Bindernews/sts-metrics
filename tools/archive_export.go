package tools

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"os"
	"time"

	"github.com/bindernews/sts-msr/orm"
	"github.com/jackc/pgx/v4/pgxpool"
)

type ArchiveExportCmd struct {
	// CLI flag set
	flags *flag.FlagSet
	// Path to archive file to write into
	OutFile string
}

func NewArchiveExportCmd() *ArchiveExportCmd {
	cmd := new(ArchiveExportCmd)
	fg := flag.NewFlagSet("json-export", flag.ExitOnError)
	defaultOut := fmt.Sprintf("json-archive_%s.tar.gz", time.Now().UTC().Format(time.RFC3339))
	fg.StringVar(&cmd.OutFile, "out", defaultOut, "Output file name")
	cmd.flags = fg
	return cmd
}

func (cmd *ArchiveExportCmd) Flags() *flag.FlagSet {
	return cmd.flags
}

func (cmd *ArchiveExportCmd) Run() error {
	ctx := context.Background()
	pool, err := pgxpool.Connect(ctx, os.Getenv(EnvPostgresConn))
	if err != nil {
		return err
	}
	// Try to open the file we'll use for writing
	tarFile, err := os.Create(cmd.OutFile)
	if err != nil {
		return err
	}
	defer tarFile.Close()

	db := orm.New(pool)
	// Pick random status value, low-enough chance of collision to be fine
	// as long as we're not running 1000s of copies of this at the same time.
	myStatus := int16((rand.Uint32() % 32760) + 1)
	// Get rows we're going to update and set them as in-progress
	rowData, err := db.ArchiveBegin(ctx, myStatus)
	if err != nil {
		return err
	}
	// Uncompress data, then tar and compress
	if err := cmd.recompress(rowData, tarFile); err != nil {
		return err
	}
	// Mark as completed
	_, err = db.ArchiveComplete(ctx, myStatus)
	if err != nil {
		return err
	}
	return nil
}

func (cmd *ArchiveExportCmd) recompress(rowData []orm.Rawjsonarchive, dst io.Writer) error {
	tarGz, err := gzip.NewWriterLevel(dst, gzip.BestCompression)
	if err != nil {
		return err
	}
	defer tarGz.Close()
	tarWr := tar.NewWriter(tarGz)
	for _, arc := range rowData {
		data := arc.Bdata.Bytes
		hdr := tar.Header{
			Typeflag: tar.TypeReg,
			Name:     arc.PlayID + ".run",
			Size:     int64(len(data)),
			Mode:     0660, // octal!
		}
		if err := tarWr.WriteHeader(&hdr); err != nil {
			return err
		}
		if _, err := tarWr.Write(data); err != nil {
			return err
		}
	}
	if err := tarWr.Close(); err != nil {
		return err
	}
	return nil
}
