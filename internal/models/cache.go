package models

const (
	CacheCommandRemove CacheCommand = "REMOVE"
	CacheCommandPurge  CacheCommand = "PURGE"
	CacheCommandAdd    CacheCommand = "ADD"
)

type CacheCommand string

type CacheMsg struct {
	Command CacheCommand `json:"command"`
	Key     interface{}  `json:"key"`
}
