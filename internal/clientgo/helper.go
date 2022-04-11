package clientgo

type ListNames struct {
	Deployment []string
	Service    []string
	Pod        []string
	Pvc        []string
}
