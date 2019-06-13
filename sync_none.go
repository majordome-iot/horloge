package horloge

type SyncNone struct {
}

func NewSyncNone() *SyncNone {
	return &SyncNone{}
}

func (s *SyncNone) Read() []Job {
	return nil
}

func (s *SyncNone) Write(jobs []*Job) error {
	return nil
}
