package jyoro

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"time"
)

type HMS struct {
	Hour   uint
	Minute uint
	Second uint
}

func (hms *HMS) MarshalJSON() ([]byte, error) {
	return json.Marshal(fmt.Sprintf("%02d:%02d:%02d", hms.Hour, hms.Minute, hms.Second))
}

func (hms *HMS) UnmarshalJSON(data []byte) error {
	var raw string
	if err := json.Unmarshal(data, &raw); err != nil {
		return fmt.Errorf("invalid time format: %s, %w", raw, err)
	}

	r := regexp.MustCompile(`^(\d+):(\d+):(\d+)$`)
	matches := r.FindStringSubmatch(raw)
	if len(matches) != 4 {
		return fmt.Errorf("invalid time format: %s", raw)
	}

	h, err := parseUintRange(matches[1], 23)
	if err != nil {
		return fmt.Errorf("invalid time format: %s, %w", raw, err)
	}
	m, err := parseUintRange(matches[2], 59)
	if err != nil {
		return fmt.Errorf("invalid time format: %s, %w", raw, err)
	}
	s, err := parseUintRange(matches[3], 59)
	if err != nil {
		return fmt.Errorf("invalid time format: %s, %w", raw, err)
	}

	hms.Hour = h
	hms.Minute = m
	hms.Second = s

	return nil
}

func parseUintRange(raw string, max uint) (uint, error) {
	v, err := strconv.ParseUint(raw, 10, 32)
	if err != nil {
		return 0, err
	}

	if v > uint64(max) {
		return 0, fmt.Errorf("value %d is out of range", v)
	}

	return uint(v), nil
}

type Duration struct {
	time.Duration
}

func (d *Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

func (d *Duration) UnmarshalJSON(data []byte) error {
	var v interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch value := v.(type) {
	case string:
		var err error
		d.Duration, err = time.ParseDuration(value)
		if err != nil {
			return err
		}
		return nil
	default:
		return errors.New("invalid duration")
	}
}

type Entry struct {
	Location string   `json:"location"`
	Port     uint     `json:"port"`
	StartAt  HMS      `json:"start_at"`
	Duration Duration `json:"duration"`
}

func (e *Entry) Match(t *time.Time, l *time.Location) bool {
	start := time.Date(t.Year(), t.Month(), t.Day(), int(e.StartAt.Hour), int(e.StartAt.Minute), int(e.StartAt.Second), 0, l)
	end := start.Add(e.Duration.Duration)
	return t.After(start) && t.Before(end)
}

type Config struct {
	Location *time.Location `json:"location"`
	Entries  []Entry        `json:"entries"`
}

func (c *Config) MarshalJSON() ([]byte, error) {
	type Alias Config
	return json.Marshal(&struct {
		Location string `json:"location"`
		*Alias
	}{
		Location: c.Location.String(),
		Alias:    (*Alias)(c),
	})
}

func (c *Config) UnmarshalJSON(data []byte) error {
	type Alias Config
	aux := &struct {
		Location string `json:"location"`
		*Alias
	}{
		Alias: (*Alias)(c),
	}
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}

	if len(aux.Location) == 0 {
		c.Location = time.UTC
		return nil
	}

	loc, err := time.LoadLocation(aux.Location)
	if err != nil {
		return err
	}
	c.Location = loc

	return nil
}

func LoadConfig(file string) (*Config, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	config := &Config{}
	if err := json.Unmarshal(data, config); err != nil {
		return nil, err
	}

	return config, nil
}
