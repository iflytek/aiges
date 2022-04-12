package xsf

type killerWrapper struct {
	callback func()
}

func (k *killerWrapper) Closeout() {
	k.callback()
}
