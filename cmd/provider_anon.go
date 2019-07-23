package cmd

func (p projectConfig) getPath() string {
	return p.Location
}

func (p projectConfig) getURL() string {
	return p.URL
}

func (p projectConfig) getToken() string {
	return p.Token
}
