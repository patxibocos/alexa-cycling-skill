package alexa

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGapMessageWithZeroGap(t *testing.T) {
	gap := 0

	message := getGapMessage(int64(gap))

	assert.Equal(t, "con el mismo tiempo", message)
}

func TestGapMessageWithSingleSecond(t *testing.T) {
	gap := 1

	message := getGapMessage(int64(gap))

	assert.Equal(t, "a 1 segundo", message)
}

func TestGapMessageWithSingleMinuteAndSecond(t *testing.T) {
	gap := 61

	message := getGapMessage(int64(gap))

	assert.Equal(t, "a 1 minuto y 1 segundo", message)
}

func TestGapMessageWithMinutesAndNoSeconds(t *testing.T) {
	gap := 60

	message := getGapMessage(int64(gap))

	assert.Equal(t, "a 1 minuto", message)
}

func TestGapMessageWithManyMinutesAndSeconds(t *testing.T) {
	gap := 122

	message := getGapMessage(int64(gap))

	assert.Equal(t, "a 2 minutos y 2 segundos", message)
}
