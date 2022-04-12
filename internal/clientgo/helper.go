package clientgo

type ListNames struct {
	Deployment []string
	Service    []string
	Pod        []string
	Pvc        []string
}

func GetNamespace(app string, name string) string {
	if app == "wordpress" {
		return name + "-wp"
	} else if app == "joomla" {
		return name + "-joomla"
	} else {
		return ""
	}
}
