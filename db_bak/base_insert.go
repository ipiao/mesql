package mesql

func (this *BaseService) InsertModels(models interface{}) (int, error) {
	return this.DB.InsertModels(models)
}

func (this *BaseService) InsertModels2(table string, models interface{}) (int, error) {
	return this.DB.InsertModels2(table, models)
}
