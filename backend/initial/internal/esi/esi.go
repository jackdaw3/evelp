package esi

type RemoteDataInit interface {
	Refresh() error
}
