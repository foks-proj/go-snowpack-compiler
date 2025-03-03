package main

import (
	"os"
	"path/filepath"
)

type FilePair struct {
	infile  string
	outfile string
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
			infile:  filepath.Join(o.indir, ent.Name()),
			outfile: filepath.Join(o.outdir, basename+"."+o.lang.OutExt()),
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
		infile:  opts.infile,
		outfile: opts.outfile,
	}
	f.files = append(f.files, fp)
	return nil
}

func (f *FilePair) run(o *Options) error {
	return ErrNotImplemented
}
