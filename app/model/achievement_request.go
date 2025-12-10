package model

type AchievementCreateRequest struct {
    Title           string                 `json:"title"`
    Description     string                 `json:"description"`
    AchievementType string                 `json:"achievementType"`
    Details         map[string]interface{} `json:"details"`
    Tags            []string               `json:"tags"`
}

type AchievementUpdateInput struct {
	Title           string                 `json:"title"`
	Description     string                 `json:"description"`
	AchievementType string                 `json:"achievementType"`
	Details         map[string]interface{} `json:"details"`   // event, tahun, etc
	Tags            []string               `json:"tags"`
}
