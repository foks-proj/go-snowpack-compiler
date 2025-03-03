package lib

type Runner struct {
	opts *Options
}

func NewRunner(o *Options) *Runner {
	return &Runner{opts: o}
}

func (r *Runner) Run() error {
	fs := &FileSet{}
	err := fs.Build(r.opts)
	if err != nil {
		return err
	}
	for _, fp := range fs.files {
		err = fp.run(r.opts)
		if err != nil {
			return err
		}
	}
	return nil
}
