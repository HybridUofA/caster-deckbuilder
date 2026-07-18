package speedrobo

// PageConfig contains the endpoint configuration embedded in the upstream card page.

type PageConfig struct {
	AjaxURL string `json:"ajax_url"`
	Nonce   string `json:"nonce"`
}
