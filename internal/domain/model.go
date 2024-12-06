package domain

type Status struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Tracker struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Project struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type User struct {
	ID    int    `json:"id"`
	Login string `json:"login"`
}

type Issue struct {
	ID         int     `json:"id"`
	Subject    string  `json:"subject"`
	Status     Status  `json:"status"`
	Tracker    Tracker `json:"tracker"`
	Project    Project `json:"project"`
	AssignedTo User    `json:"assigned_to"`
	Author     User    `json:"author"`
	Children   []Issue `json:"children,omitempty"`
}

type EvaluationResult struct {
	Message string
	Issue   *Issue
}
