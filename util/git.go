package util

import "fmt"

func PublishChanges(fileName, message string) error {
	_, err := CallCmd("git", "pull")
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	_, err = CallCmd("git", "add", fileName)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	_, err = CallCmd("git", "commit", "-m", message)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	_, err = CallCmd("git", "push")
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	return nil
}
