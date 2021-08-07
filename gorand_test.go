package gorand_test

import (
	"math/rand"
	"reflect"
	"testing"
	"time"

	"github.com/loghole/gorand"
)

func BenchmarkRand(b *testing.B) {
	bc := 4

	for i := 0; i < b.N; i++ {
		rand.Int()
	}

	b.Run("pkg", func(b *testing.B) {
		b.ResetTimer()

		b.SetParallelism(bc)

		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				rand.Int()
			}
		})
	})

	b.Run("custom", func(b *testing.B) {
		rnd := rand.New(gorand.NewSource(time.Now().UnixNano()))

		b.SetParallelism(bc)

		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				rnd.Int()
			}
		})
	})
}

func TestNewSource(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		want []int64
		seed int64
	}{
		{
			name: "pass",
			seed: 1,
			want: []int64{8855534638979991142},
		},
	}
	for _, rt := range tests {
		tt := rt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var (
				src = gorand.NewSource(tt.seed)
				got = make([]int64, 0, len(tt.want))
			)

			for i := 0; i < len(tt.want); i++ {
				got = append(got, src.Int63())
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSource() = %v, want %v", got, tt.want)
			}
		})
	}
}
