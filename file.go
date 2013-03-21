// Copyright 2012 Ruben Pollan <meskio@sindominio.net>
// Use of this source code is governed by a LGPL licence
// version 3 or later that can be found in the LICENSE file.

package epubgo

import (
	"archive/zip"
	"encoding/xml"
	"errors"
	"io"
)

func openFile(file *zip.Reader, path string) (io.ReadCloser, error) {
	for _, f := range file.File {
		if f.Name == path {
			return f.Open()
		}
	}
	return nil, errors.New("File " + path + " not found")
}

type rootfile struct {
	Path string `xml:"full-path,attr"`
}
type container_xml struct {
	// FIXME: only support for one rootfile, can it be more than one?
	Rootfile rootfile `xml:"rootfiles>rootfile"`
}

func openOPF(file *zip.Reader) (io.ReadCloser, error) {
	f, err := openFile(file, "META-INF/container.xml")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var c container_xml
	decoder := xml.NewDecoder(f)
	err = decoder.Decode(&c)
	if err != nil {
		return nil, err
	}

	return openFile(file, c.Rootfile.Path)
}
