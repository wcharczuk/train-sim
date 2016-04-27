package simulation

import (
	. "testing"

	"github.com/blendlabs/go-assert" //OH SHIT THIS IS A THING??
)

func Test(t *T) {
	assert := assert.New(t)
	g := NewGaussian(3.0, 1)

	assert.NotZero(g.Pdf(1))
}
