package main

import "fmt"

var graph = map[string][]string{
	"a": {"b", "l"},
	"b": {"l"},
	"c": {"a", "e", "i", "o"},
	"e": {},
	"g": {"h"},
	"h": {"e", "t"},
	"i": {"g", "l"},
	"l": {"e"},
	"n": {"c", "p", "t"},
	"o": {"g", "p"},
	"p": {"a", "r", "y"},
	"r": {"a", "e", "i"},
	"t": {"a", "b"},
	"u": {"n"},
	"y": {"g", "r", "t"},
}

type visitor struct {
	graph    map[string][]string
	visited  map[string]struct{}
	visiting map[string]struct{}
	cycle    []string
	sorted   []string
	offset   int
}

func (v *visitor) visit(n string) (cycleStart string) {
	if _, ok := v.visiting[n]; ok {
		if v.cycle == nil {
			v.cycle = make([]string, 0)
		}
		return n
	} else if _, ok := v.visited[n]; ok {
		return ""
	}
	v.visiting[n] = struct{}{}
	for _, m := range v.graph[n] {
		cycleStart := v.visit(m)
		if v.cycle != nil {
			if cycleStart != "" {
				if n == cycleStart {
					cycleStart = ""
				}
				v.cycle = append(v.cycle, n)
			}
			return cycleStart
		}
	}
	delete(v.visiting, n)
	v.visited[n] = struct{}{}
	v.offset--
	v.sorted[v.offset] = n
	return ""
}

func toposort(g map[string][]string) (sorted, cycle []string) {
	v := &visitor{
		graph:    g,
		offset:   len(g),
		sorted:   make([]string, len(g)),
		visited:  map[string]struct{}{},
		visiting: map[string]struct{}{},
	}
	for n := range g {
		if _, ok := v.visited[n]; ok {
			continue
		}
		v.visit(n)
		if v.cycle != nil {
			return nil, v.cycle
		}
	}
	return v.sorted, nil
}

func main() {
	sorted, cycle := toposort(graph)
	fmt.Println("cycle", cycle)
	fmt.Println("sorted", sorted)

	fmt.Println("--")

	graph["e"] = []string{"u"}
	sorted, cycle = toposort(graph)
	fmt.Println("cycle", cycle)
	fmt.Println("sorted", sorted)
}
