package lib

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type File struct {
	name string
}

func (f *File) isStdPipe() bool {
	return f.name == "" || f.name == "-"
}

func (f *File) Filename(def string) string {
	if f.isStdPipe() {
		return def
	}
	return f.name
}

type Infile struct {
	File
}

func (i *Infile) Read() ([]byte, error) {

	var rc io.ReadCloser
	if i.isStdPipe() {
		rc = os.Stdin
	} else {
		var err error
		rc, err = os.Open(i.name)
		if err != nil {
			return nil, err
		}
	}
	defer rc.Close()
	return io.ReadAll(rc)
}

func (i *Infile) Name() string {
	if i.isStdPipe() {
		return "<stdin>"
	}
	return i.name
}

type Outfile struct {
	File
}

func newOutfile(f string) Outfile { return Outfile{File: File{name: f}} }
func newInfile(f string) Infile   { return Infile{File: File{name: f}} }

type FilePair struct {
	infile  Infile
	outfile Outfile
}

type FileSet struct {
	files []FilePair
}

func (f *FileSet) buildFromDir(o *Options) error {
	ents, err := os.ReadDir(o.indir)
	if err != nil {
		return err
	}
	for _, ent := range ents {
		if ent.IsDir() {
			continue
		}
		ext := filepath.Ext(ent.Name())
		if ext != "."+o.ext {
			continue
		}
		basename := ent.Name()[:len(ent.Name())-len(ext)]

		fp := FilePair{
			infile:  newInfile(filepath.Join(o.indir, ent.Name())),
			outfile: newOutfile(filepath.Join(o.outdir, basename+"."+o.lang.OutExt())),
		}
		f.files = append(f.files, fp)
	}
	return nil
}

func (f *FileSet) Build(opts *Options) error {
	if opts.indir != "" {
		return f.buildFromDir(opts)
	}
	fp := FilePair{
		infile:  newInfile(opts.infile),
		outfile: newOutfile(opts.outfile),
	}
	f.files = append(f.files, fp)
	return nil
}

func (f *FilePair) run(o *Options) error {
	md := NewMetadata(f, o)
	return md.run()
}

type Metadata struct {
	infile  Infile
	outfile Outfile
	lang    Language
	pkg     string
}

func NewMetadata(fp *FilePair, o *Options) Metadata {
	return Metadata{
		infile:  fp.infile,
		outfile: fp.outfile,
		lang:    o.lang,
		pkg:     o.pkg,
	}
}

func (m *Metadata) run() error {
	indat, err := m.infile.Read()
	if err != nil {
		return err
	}
	lexer := Lex(indat, m.infile.Name())
	for tok := range lexer.tokens {
		fmt.Printf("tok: %v\n", tok)
	}
	return ErrNotImplemented
}
