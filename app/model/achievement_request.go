package model

type AchievementCreateRequest struct {
    Title           string                 `json:"title"`
    Description     string                 `json:"description"`
    AchievementType string                 `json:"achievementType"`
    Details         map[string]interface{} `json:"details"`
    Tags            []string               `json:"tags"`
}
