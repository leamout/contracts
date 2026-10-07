package agent

import (
	"fmt"
	"strings"
	"time"
)

type RoutingManifest struct {
	SchemaVersion int                 `json:"schema_version"`
	Name          string              `json:"name"`
	Description   string              `json:"description,omitempty"`
	DefaultAgent  string              `json:"default_agent"`
	Routes        []Route             `json:"routes,omitempty"`
	Schedules     map[string]Schedule `json:"schedules,omitempty"`
}

type Route struct {
	Schedule string `json:"schedule"`
	Agent    string `json:"agent"`
}

type Schedule struct {
	Timezone string         `json:"timezone"`
	Hours    []ScheduleHour `json:"hours"`
}

type ScheduleHour struct {
	Days  []string `json:"days"`
	Start string   `json:"start"`
	End   string   `json:"end"`
}

func (m RoutingManifest) Validate() error {
	if m.SchemaVersion != SchemaVersion {
		return fmt.Errorf("unsupported routing schema_version %d", m.SchemaVersion)
	}
	if strings.TrimSpace(m.Name) == "" {
		return fmt.Errorf("routing name is required")
	}
	if strings.TrimSpace(m.DefaultAgent) == "" {
		return fmt.Errorf("routing default_agent is required")
	}

	for name, schedule := range m.Schedules {
		if strings.TrimSpace(name) == "" {
			return fmt.Errorf("schedule name is required")
		}
		if err := schedule.Validate(); err != nil {
			return fmt.Errorf("schedule %q: %w", name, err)
		}
	}

	for _, route := range m.Routes {
		if strings.TrimSpace(route.Schedule) == "" {
			return fmt.Errorf("route schedule is required")
		}
		if strings.TrimSpace(route.Agent) == "" {
			return fmt.Errorf("route agent is required")
		}
		if _, ok := m.Schedules[route.Schedule]; !ok {
			return fmt.Errorf("route references unknown schedule %q", route.Schedule)
		}
	}
	return nil
}

func (m RoutingManifest) ValidateAgentReferences(aliases map[string]struct{}) error {
	if _, ok := aliases[m.DefaultAgent]; !ok {
		return fmt.Errorf("routing references unknown default agent %q", m.DefaultAgent)
	}
	for _, route := range m.Routes {
		if _, ok := aliases[route.Agent]; !ok {
			return fmt.Errorf("routing references unknown agent %q", route.Agent)
		}
	}
	return nil
}

func (s Schedule) Validate() error {
	if strings.TrimSpace(s.Timezone) == "" {
		return fmt.Errorf("timezone is required")
	}
	if _, err := time.LoadLocation(s.Timezone); err != nil {
		return fmt.Errorf("invalid timezone %q", s.Timezone)
	}
	if len(s.Hours) == 0 {
		return fmt.Errorf("hours are required")
	}
	for _, hours := range s.Hours {
		if err := hours.Validate(); err != nil {
			return err
		}
	}
	return nil
}

func (h ScheduleHour) Validate() error {
	if len(h.Days) == 0 {
		return fmt.Errorf("schedule days are required")
	}
	for _, day := range h.Days {
		switch day {
		case "mon", "tue", "wed", "thu", "fri", "sat", "sun":
		default:
			return fmt.Errorf("invalid schedule day %q", day)
		}
	}
	start, err := time.Parse("15:04", h.Start)
	if err != nil {
		return fmt.Errorf("invalid schedule start %q", h.Start)
	}
	end, err := time.Parse("15:04", h.End)
	if err != nil {
		return fmt.Errorf("invalid schedule end %q", h.End)
	}
	if !end.After(start) {
		return fmt.Errorf("schedule end must be after start")
	}
	return nil
}
