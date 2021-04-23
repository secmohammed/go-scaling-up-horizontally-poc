package entity

import "time"

//LogLevel is used to indicate that type is a string
type LogLevel string

const (
    //LogLevelInfo is used to declare a constant containing the info level
    LogLeveLInfo LogLevel = "INFO"
    //LogLevelError is used to declare a constant containing the error level
    LogLevelError LogLevel = "ERROR"
    //LogLevelPanic is used to declare a constant containing the error level
    LogLevelPanic LogLevel = "PANIC"
)

//LogEntry is used to structure the logentry attributes
type LogEntry struct {
    Level     LogLevel
    Timestamp time.Time
    Source    string
    Message   string
}
