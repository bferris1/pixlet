package random

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

const (
	ModuleName    = "random"
	threadRandKey = "tidbyt.dev/pixlet/runtime/random"
)

var (
	once   sync.Once
	module starlark.StringDict
)

func AttachToThread(t *starlark.Thread) {
	t.SetLocal(
		threadRandKey,
		rand.New(
			rand.NewSource(rand.Int63()),
		),
	)
}

func LoadModule() (starlark.StringDict, error) {
	once.Do(func() {
		rand.Seed(time.Now().UnixNano())
		module = starlark.StringDict{
			ModuleName: &starlarkstruct.Module{
				Name: ModuleName,
				Members: starlark.StringDict{
					"number": starlark.NewBuiltin("number", randomNumber),
					"seed":   starlark.NewBuiltin("seed", randomSeed),
				},
			},
		}
	})

	return module, nil
}

func randomSeed(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var starSeed starlark.Int

	if err := starlark.UnpackArgs(
		"seed",
		args, kwargs,
		"seed", &starSeed,
	); err != nil {
		return nil, fmt.Errorf("unpacking arguments for seed: %w", err)
	}

	seed, ok := starSeed.Int64()
	if !ok {
		return nil, fmt.Errorf("casting seed to int64")
	}

	rng, ok := thread.Local(threadRandKey).(*rand.Rand)
	if !ok || rng == nil {
		return nil, fmt.Errorf("RNG not set (very bad)")
	}

	rng.Seed(seed)

	return starlark.None, nil
}

func randomNumber(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var (
		starMin starlark.Int
		starMax starlark.Int
	)

	if err := starlark.UnpackArgs(
		"number",
		args, kwargs,
		"min", &starMin,
		"max", &starMax,
	); err != nil {
		return nil, fmt.Errorf("unpacking arguments for random number: %w", err)
	}

	min, ok := starMin.Int64()
	if !ok {
		return nil, fmt.Errorf("casting min to an int64")

	}

	max, ok := starMax.Int64()
	if !ok {
		return nil, fmt.Errorf("casting max to an int64")

	}

	if min < 0 {
		return nil, fmt.Errorf("min has to be 0 or greater")
	}

	if max < min {
		return nil, fmt.Errorf("max is less then min")
	}

	rng, ok := thread.Local(threadRandKey).(*rand.Rand)
	if !ok || rng == nil {
		return nil, fmt.Errorf("RNG not set (very bad!)")
	}

	return starlark.MakeInt64(rng.Int63n(max-min+1) + min), nil
}
