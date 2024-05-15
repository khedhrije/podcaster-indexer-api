package retry

import (
	"time"

	"github.com/rs/zerolog/log"
)

type FunctionToRetry func() (interface{}, error)

// ExecuteWithBackoffRetry execute the function passed in parameters with a number of retries and a fixed backoff delay.
// It returns the number of function calls (= 1 if everything was ok) and the error only if none of the calls worked
func ExecuteWithBackoffRetry(function FunctionToRetry, maxRetry int, delay time.Duration) (interface{}, int, error) {
	ticker := time.NewTicker(delay)
	defer ticker.Stop()
	result, err := function()
	nbTry := 1
	isNextCallWorked := false

	if err != nil {
		log.Debug().Msgf("Call is not ok after 1 try, will retry %d times", maxRetry)

		for ; nbTry <= maxRetry && !isNextCallWorked; nbTry++ {
			select {
			case <-ticker.C:
				{
					result, err = function()
					if err == nil {
						isNextCallWorked = true
						log.Debug().Msgf("Call OK after %d try", nbTry+1)
					} else {
						log.Debug().Msgf("Call is not ok after %d try", nbTry+1)
					}
				}
			}
		}
		if isNextCallWorked {
			err = nil
		}
	}
	return result, nbTry, err
}
