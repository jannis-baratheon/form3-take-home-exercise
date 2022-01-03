package form3apiclient

type context struct {
	Url string
}

func createContext(url string) context {
	return context{Url: url}
}