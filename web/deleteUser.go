package main


//Delete a user account
func deleteUser(username string) int  {
	//TODO: for later stages, we'll have to add Locks here
	debugPrint("Deleting User: " + username +"Account")
	delete(userdata,username)
	return 1
}
