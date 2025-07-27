package provider

type BaseProvider struct {
	ChatCompleteURL    string
	ChatCompleteMethod string
	ModelURL           string
	ModelMethod        string
	Headers            map[string]string
}
