package collector

import (
	"context"
	"encoding/json"
	"testing"
)

// stubCollector 是测试用的抓取器；不做任何 IO。
type stubCollector struct {
	kind string
}

func (s *stubCollector) Kind() string { return s.kind }
func (s *stubCollector) Fetch(context.Context, json.RawMessage) ([]Item, error) {
	return nil, nil
}

func TestRegistryRoundtrip(t *testing.T) {
	// 隔离全局 registry：保存 → 清空 → 恢复
	registryMu.Lock()
	saved := registry
	registry = map[string]Collector{}
	registryMu.Unlock()
	defer func() {
		registryMu.Lock()
		registry = saved
		registryMu.Unlock()
	}()

	Register(&stubCollector{kind: "test_a"})
	Register(&stubCollector{kind: "test_b"})

	if c, ok := Get("test_a"); !ok || c.Kind() != "test_a" {
		t.Fatalf("Get(test_a) = %v, %v", c, ok)
	}
	if _, ok := Get("missing"); ok {
		t.Fatalf("Get(missing) should return false")
	}

	kinds := Kinds()
	if len(kinds) != 2 {
		t.Fatalf("Kinds() len = %d, want 2", len(kinds))
	}
}

func TestRegisterOverwrites(t *testing.T) {
	registryMu.Lock()
	saved := registry
	registry = map[string]Collector{}
	registryMu.Unlock()
	defer func() {
		registryMu.Lock()
		registry = saved
		registryMu.Unlock()
	}()

	Register(&stubCollector{kind: "dup"})
	c1, _ := Get("dup")
	Register(&stubCollector{kind: "dup"})
	c2, _ := Get("dup")
	if c1 == c2 {
		t.Fatalf("second Register should replace, got same pointer")
	}
}
