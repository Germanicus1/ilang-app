package models

type Game struct {
    ID          string `json:"id"`
    Title       string `json:"title"`
    Description string `json:"description"`
    SubjectID   string `json:"subject_id"`
    Difficulty  int    `json:"difficulty_level"`
    CreatedAt   string `json:"created_at"`
}