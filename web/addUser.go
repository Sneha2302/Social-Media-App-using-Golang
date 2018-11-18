package main

//function to add user to data on registration
func addUser(usrname string, pwd string) int  {
	_, ok := userdata[usrname]
	if(ok){
		debugPrint("Debug: User already exists")
		return 0
	}
	usr := User{username:usrname,password:pwd}
	usr.follows = make(map[string]bool)
	userdata[usrname] = usr
	debugPrint("Debug: User added")
	return 1
}