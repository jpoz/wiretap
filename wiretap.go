package wiretap

type Wiretap struct {
	API   *APIServer
	Proxy *SocksProxy
}

func NewWiretap() *Wiretap {
	return &Wiretap{
		API:   NewAPIServer(),
		Proxy: &SocksProxy{},
	}
}

func (w *Wiretap) ListenAndServe() error {
	return w.API.ListenAndServe()
}
