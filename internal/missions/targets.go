package missions

type Target struct {
	Name     string   `json:"name"`
	Country  string   `json:"country"`
	Notes    []string `json:"notes"`
	Complete bool     `json:"complete"`
}
