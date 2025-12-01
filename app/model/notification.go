package model

type Notification struct {
    Title   string `json:"title"`
    Message string `json:"message"`
    SentAt  int64  `json:"sent_at"`
}
