package joomla

func (jo *Joomla) Create(wname string, port int32) error {
	err := CreateJoomlaService(wname, port)
	if err == nil {
		CreateSecretKey(wname)
		err = CreateDatabasePvc(wname)
		if err != nil {
			return err
		}
		err = CreateJoomlaPVC(wname)
		if err != nil {
			return err
		}
		err = CreateDatabaseService(wname)
		if err != nil {
			return err
		}
		err = CreateDatabaseDeployment(wname)
		if err != nil {
			return err
		}
		err = CreateJoomlaDeployment(wname)
		if err != nil {
			return err
		}
		return nil
	}
	return err
}
