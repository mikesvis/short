package service

type service struct {
	storage storageURL
}

func New(s storageURL) service {
	return service{storage: s}
}
