package http

type Task struct {
	Name      string   `json:"name"`
	Cron      string   `json:"cron"`
	Uri       string   `json:"uri"`
	Method    string   `json:"method"`
	Locations []string `json:"locations"`
}

type TaskIds struct {
	TaskIds []string `json:"task_ids"`
}
