package parser

type Names struct {
	nameToID map[string]int
	names    []string
}

func (n *Names) IDFromName(name string) int {
	if n.nameToID == nil {
		n.nameToID = make(map[string]int, 100)
	}
	id, ok := n.nameToID[name]
	if !ok {
		n.names = append(n.names, name)
		id = len(n.names) - 1
		n.nameToID[name] = id
	}
	return id + 1
}

func (n *Names) NameFromID(id int) string {
	return n.names[id-1]
}
