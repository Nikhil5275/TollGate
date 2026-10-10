package ring

import (
 "fmt"
 "hash/fnv"
 "sort"
)

type point struct { hash uint64; node string }
type Ring struct { points []point }
func hash(key string) uint64 { h := fnv.New64a(); _, _ = h.Write([]byte(key)); return h.Sum64() }

// New spreads each server across 128 virtual nodes.
func New(nodes []string) *Ring {
 r := &Ring{}
 for _, node := range nodes {
  for i := 0; i < 128; i++ { r.points = append(r.points, point{hash(fmt.Sprintf("%s#%d", node, i)), node}) }
 }
 sort.Slice(r.points, func(i, j int) bool { return r.points[i].hash < r.points[j].hash })
 return r
}
func (r *Ring) Owner(key string) string {
 if len(r.points) == 0 { return "" }
 h := hash(key)
 i := sort.Search(len(r.points), func(i int) bool { return r.points[i].hash >= h })
 if i == len(r.points) { i = 0 }
 return r.points[i].node
}
