package services

type AboutService struct{}

func NewAboutService() *AboutService {
	return &AboutService{}
}

func (service *AboutService) GetWelcomeMessage() string {
	return "Welcome to the website starter! (About)"
}
