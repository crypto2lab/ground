package runtime

import (
	"context"
	"fmt"

	"github.com/crypto2lab/ground/internal/chainspec"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

type Instance struct {
	ctx           context.Context
	runtime       wazero.Runtime
	runtimeModule api.Module
	runtimeBytes  []byte
}

func NewInstance(ctx context.Context, runtimeBytes []byte) *Instance {
	return &Instance{
		ctx:          ctx,
		runtime:      wazero.NewRuntime(ctx),
		runtimeBytes: runtimeBytes,
	}
}

func (instance *Instance) Instantiate() error {
	wasi_snapshot_preview1.MustInstantiate(instance.ctx, instance.runtime)
	mod, err := instance.runtime.Instantiate(instance.ctx, instance.runtimeBytes)
	if err != nil {
		return fmt.Errorf("instantiating runtime: %w", err)
	}

	instance.runtimeModule = mod
	return nil
}

func (instance *Instance) Initialize(genesis chainspec.Genesis) error {
	return nil
}

func (instance *Instance) Close() error {
	return instance.runtime.Close(instance.ctx)
}
