package settings

type Setting struct {
	Key       string `json:"key"`
	Value     string `json:"value"`
	UpdatedAt string `json:"updated_at"`
}

type UpdateSettingsInput struct {
	Settings map[string]string `json:"settings"`
}
