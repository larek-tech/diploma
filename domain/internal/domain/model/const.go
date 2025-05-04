package model

// SourceType enum defining the way of parsing source.
type SourceType uint8

const (
	// TypeUndefined undefined source type.
	TypeUndefined SourceType = iota
	// TypeWeb web source.
	TypeWeb
	// TypeSingleFile single file source.
	TypeSingleFile
	// TypeArchivedFiles many archived files source.
	TypeArchivedFiles
	// TypeWithCredentials source with credentials
	TypeWithCredentials
)

// SourceStatus enum defining the status of source parsing.
type SourceStatus uint8

const (
	// StatusUndefined undefined status.
	StatusUndefined SourceStatus = iota
	// StatusReady source is ready.
	StatusReady
	// StatusParsing source is being parsed.
	StatusParsing
	// StatusFailed source parsing failed.
	StatusFailed
)
