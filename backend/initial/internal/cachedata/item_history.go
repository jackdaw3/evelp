package cachedata

type itemHistroy struct {
}

func (i *itemHistroy) invoke() func() {
	return func() {}
}
