//go:generate mockgen -source ./service.go -destination=./mocks/sender.go -package=mock_sender

package answer

import "fmt"

type Sender interface {
	sendAsyncMessage(message DeleteMessage) error
	sendMessage(message DeleteMessage) error
	sendMessages(messages []DeleteMessage) error
}

type Repository interface {
	GetAnswer(URL string, Method string) *Answer
	VerifyAnswer(URL string, Method string, success bool) error
}

type Service struct {
	repo   Repository
	sender Sender
}

func NewService(repo Repository, sender Sender) *Service {
	return &Service{
		repo:   repo,
		sender: sender,
	}
}

func (s Service) Verify(answerURL string, answerMethod string, success bool, sync bool) {
	answer := s.repo.GetAnswer(answerURL, answerMethod)
	err := s.repo.VerifyAnswer(answerURL, answerMethod, success)
	if err != nil {
		fmt.Println("Error in verify", err)
	}

	if sync {
		err = s.sender.sendMessage(
			DeleteMessage{
				answer.URL,
				answer.Method,
				success,
			},
		)

		if err != nil {
			fmt.Println("Send sync message error: ", err)
		}

		return
	}

	err = s.sender.sendAsyncMessage(
		DeleteMessage{
			answer.URL,
			answer.Method,
			success,
		},
	)

	if err != nil {
		fmt.Println("Send async message error: ", err)
	}

}
