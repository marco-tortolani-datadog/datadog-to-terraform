package types

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)
var (
	REQUIRED_TAGS = []string{"team:${var.team}", "env:${var.env_short}"}
)
type ThresholdCount struct {
	Ok               *json.Number `json:"ok,omitempty" hcl:"ok"`
	Critical         *json.Number `json:"critical,omitempty" hcl:"critical"`
	Warning          *json.Number `json:"warning,omitempty" hcl:"warning"`
	Unknown          *json.Number `json:"unknown,omitempty" hcl:"unknown"`
	CriticalRecovery *json.Number `json:"critical_recovery,omitempty" hcl:"critical_recovery"`
	WarningRecovery  *json.Number `json:"warning_recovery,omitempty" hcl:"warning_recovery"`
}

type ThresholdWindows struct {
	RecoveryWindow *string `json:"recovery_window,omitempty" hcl:"recovery_window"`
	TriggerWindow  *string `json:"trigger_window,omitempty" hcl:"trigger_window"`
}

type NoDataTimeframe int

func (tf *NoDataTimeframe) UnmarshalJSON(data []byte) error {
	s := string(data)
	if s == "false" || s == "null" {
		*tf = 0
	} else {
		i, err := strconv.ParseInt(s, 10, 32)
		if err != nil {
			return err
		}
		*tf = NoDataTimeframe(i)
	}
	return nil
}

type Options struct {
	// NoDataTimeframe   NoDataTimeframe   `json:"no_data_timeframe,omitempty" hcl:"no_data_timeframe"`
	// NotifyAudit       *bool             `json:"notify_audit,omitempty" hcl:"notify_audit"`
	// NotifyNoData      *bool             `json:"notify_no_data,omitempty" hcl:"notify_no_data"`
	// RenotifyInterval  *int              `json:"renotify_interval,omitempty" hcl:"renotify_interval"`
	// NewHostDelay      *int              `json:"new_host_delay,omitempty" hcl:"new_host_delay"`
	// EvaluationDelay   *int              `json:"evaluation_delay,omitempty" hcl:"evaluation_delay"`
	// TimeoutH          *int              `json:"timeout_h,omitempty" hcl:"timeout_h"`
	// EscalationMessage *string           `json:"escalation_message,omitempty" hcl:"escalation_message"`
	// Thresholds        *ThresholdCount   `json:"thresholds,omitempty" hcl:"monitor_thresholds"`
	// ThresholdWindows  *ThresholdWindows `json:"threshold_windows,omitempty" hcl:"monitor_threshold_windows"`
	// IncludeTags       *bool             `json:"include_tags,omitempty" hcl:"include_tags"`
	// RequireFullWindow *bool             `json:"require_full_window,omitempty" hcl:"require_full_window"`
	// Locked            *bool             `json:"locked,omitempty" hcl:"locked"`
	// EnableLogsSample  *bool             `json:"enable_logs_sample,omitempty" hcl:"enable_logs_sample"`
}

// Monitor allows watching a metric or check that you care about,
// notifying your team when some defined threshold is exceeded
type Monitor struct {
	Type     *string  `json:"type,omitempty" hcl:"type"`
	Query    *string  `json:"query,omitempty" hcl:"query"`
	Name     *string  `json:"name,omitempty" hcl:"name"`
	Message  *string  `json:"message,omitempty" hcl:"message"`
	Tags     []string `json:"tags" hcl:"tags"`
	Priority *int     `json:"priority" hcl:"priority"`
	*Options `json:"options,omitempty" hcl:",squash"`
}

// Write a function for the monitor struct that adds the required tags to the monitor.tags array
func (m *Monitor) AddRequiredTags() {
	m.Tags = append(m.Tags, REQUIRED_TAGS...)
}

// Write a function that returns the name of a monitor in lowercase, removing emojis, and replacing spaces with underscores
func (m *Monitor) GetLowercaseName() string {
	return strings.ToLower(strings.ReplaceAll(*m.Name, " ", "_"))
}

func (m *Monitor) MakeQueryHeredoc() {
	*m.Query =  "<<-EOF\n" + strings.Trim(*m.Query, "\"") + "\nEOF"
}

// Write a function that takes the multiline string monitor.Message and converts it into heredoc format by stripping the leading and trailing quotation marks and adding a leading and trailing EOF markers
func (m *Monitor) MakeMessageHeredoc() {
	*m.Message =  "<<-EOF\n" + strings.Trim(*m.Message, "\"") + "\nEOF"
}

// Write a function that requests user input on the command line for whether the monitor should be muted, and if so it should add a tag "tf-mute:true" to the monitor.tags array
func (m *Monitor) AskForMuteTag() {
	var mute string
	for {
		print("Would you like to mute this monitor (tf-mute:true)? (y/n): ")
		_, err := fmt.Scan(&mute)
		if err != nil {
			fmt.Println("Please enter a valid input")
			continue
		}
		if strings.ToLower(mute) == "y" || strings.ToLower(mute) == "yes" {
			m.Tags = append(m.Tags, "tf-mute:true")
		}
		break
	}
}

// write a function that the user for a priority level for an integer 1-5 or the option to skip for the monitor.priority field if it is not already set
func (m *Monitor) AskForPriority() {
	if m.Priority != nil {
		return
	}
	var priority int
	for {
		print("Enter a priority level for this monitor (1-5) or 0 to skip: ")
		_, err := fmt.Scan(&priority)
		if err != nil {
			fmt.Println("Invalid input, skipping priority")
			break
		}
		if priority >= 1 && priority <= 6 {
			m.Priority = &priority
		}
		break
	}
}

// write a function that checks for this string "${var.slack_channel_notify}" in the monitor.Message field and if it is not present it should add it to the end of the string in a new line
func (m *Monitor) AddSlackChannelNotify() {
	if !strings.Contains(*m.Message, "${var.slack_channel_notify}") {
		*m.Message += "\t${var.slack_channel_notify}"
	}
}