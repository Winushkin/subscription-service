package queries_test

import (
	"errors"
	"subscription-service/internal/entities"
	"subscription-service/internal/repository/queries"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestFullSubSelect(t *testing.T) {
	t.Run("check default sql string", func(t *testing.T) {
		sql, args, err := queries.FullSubscriptionSelect().ToSql()
		expected := "SELECT id, service_name, price, start_date, end_date, created_at, updated_at FROM subscriptions"
		require.NoError(t, err)
		require.Empty(t, args, "Слайс аргументов селекта должен быть пустым")
		require.Equal(t, expected, sql, "Ожидаемый и полученный результаты не совпадают")
	})
}

func TestUpdateSubscription(t *testing.T) {
	testName, testPrice, testStartDate, testEndDate := "NewName", 100, "05-2026",
		"05-2026"
	userID, _ := uuid.Parse("60601fee-2bf1-4721-ae6f-7636e79a0cba")

	tests := []struct {
		name          string
		updates       entities.UpdateSubscriptionRequest
		expectedError error
		expectedArgs  []any
		expectedSQL   string
	}{
		{
			name: "full updates",
			updates: entities.UpdateSubscriptionRequest{
				ServiceName: &testName,
				Price:       &testPrice,
				StartDate:   &testStartDate,
				EndDate:     &testEndDate,
			},
			expectedError: nil,
			expectedArgs: []any{
				testName,
				testPrice,
				testStartDate,
				testEndDate,
				userID.String(),
			},
			expectedSQL: "UPDATE subscriptions SET service_name = $1, price = $2, start_date = $3, end_date = $4 WHERE id = $5",
		},
		{
			name:          "empty updates",
			updates:       entities.UpdateSubscriptionRequest{},
			expectedError: errors.New("update statements must have at least one Set clause"),
			expectedArgs:  []any{},
			expectedSQL:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sql, args, err := queries.UpdateSubscription(userID, tt.updates).ToSql()

			if tt.expectedError != nil {
				require.Error(t, err)
				require.Equal(t, tt.expectedError.Error(), err.Error())
				return 
			}
			
			require.NoError(t, err)
			require.Equal(t, tt.expectedSQL, sql)
			require.Equal(t, tt.expectedArgs, args)
		})
	}
}
