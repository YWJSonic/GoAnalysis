package goanalysis

func (s *source) OnIdentifierList() []string {
	results := []string{}
	results = append(results, s.scanIdentifiers())
	if s.ch == ',' {
		results = append(results, s.scanIdentifiers())
	}

	return results
}
