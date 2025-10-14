package usecase_test

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/dto"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/usecase"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/usecase/mocks"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/pkg/transaction"
	"golang.org/x/net/context"
)

func Test_CreateProfileV2_Success(t *testing.T) {
	// otel.SilentModeInit() // отключить open telemetry
	transaction.IsUnitTest = true // отключить транзакции

	// Настраиваем поведение Postgres
	postgres := new(mocks.Postgres) // Создаётся мок-объект для интерфейса Postgres, сгенерированный с помощью mockery
	postgres.On("CreateProfile", mock.Anything, mock.Anything).Return(nil)
	postgres.On("CreateProperty", mock.Anything, mock.Anything).Return(nil)

	// Эти defer гарантируют, что после выполнения теста будет проверено, действительно ли методы мока вызывались.
	defer postgres.AssertCalled(t, "CreateProfile", mock.Anything, mock.Anything)
	defer postgres.AssertCalled(t, "CreateProperty", mock.Anything, mock.Anything)

	// Создаётся экземпляр UseCase с внедрённым моком Postgres.
	u := usecase.New(postgres)

	{ // Сам тест
		input := dto.CreateProfileInput{
			Name:  "John Doe",
			Age:   30,
			Email: "john.doe@example.com",
			Phone: "+1234567890",
		}

		// Вызов тестируемого метода CreateProfile
		// Вместо реальной БД используется мок, и транзакции отключены.
		// Метод должен успешно создать профиль и вернуть ID.
		actual, err := u.CreateProfile(context.Background(), input)
		require.NoError(t, err)        // Проверяет, что метод не вернул ошибку.
		require.NotEmpty(t, actual.ID) // Проверяет, что в ответе вернулся непустой ID профиля.
	}
}
