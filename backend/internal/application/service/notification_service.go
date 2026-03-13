package service

import (
	"context"
	"fmt"
	"log"
	"project-mgmt/backend/internal/domain/entity"
	"project-mgmt/backend/internal/domain/repository"
	"time"
)

type NotificationService interface {
	Create(ctx context.Context, userID int, nType, title, body string, refID *int, refType *string) error
	ListByUser(ctx context.Context, userID, limit, offset int) ([]entity.Notification, int, error)
	GetUnreadCount(ctx context.Context, userID int) (int, error)
	MarkRead(ctx context.Context, id int) error
	MarkAllRead(ctx context.Context, userID int) error
}

type notificationService struct {
	repo     repository.NotificationRepository
	userRepo repository.UserRepository
}

func NewNotificationService(
	repo repository.NotificationRepository,
	userRepo repository.UserRepository,
) NotificationService {
	return &notificationService{repo: repo, userRepo: userRepo}
}

func (s *notificationService) Create(ctx context.Context, userID int, nType, title, body string, refID *int, refType *string) error {
	n := &entity.Notification{
		UserID:  userID,
		Type:    nType,
		Title:   title,
		Body:    body,
		RefID:   refID,
		RefType: refType,
	}
	return s.repo.Create(ctx, n)
}

func (s *notificationService) ListByUser(ctx context.Context, userID, limit, offset int) ([]entity.Notification, int, error) {
	return s.repo.ListByUser(ctx, userID, limit, offset)
}

func (s *notificationService) GetUnreadCount(ctx context.Context, userID int) (int, error) {
	return s.repo.GetUnreadCount(ctx, userID)
}

func (s *notificationService) MarkRead(ctx context.Context, id int) error {
	return s.repo.MarkRead(ctx, id)
}

func (s *notificationService) MarkAllRead(ctx context.Context, userID int) error {
	return s.repo.MarkAllRead(ctx, userID)
}

// StartCronJobs launches background goroutine-based periodic jobs.
// Call once at server startup.
func StartCronJobs(notifSvc NotificationService, userRepo repository.UserRepository) {
	go runCron(notifSvc, userRepo)
	log.Println("[CRON] Notification cron scheduler started")
}

func runCron(notifSvc NotificationService, userRepo repository.UserRepository) {
	// Run immediately then tick
	checkProjectDeadlines(notifSvc, userRepo)
	checkStaleIssues(notifSvc, userRepo)

	// Tick every 24 hours for daily checks
	dailyTicker := time.NewTicker(24 * time.Hour)
	// Tick every Friday for weekly checks
	weeklyTicker := time.NewTicker(7 * 24 * time.Hour)

	for {
		select {
		case <-dailyTicker.C:
			checkProjectDeadlines(notifSvc, userRepo)
			checkStaleIssues(notifSvc, userRepo)
		case <-weeklyTicker.C:
			checkTimesheetReminders(notifSvc, userRepo)
		}
	}
}

// checkTimesheetReminders notifies all active users on Fridays to submit their timesheets.
func checkTimesheetReminders(notifSvc NotificationService, userRepo repository.UserRepository) {
	ctx := context.Background()
	now := time.Now()
	// Only run on Friday (weekday == 5)
	if now.Weekday() != time.Friday {
		return
	}

	users, err := userRepo.List(ctx)
	if err != nil {
		log.Printf("[CRON] Failed to list users for timesheet reminder: %v", err)
		return
	}

	for _, u := range users {
		if !u.IsActive {
			continue
		}
		_ = notifSvc.Create(ctx, int(u.ID),
			entity.NotifTypeTimesheetReminder,
			"Nhắc nhở: Submit Timesheet",
			fmt.Sprintf("Hôm nay là thứ 6 (%s). Vui lòng ghi chép và submit giờ làm việc trong tuần này.", now.Format("02/01/2006")),
			nil, nil,
		)
	}
	log.Printf("[CRON] Friday timesheet reminders sent to %d users", len(users))
}

// checkProjectDeadlines sends notifications for projects with deadlines within 7 days.
func checkProjectDeadlines(notifSvc NotificationService, userRepo repository.UserRepository) {
	// Placeholder — would query projects table and notify PM and members
	// Requires projectRepo to be passed in; skipped for now to keep DI simple
	log.Println("[CRON] checkProjectDeadlines: placeholder (requires project repo)")
}

// checkStaleIssues notifies about open issues not updated in 14+ days.
func checkStaleIssues(notifSvc NotificationService, userRepo repository.UserRepository) {
	// Placeholder — would query issues table
	log.Println("[CRON] checkStaleIssues: placeholder (requires issue repo)")
}
