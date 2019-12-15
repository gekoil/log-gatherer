package internal

type RequestOptions struct {
	StdOut bool   `json:"stdOut,omitempty"`
	StdErr bool   `json:"stdErr,omitempty"`
	Filter Filter `json:"filter"`
}

type Filter struct {
	ContainerName  string `json:"name,omitempty"`
	ContainerLabel string `json:"label,omitempty"`
	ContainerId    string `json:"id,omitempty"`
}

type Containers map[string]chan int
