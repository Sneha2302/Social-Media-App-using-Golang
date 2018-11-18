package main


//Returns users password
func getPassword(usrname string) (bool, string){

	user, ok := userdata[usrname]
	if(!ok){
		debugPrint("No such user")
		return false, "No such User"
	}
	return true, user.password

}