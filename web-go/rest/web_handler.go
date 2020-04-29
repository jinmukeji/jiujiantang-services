package rest

type webHandler struct {
	res *Resource
}

func newWebHandler(configFile string) *webHandler {
	res, _ := LoadResourceFile(configFile)
	return &webHandler{
		res: res,
	}
}
