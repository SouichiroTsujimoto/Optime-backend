package models

import (
	"time"

	"gorm.io/gorm"
)

// Event は予定調整のイベントを表す
type Event struct {
	ID          string    `gorm:"primaryKey;type:varchar(36)" json:"id"`
	Title       string    `gorm:"not null;type:varchar(255)" json:"title"`
	Description string    `gorm:"type:text" json:"description"`
	CreatorName string    `gorm:"type:varchar(100)" json:"creator_name"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// リレーション
	CandidateDates []CandidateDate `gorm:"foreignKey:EventID;constraint:OnDelete:CASCADE" json:"candidate_dates,omitempty"`
	Participants   []Participant   `gorm:"foreignKey:EventID;constraint:OnDelete:CASCADE" json:"participants,omitempty"`
}

// CandidateDate はイベントの候補日を表す
type CandidateDate struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	EventID   string         `gorm:"not null;type:varchar(36);index" json:"event_id"`
	DateTime  time.Time      `gorm:"not null" json:"date_time"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// リレーション
	Responses []Response `gorm:"foreignKey:CandidateDateID;constraint:OnDelete:CASCADE" json:"responses,omitempty"`
}

// Participant は参加者を表す
type Participant struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	EventID   string    `gorm:"not null;type:varchar(36);index" json:"event_id"`
	Name      string    `gorm:"not null;type:varchar(100)" json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// リレーション
	Responses []Response `gorm:"foreignKey:ParticipantID;constraint:OnDelete:CASCADE" json:"responses,omitempty"`
}

// Response は参加者の各候補日に対する回答を表す
type Response struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	ParticipantID   uint      `gorm:"not null;index" json:"participant_id"`
	CandidateDateID uint      `gorm:"not null;index" json:"candidate_date_id"`
	Status          string    `gorm:"not null;type:varchar(20)" json:"status"` // "available", "maybe", "unavailable"
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
