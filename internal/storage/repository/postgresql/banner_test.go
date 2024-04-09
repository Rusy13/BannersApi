package postgresql

import (
	"HW1/internal/storage/repository"
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"testing"
)

func Test_postgresDBRepo_GetByID(t *testing.T) {
	t.Parallel()

	var (
		id = 1
	)

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		// arrange

		s := SetUp(t)
		defer s.TearDown()
		s.mockDB.EXPECT().Get(gomock.Any(), gomock.Any(), "SELECT id,pvzname,address,email FROM pvz where id=$1", gomock.Any()).Return(nil)
		// act
		reservation, err := s.repo.GetByID(context.Background(), int64(id))
		// assert

		require.NoError(t, err)
		assert.Equal(t, int64(0), reservation.ID)
	})

	t.Run("fail", func(t *testing.T) {
		t.Parallel()
		t.Run("not found", func(t *testing.T) {
			t.Parallel()
			// arrange
			s := SetUp(t)
			defer s.TearDown()

			s.mockDB.EXPECT().Get(gomock.Any(), gomock.Any(), "SELECT id,pvzname,address,email FROM pvz where id=$1", gomock.Any()).Return(repository.ErrObjectNotFound)

			// act
			reservation, err := s.repo.GetByID(context.Background(), int64(id))
			// assert
			require.EqualError(t, err, "not found")

			assert.Nil(t, reservation)
		})

		t.Run("internal error", func(t *testing.T) {
			t.Parallel()
			// arrange
			s := SetUp(t)
			defer s.TearDown()

			s.mockDB.EXPECT().Get(gomock.Any(), gomock.Any(), "SELECT id,pvzname,address,email FROM pvz where id=$1", gomock.Any()).Return(assert.AnError)

			// act
			reservation, err := s.repo.GetByID(context.Background(), int64(id))
			// assert
			require.EqualError(t, err, "assert.AnError general error for testing")

			assert.Nil(t, reservation)
		})
	})
}

func Test_Delete(t *testing.T) {
	var (
		id int64 = 1
	)

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		// arrange

		s := SetUp(t)
		defer s.TearDown()
		s.mockDB.EXPECT().Exec(gomock.Any(), "DELETE FROM pvz WHERE id=$1", gomock.Any()).Return(nil, nil)
		// act
		reservation := s.repo.Delete(context.Background(), int64(id))
		// assert

		//require.NoError(t, reservation)
		assert.Equal(t, nil, reservation)
	})

	t.Run("fail", func(t *testing.T) {
		t.Parallel()
		t.Run("not found", func(t *testing.T) {
			t.Parallel()
			// arrange
			s := SetUp(t)
			defer s.TearDown()

			s.mockDB.EXPECT().Exec(gomock.Any(), "DELETE FROM pvz WHERE id=$1", gomock.Any()).Return(nil, repository.ErrObjectNotFound)

			// act
			reservation := s.repo.Delete(context.Background(), int64(id))
			// assert
			require.EqualError(t, reservation, "not found")

			//assert.Nil(t, reservation)
		})

		t.Run("internal error", func(t *testing.T) {
			t.Parallel()
			// arrange
			s := SetUp(t)
			defer s.TearDown()

			s.mockDB.EXPECT().Exec(gomock.Any(), "DELETE FROM pvz WHERE id=$1", gomock.Any()).Return(nil, assert.AnError)

			// act
			reservation := s.repo.Delete(context.Background(), int64(id))
			// assert
			require.EqualError(t, reservation, "assert.AnError general error for testing")

			//assert.Nil(t, reservation)
		})
	})
}
