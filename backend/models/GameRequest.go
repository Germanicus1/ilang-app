package models

type GameRequest struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	SubjectID   string `json:"subject_id,omitempty"`
	Difficulty  int    `json:"difficulty_level,omitempty"`
}