package answer

type Answer struct {
	URL    string
	Method string
}

type PgRepository struct {
	connector string
}

func NewRepository(connector string) *PgRepository {
	return &PgRepository{
		connector: connector,
	}
}

func (r *PgRepository) GetAnswer(URL string, Method string) *Answer {
	return &Answer{
		URL,
		Method,
	}
}

func (r *PgRepository) VerifyAnswer(_ string, _ string, _ bool) error {
	return nil
}
