package memory

import (
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/ebitengine/purego"
	"github.com/smtdfc/nagare/core/logger"
	"github.com/smtdfc/nagare/pkgs/paths"
)

func getVectorModule() string {
	var filename string

	switch runtime.GOOS {
	case "darwin":
		filename = "libnagare_vector.dylib"
	case "windows":
		filename = "nagare_vector.dll"
	case "linux", "freebsd", "netbsd":
		filename = "libnagare_vector.so"
	default:
		panic(fmt.Errorf("GOOS=%s is not supported", runtime.GOOS))
	}

	return filepath.Join(paths.ModulesDir, filename)
}

type VectorIndex uintptr

type VectorMemory struct {
	lib uintptr

	vectorIndexCreate func(dim uintptr, bits uintptr) uintptr
	vectorIndexLoad   func(path string) uintptr
	// vectorIndexAdd     func(index uintptr, data uintptr, ids uintptr, count uintptr, dim uintptr) int32
	vectorIndexSync      func(index uintptr, path string) int32
	vectorIndexLastError func() string

	// vectorIndexDestroy func(index uintptr)
	index  VectorIndex
	logger *logger.BaseLogger
}

func (v *VectorMemory) Create(dim, bits uintptr) error {
	index := v.vectorIndexCreate(dim, bits)

	if index == 0 {
		err := v.vectorIndexLastError()

		if err == "" {
			return fmt.Errorf("failed to create vector index")
		}

		return fmt.Errorf("failed to create vector index: %s", err)
	}

	v.index = VectorIndex(index)
	return nil
}

func (v *VectorMemory) Load(path string) error {
	index := v.vectorIndexLoad(path)

	if index == 0 {
		return fmt.Errorf("failed to load vector index")
	}

	v.index = VectorIndex(index)
	return nil
}

func (v *VectorMemory) Sync(path string) error {
	code := v.vectorIndexSync(uintptr(v.index), path)

	if code != 0 {
		err := v.vectorIndexLastError()
		if err != "" {
			return fmt.Errorf("failed to sync vector index: %s", err)
		}
		return fmt.Errorf("failed to sync vector index: code=%d", code)
	}

	return nil
}

// @Injectable
func NewVectorMemory(logger *logger.BaseLogger) (*VectorMemory, error) {
	logger = logger.With("module", "vector-memory")

	lib, err := openLibrary(getVectorModule())
	if err != nil {
		logger.Error("Failed to open vector module", "error", err)
		return nil, err
	}

	vm := &VectorMemory{
		lib:    lib,
		logger: logger,
	}

	purego.RegisterLibFunc(
		&vm.vectorIndexCreate,
		lib,
		"vector_index_create",
	)

	purego.RegisterLibFunc(
		&vm.vectorIndexLoad,
		lib,
		"vector_index_load",
	)

	// purego.RegisterLibFunc(
	// 	&vm.vectorIndexAdd,
	// 	lib,
	// 	"vector_index_add",
	// )

	purego.RegisterLibFunc(
		&vm.vectorIndexSync,
		lib,
		"vector_index_sync",
	)

	purego.RegisterLibFunc(
		&vm.vectorIndexLastError,
		lib,
		"vector_index_last_error",
	)

	// purego.RegisterLibFunc(
	// 	&vm.vectorIndexDestroy,
	// 	lib,
	// 	"vector_index_destroy",
	// )

	return vm, nil
}
