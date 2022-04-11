package wordpress

func (wp *WordPress) Create(wname string, port int32) error {
	err := CreateWordpressService(wname, port)
	if err == nil {
		CreateSecretKey(wname)
		err = CreateDatabasePvc(wname)
		if err != nil {
			return err
		}
		CreateWordpressPVC(wname)
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
		err = CreateWordpressDeployment(wname)
		if err != nil {
			return err
		}
		return nil
	}
	return err
}
