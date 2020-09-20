package database

func Setup() {
	var db = new(Mysql)
	db.Setup()
}
