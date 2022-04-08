package joomla

func (oc *Joomla) Create(wname string, port int32) error {
	err := CreateJoomlaService(wname, port)
	if err == nil {
		CreateSecretKey()
		CreateDatabasePvc(wname)
		CreateJoomlaPVC(wname)
		CreateDatabaseService(wname)
		CreateDatabaseDeployment(wname)
		CreateJoomlaDeployment(wname)
		return nil
	}
	return err
}
