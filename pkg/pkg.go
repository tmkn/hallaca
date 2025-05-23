package pkg

type Pkg struct {
	Parent       *Pkg
	Name         string
	Version      string
	Dependencies []*Pkg
	IsLoop       bool
}
