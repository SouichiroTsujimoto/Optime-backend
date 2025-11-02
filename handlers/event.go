package handlers

import (
	"time"

	"yotei-backend/database"
	"yotei-backend/models"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// CreateEventRequest はイベント作成リクエストの構造体
type CreateEventRequest struct {
	Title          string   `json:"title" validate:"required"`
	CandidateDates []string `json:"candidate_dates" validate:"required,min=1"` // ISO 8601形式の日時文字列の配列
}

// CreateEventResponse はイベント作成レスポンスの構造体
type CreateEventResponse struct {
	ID             string                  `json:"id"`
	Title          string                  `json:"title"`
	CandidateDates []CandidateDateResponse `json:"candidate_dates"`
	CreatedAt      time.Time               `json:"created_at"`
}

// CandidateDateResponse は候補日のレスポンス構造体
type CandidateDateResponse struct {
	ID       uint      `json:"id"`
	DateTime time.Time `json:"date_time"`
}

// CreateEvent はイベントと候補日を作成する
func CreateEvent(c *fiber.Ctx) error {
	var req CreateEventRequest

	// リクエストボディをパース
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "リクエストの形式が正しくありません",
		})
	}

	// バリデーション
	if req.Title == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "タイトルは必須です",
		})
	}

	if len(req.CandidateDates) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "候補日を最低1つ指定してください",
		})
	}

	// UUIDを生成
	eventID := uuid.New().String()

	// 候補日をパース
	var candidateDates []models.CandidateDate
	for _, dateStr := range req.CandidateDates {
		parsedTime, err := time.Parse(time.RFC3339, dateStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "候補日の形式が正しくありません。ISO 8601形式（例: 2024-01-01T10:00:00Z）で指定してください",
			})
		}

		candidateDates = append(candidateDates, models.CandidateDate{
			EventID:  eventID,
			DateTime: parsedTime,
		})
	}

	// イベントを作成
	event := models.Event{
		ID:             eventID,
		Title:          req.Title,
		CandidateDates: candidateDates,
	}

	// データベースに保存
	if err := database.DB.Create(&event).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "イベントの作成に失敗しました",
		})
	}

	// レスポンスを作成
	var candidateDateResponses []CandidateDateResponse
	for _, cd := range event.CandidateDates {
		candidateDateResponses = append(candidateDateResponses, CandidateDateResponse{
			ID:       cd.ID,
			DateTime: cd.DateTime,
		})
	}

	response := CreateEventResponse{
		ID:             event.ID,
		Title:          event.Title,
		CandidateDates: candidateDateResponses,
		CreatedAt:      event.CreatedAt,
	}

	return c.Status(fiber.StatusCreated).JSON(response)
}

// GetEvent はイベントの詳細を取得する
func GetEvent(c *fiber.Ctx) error {
	eventID := c.Params("id")

	var event models.Event
	if err := database.DB.Preload("CandidateDates").Preload("Participants").First(&event, "id = ?", eventID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "イベントが見つかりません",
		})
	}

	return c.JSON(event)
}

// // GetAllEvents はすべてのイベントを取得する
// func GetAllEvents(c *fiber.Ctx) error {
// 	var events []models.Event
// 	if err := database.DB.Preload("CandidateDates").Find(&events).Error; err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 			"error": "イベントの取得に失敗しました",
// 		})
// 	}

// 	return c.JSON(events)
// }
