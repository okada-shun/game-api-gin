package database

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (s *DatabaseSuite) TestGetGachaIds() {
	expectedIds := make([]int, 0)
	for i := 1; i < 4; i++ {
		for j := 0; j < 10; j++ {
			expectedIds = append(expectedIds, i)
		}
	}
	actualIds, err := s.db.GetGachaIds()
	require.NoError(s.T(), err)
	assert.Equal(s.T(), expectedIds, actualIds)
}
