package tar

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"

	"github.com/charmbracelet/log"
	"github.com/tmkn/hallaca/pkg"
)

func DownloadTar(pkg *pkg.Pkg) {
	resp, err := http.Get(pkg.Metadata.Dist.Tarball)
	if err != nil {
		panic(fmt.Errorf("failed to download file: %v", err))
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		panic(fmt.Errorf("bad status: %s", resp.Status))
	}

	gzReader, err := gzip.NewReader(resp.Body)
	if err != nil {
		panic(fmt.Errorf("failed to create gzip reader: %v", err))
	}
	defer gzReader.Close()

	tarReader := tar.NewReader(gzReader)

	log.Infof("Files in archive for %s@%s:", pkg.Name, pkg.Version)
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(fmt.Errorf("error reading tar: %v", err))
		}

		log.Infof("%s", header.Name)
	}
}
