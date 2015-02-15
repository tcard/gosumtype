package example

//go:generate gosumtype -test MyTree

type MyTree interface {
	Name() string
	isMyTree()
}

type Leaf string

func (l Leaf) Name() string { return string(l) }
func (l Leaf) isMyTree()    {}

type AnonBranch []MyTree

func (b AnonBranch) Name() string { return "AnonBranch" }
func (b AnonBranch) isMyTree()    {}
