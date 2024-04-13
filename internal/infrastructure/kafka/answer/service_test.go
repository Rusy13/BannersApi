package answer

import (
	"errors"
	"testing"
)

// mockRepository реализует интерфейс Repository для тестирования Service.
type mockRepository struct{}

func (m *mockRepository) GetAnswer(URL string, Method string) *Answer {
	return &Answer{
		URL:    URL,
		Method: Method,
	}
}

func (m *mockRepository) VerifyAnswer(URL string, Method string, success bool) error {
	// Возвращаем ошибку только в случае, если URL начинается с "error".
	if URL == "error" {
		return errors.New("mock repository error")
	}
	return nil
}

// mockSender реализует интерфейс Sender для тестирования Service.
type mockSender struct {
	asyncMessageSent  bool
	syncMessageSent   bool
	messagesSentCount int
}

func (m *mockSender) sendAsyncMessage(message DeleteMessage) error {
	m.asyncMessageSent = true
	return nil
}

func (m *mockSender) sendMessage(message DeleteMessage) error {
	m.syncMessageSent = true
	return nil
}

func (m *mockSender) sendMessages(messages []DeleteMessage) error {
	m.messagesSentCount = len(messages)
	return nil
}

func TestService_Verify(t *testing.T) {
	// Создание экземпляров mockRepository и mockSender.
	repo := &mockRepository{}
	sender := &mockSender{}

	// Создание экземпляра сервиса.
	service := NewService(repo, sender)

	// Тестирование синхронной отправки.
	service.Verify("testURL", "testMethod", true, true)
	if !sender.syncMessageSent {
		t.Error("Expected sync message to be sent")
	}

	// Тестирование асинхронной отправки.
	service.Verify("testURL", "testMethod", true, false)
	//time.Sleep(time.Millisecond * 10000) // Подождать 100 миллисекунд
	if !sender.asyncMessageSent {
		t.Error("Expected async message to be sent")
	}

}
