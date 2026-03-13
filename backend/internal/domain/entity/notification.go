package entity

import "time"

const (
	NotifTypeTimesheetReminder = "TIMESHEET_REMINDER"
	NotifTypeEffortOverrun     = "EFFORT_OVERRUN"
	NotifTypeDeadlineSoon      = "DEADLINE_SOON"
	NotifTypeIssueStale        = "ISSUE_STALE"
)

const (
	NotifRefProject   = "project"
	NotifRefWBSNode   = "wbs_node"
	NotifRefIssue     = "issue"
	NotifRefTimesheet = "timesheet"
)

type Notification struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Type      string    `json:"type"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	RefID     *int      `json:"ref_id,omitempty"`
	RefType   *string   `json:"ref_type,omitempty"`
	IsRead    bool      `json:"is_read"`
	CreatedAt time.Time `json:"created_at"`
}
