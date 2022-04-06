package wordpress

func (wp *WordPress) Create(wname string, port int32) error {
	err := CreateWordpressService(wname, port)
	if err == nil {
		CreateSecretKey()
		CreateDatabasePvc(wname)
		CreateWordpressPVC(wname)
		CreateDatabaseService(wname)
		CreateDatabaseDeployment(wname)
		CreateWordpressDeployment(wname)
		return nil
	}
	return err
}
