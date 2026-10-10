package ring

import "testing"

func TestStableOwnership(t *testing.T) {
 a := New([]string{"a", "b", "c"})
 b := New([]string{"c", "a", "b"})
 for _, key := range []string{"alice", "bob", "carol", "dave"} {
  if a.Owner(key) != b.Owner(key) { t.Fatalf("inconsistent owner for %s", key) }
 }
}
