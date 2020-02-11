package tests

type apiTest struct {
	login         string
	password      string
	ip            string
	subnet        string
	responseError error
}

func (a *apiTest) loginIs(login string) error {
	a.login = login
	return nil
}

func (a *apiTest) passwordIs(pass string) error {
	a.password = pass
	return nil
}

func (a *apiTest) ipIs(ip string) error {
	a.ip = ip
	return nil
}

func (a *apiTest) subnetIs(subnet string) error {
	a.subnet = subnet
	return nil
}
