package retry

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_BackoffRetry_WhenOk(t *testing.T) {
	result, nbTry, err := ExecuteWithBackoffRetry(func() (interface{}, error) {
		return 12, nil
	}, 2, time.Duration(100)*time.Millisecond)

	assert.Equal(t, 1, nbTry, "Should be tried once since it's ok")
	assert.Equal(t, 12, result.(int))
	assert.NoError(t, err, "Should be ok")
}

func Test_BackoffRetry_WhenKo(t *testing.T) {
	_, nbTry, err := ExecuteWithBackoffRetry(func() (interface{}, error) {
		return nil, errors.New("Fake error")
	}, 2, time.Duration(100)*time.Millisecond)

	assert.Equal(t, 3, nbTry, "Should be tried 3 times since it's ko and there are 2 retries")
	assert.Error(t, err, "Should be ko")
}
