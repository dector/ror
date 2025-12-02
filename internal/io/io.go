package io

type IO interface {
	Files() FileSystem
	Shell() Shell
	Std() Standard
}

type RealIO struct {
	files *RealFiles
	shell *RealShell
	std   *RealStandard
}

func NewIO() IO {
	return &RealIO{
		files: &RealFiles{},
		shell: &RealShell{},
		std:   &RealStandard{},
	}
}

func (r *RealIO) Files() FileSystem { return r.files }
func (r *RealIO) Shell() Shell      { return r.shell }
func (r *RealIO) Std() Standard     { return r.std }
