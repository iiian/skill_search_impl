package db

type Candidate struct {
	ID     string
	Skills []string
}

type Database interface {
	Add(...*Candidate)
	Search([]string) (*Candidate, bool)
}

type SearchCache interface {
	Get([]string) (*CacheRecord, bool)
	Put([]string, *CacheRecord)
	HandleCacheUpdate(*Candidate)
}
