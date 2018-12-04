package main

func GetPrimary(view int, nservers int) int {
	return view % nservers
}
