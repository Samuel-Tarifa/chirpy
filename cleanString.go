package main

import "strings"


func cleanString(s string) string {
	badWords:=badWordsMap{
		"kerfuffle":true,
		"sharbert":true,
		"fornax":true,
	}
	arr:=strings.Split(s, " ")
	for i,item:= range arr{
		if _,ok:=badWords[strings.ToLower(item)];ok{
			arr[i]="****"
		}
	}
	return strings.Join(arr, " ")
}