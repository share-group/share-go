package service

type systemService struct{}

var SystemService = newSystemService()

func newSystemService() *systemService {
	return &systemService{}
}

func (s *systemService) Get(id int64) int {
	return 1
}
