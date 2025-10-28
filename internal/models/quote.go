package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Quote represents a generated quote stored in the database
type Quote struct {
	ID         uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	Tag        string    `gorm:"type:varchar(50);not null;index" json:"tag"`
	TagSource  string    `gorm:"type:varchar(20);not null" json:"tag_source"` // "preset" or "custom"
	QuoteText  string    `gorm:"type:text;not null" json:"quote_text"`
	Author     *string   `gorm:"type:varchar(255)" json:"author,omitempty"`
	Source     string    `gorm:"type:varchar(50);not null" json:"source"` // "openrouter"
	CreatedAt  time.Time `json:"created_at"`
	LatencyMs  int       `json:"latency_ms"`
	ClientIP   string    `gorm:"type:varchar(45)" json:"client_ip"`
	UserAgent  string    `gorm:"type:text" json:"user_agent"`
}

// BeforeCreate will set a UUID rather than numeric ID
func (q *Quote) BeforeCreate(tx *gorm.DB) error {
	if q.ID == uuid.Nil {
		q.ID = uuid.New()
	}
	return nil
}

// ValidTags is the list of preset emotion tags
var ValidTags = []string{
	"joy", "sadness", "anger", "fear", "surprise", "love", "gratitude", "resilience",
	"optimism", "melancholy", "confidence", "anxiety", "curiosity", "hope", "calm",
	"nostalgia", "wonder", "determination", "humor", "serenity", "loneliness", "pride",
	"forgiveness", "humility", "ambition", "compassion", "playful", "boredom", "zeal", "contentment",
}

// IsValidTag checks if a tag is in the preset list
func IsValidTag(tag string) bool {
	for _, validTag := range ValidTags {
		if tag == validTag {
			return true
		}
	}
	return false
}

// GetTagSource determines if a tag is preset or custom
func GetTagSource(tag string) string {
	if IsValidTag(tag) {
		return "preset"
	}
	return "custom"
}
