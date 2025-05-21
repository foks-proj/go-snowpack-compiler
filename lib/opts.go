package lib

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

type Language int

const (
	LangNone       Language = iota
	LangGeneric    Language = iota
	LangGo         Language = iota
	LangTypeScript Language = iota
)

type Options struct {
	lang Language

	langRaw string
	infile  string
	outfile string
	indir   string
	outdir  string
	pkg     string
	ext     string

	verbose bool
}

func (l Language) OutExt() string {
	switch l {
	case LangGo:
		return "go"
	case LangTypeScript:
		return "ts"
	default:
		return ""
	}
}

var ErrNotImplemented = errors.New("not implemented")

func ParseOpts() (*Options, error) {
	var opts Options
	err := makeCommand(&opts).Execute()
	if err != nil {
		return nil, err
	}
	return &opts, nil
}

func (o *Options) Run() error {
	runner := NewRunner(o)
	return runner.Run()
}

func isDir(d string) bool {
	st, err := os.Stat(d)
	if err != nil {
		return false
	}
	return st.IsDir()
}

func (o *Options) check() error {

	switch o.langRaw {
	case "go":
		o.lang = LangGo
	case "ts":
		o.lang = LangTypeScript
	default:
		return fmt.Errorf("unsupported language: %s", o.langRaw)
	}

	if (o.indir != "" && o.outdir == "") || (o.indir == "" && o.outdir != "") {
		return errors.New("must specify output directory with input directory")
	}
	if (o.indir != "" || o.outdir != "") && (o.infile != "" || o.outfile != "") {
		return errors.New("cannot use input or output file with input directory")
	}

	if o.indir != "" && !isDir(o.indir) {
		return fmt.Errorf("input directory %s does not exist", o.indir)
	}
	if o.outdir != "" && !isDir(o.outdir) {
		return fmt.Errorf("output directory %s does not exist", o.outdir)
	}

	if o.indir == "" && o.infile == "" {
		o.infile = "-"
	}
	if o.outdir == "" && o.outfile == "" {
		o.outfile = "-"
	}

	if o.pkg == "" {
		return errors.New("must specify package name")
	}

	if o.ext == "" {
		o.ext = "snowp"
	}

	return nil
}

func makeCommand(opts *Options) *cobra.Command {
	ret := &cobra.Command{
		Use:   "snowpc",
		Short: "Snowpack RPC compiler compile .snowp files",
		RunE: func(cmd *cobra.Command, args []string) error {
			return opts.check()
		},
	}
	ret.Flags().StringVarP(&opts.langRaw, "lang", "l", "go", "output language")
	ret.Flags().StringVarP(&opts.infile, "infile", "i", "", "input file")
	ret.Flags().StringVarP(&opts.outfile, "outfile", "o", "", "output file")
	ret.Flags().StringVarP(&opts.indir, "input-dir", "I", "", "input directory")
	ret.Flags().StringVarP(&opts.outdir, "output-dir", "O", "", "output directory")
	ret.Flags().StringVarP(&opts.pkg, "package", "p", "", "package name")
	ret.Flags().StringVarP(&opts.ext, "ext", "e", ".snowp", "file extension")
	ret.Flags().BoolVarP(&opts.verbose, "verbose", "v", false, "verbose output")
	return ret
}
