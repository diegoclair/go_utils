package dateutils

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetFirstAndLastOfAMonth(t *testing.T) {
	year, _, _ := GetDateNowTime().Date()

	// Test case for January
	firstDay, lastDay := GetFirstAndLastOfAMonth(time.January)
	firstDayExpected := fmt.Sprintf("%d-01-01", year)
	lastDayExpected := fmt.Sprintf("%d-01-31", year)
	assert.Equal(t, firstDayExpected, firstDay)
	assert.Equal(t, lastDayExpected, lastDay)
}
