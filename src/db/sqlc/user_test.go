package db

import (
	"context"
	"testing"
	"time"

	"github.com/rhmnmbr/fling-service/util"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {
	hashedPassword, err := util.HashPassword(util.RandomString(6))
	require.NoError(t, err)

	arg := CreateUserParams{
		Email:          util.RandomEmail(),
		HashedPassword: hashedPassword,
		Phone:          util.RandomPhone(),
		FirstName:      util.RandomName(),
		BirthDate:      util.RandomBirthDate(),
		Gender:         GenderEnum(util.RandomGender()),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Email, user.Email)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.Phone, user.Phone)
	require.Equal(t, arg.FirstName, user.FirstName)
	require.Equal(t, arg.BirthDate.Format("2006-01-02"), user.BirthDate.Format("2006-01-02"))
	require.Equal(t, arg.Gender, user.Gender)

	require.NotZero(t, user.CreatedAt)

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user1 := createRandomUser(t)
	user2, err := testQueries.GetUserByEmail(context.Background(), user1.Email)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.Email, user2.Email)
	require.Equal(t, user1.HashedPassword, user2.HashedPassword)
	require.Equal(t, user1.Phone, user2.Phone)
	require.Equal(t, user1.FirstName, user2.FirstName)
	require.Equal(t, user1.BirthDate.Format("2006-01-02"), user2.BirthDate.Format("2006-01-02"))
	require.Equal(t, user1.Gender, user2.Gender)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
}
