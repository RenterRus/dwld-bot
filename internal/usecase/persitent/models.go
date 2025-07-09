package persitent

type Upsert struct {
	name, host       string
	allowedRootLinks []string
	port             int
}

type Servers interface {
	UpsertServer(server *Upsert) error
	LoadServers() ([]*Upsert, error)
}
