package wiretap

type Wiretap struct {
	API   *APIServer
	Proxy *SocksProxy
}

func NewWiretap() (*Wiretap, error) {
	var err error
	wiretap := &Wiretap{
		API: NewAPIServer(),
	}
	wiretap.Proxy, err = NewSocksProxy()
	return wiretap, err
}

func (w *Wiretap) ListenAndServe() error {
	errChan := make(chan error)

	go func() {
		errChan <- w.API.ListenAndServe()
	}()
	go func() {
		errChan <- w.Proxy.ListenAndServe()
	}()

	return <-errChan
}
