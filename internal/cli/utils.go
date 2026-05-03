package cli

func normalizeNamespace(ns string) string {
	if ns == "" {
		return "default"
	}
	return ns
}
