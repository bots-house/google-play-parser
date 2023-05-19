package models

type Permission struct {
	Type    string `json:"type"`
	Summary string `json:"summary,omitempty"`
}
