package models

import "encoding/json"

type Call struct {
	Path         string            `bson:"path" json:"path"`
	Method       string            `bson:"method" json:"method"`
	ResponseBody json.RawMessage   `bson:"response_body" json:"response_body"`
	StringBody   string            `bson:"string_body" json:"string_body"`
	ContentType  string            `bson:"content_type" json:"content_type"`
	Header       map[string]string `bson:"header" json:"header"`
}

type Job struct {
	Id   string `bson:"_id" json:"id"`
	Call Call   `bson:"call" json:"call"`

	// settings
	MaxRetries      int   `bson:"maxRetries" json:"maxRetries"`
	Delay           int64 `bson:"delay" json:"delay"`
	Repeat          bool  `bson:"repeat" json:"repeat"`
	InitialPriority int   `bson:"initial_priority" json:"initial_priority"`

	// state
	Tries     int   `bson:"tries" json:"tries"`
	Priority  int   `bson:"priority" json:"priority"`
	Timestamp int64 `bson:"timestamp" json:"timestamp"`
}
