package api_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/krestkrest/word-of-wisdom/api"
)

func testChallenge(complexity uint8) *api.MessageChallenge {
	return &api.MessageChallenge{
		Address:    "192.168.0.1:53011",
		Nonce:      api.Nonce{},
		UnixTime:   time.Now().UnixNano(),
		Complexity: complexity,
	}
}

func TestMessageChallenge_Check(t *testing.T) {
	challenge := testChallenge(9)
	solution := challenge.FindSolution()
	response := api.NewMessageResponse(challenge, solution)
	assert.NoError(t, response.CheckSolution())
}

func benchmarkMessageChallengeFindSolution(b *testing.B, complexity uint8) {
	challenge := testChallenge(complexity)
	for i := 0; i < b.N; i++ {
		challenge.FindSolution()
	}
}

func BenchmarkMessageChallenge_FindSolution8(b *testing.B) {
	benchmarkMessageChallengeFindSolution(b, 8)
}

func BenchmarkMessageChallenge_FindSolution16(b *testing.B) {
	benchmarkMessageChallengeFindSolution(b, 16)
}

func BenchmarkMessageChallenge_FindSolution20(b *testing.B) {
	benchmarkMessageChallengeFindSolution(b, 20)
}

func BenchmarkMessageChallenge_FindSolution22(b *testing.B) {
	benchmarkMessageChallengeFindSolution(b, 22)
}
