package service

type Blockable interface {
	Block() error
	Unblock() error
}
