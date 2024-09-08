package db

// 全局私有单例
var idGenerator = newIDGenerator()
func GenId() uint64{
	return idGenerator.generateID()
}

func Execute(sql string){

}
