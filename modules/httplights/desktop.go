package httplights

import (
	"net/http"
	"time"

	"github.com/collinux/gohue"
	"github.com/gorilla/mux"
)

// DesktopTriggerSchedule provides the start time and duration that the trigger
// should be active.
type DesktopTriggerSchedule struct {
	// StartTime should be in time.Kitchen format. This is the time after which
	// triggering the desktop lights switch will become active
	StartTime string

	// ActiveFor is the duration after the start time that the scene may be
	// triggered.
	ActiveFor time.Duration
}

// DesktopTrigger is a HTTP light module that adds a endpoint that should be
// triggered by putting the desktop computer into standby mode. If it's past
// the start time and within the duration of the schedule, the scene will be
// triggered.
//
// Good for night-light triggering.
type DesktopTrigger struct {
	hueBridge *hue.Bridge

	// SceneName specifies the name of the scene to trigger.
	SceneName string

	// Schedule configures when the trigger should be active for a specific day
	// of the week.
	Schedules map[time.Weekday]DesktopTriggerSchedule

	// DefaultTime is the default schedule that should be followed when no
	// specific weekday schedule is configured.
	DefaultSchedule DesktopTriggerSchedule
}

func (t *DesktopTrigger) getSchedule(date time.Time) DesktopTriggerSchedule {
	schedule, hasDaySchedule := t.Schedules[date.Weekday()]

	if !hasDaySchedule {
		schedule = t.DefaultSchedule
	}

	return schedule
}

// Check if we're within the schedule for the given date. A 'now' date should
// be provided which will be used as the reference time
func (t *DesktopTrigger) withinSchedule(date, now time.Time) bool {
	schedule := t.getSchedule(date)
	if schedule.StartTime == "" {
		return false
	}

	start, err := time.Parse(time.Kitchen, schedule.StartTime)
	if err != nil {
		return false
	}

	// Calculate start and end times for the schedule
	y, mon, d := date.Date()
	h, min, s := start.Clock()

	startTime := time.Date(y, mon, d, h, min, s, 0, date.Location())
	endingTime := startTime.Add(schedule.ActiveFor)

	return now.After(startTime) && now.Before(endingTime)
}

// SetHueBridge implements the HTTPLightsModule interface.
func (t *DesktopTrigger) SetHueBridge(bridge *hue.Bridge) {
	t.hueBridge = bridge
}

// RegisterInRouter implements the HTTPLightsModule interface.
func (t *DesktopTrigger) RegisterInRouter(router *mux.Router) {
	router.Handle("/desktop-standby", t).Methods("POST")
}

// ServeHTTP implements the HTTPLightsModule and http.Handler interface.
func (t *DesktopTrigger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	yesterday := now.AddDate(0, 0, -1)

	// Check if we're within the schedule for today or from yesterday. We check
	// yesterdays schedule for the case where it crosses the midnight boundary
	// to today.
	if !t.withinSchedule(now, now) && !t.withinSchedule(yesterday, now) {
		return
	}

	// Trigger lights
	t.hueBridge.RecallSceneByName(t.SceneName)
}
