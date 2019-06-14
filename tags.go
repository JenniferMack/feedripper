package feedpub

// sort tags by priority
func (t tags) Len() int           { return len(t) }
func (t tags) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }
func (t tags) Less(i, j int) bool { return t[i].Priority < t[j].Priority }

func (t tags) priOutOfRange() bool {
	idx := uint(0)
	cmp := uint(len(t) - 1)

	for _, v := range t {
		if v.Priority > idx {
			idx = v.Priority
		}
	}
	return idx > cmp
}
