package main

var tmpDirCreated string

func TempDir() (string, error) {
	if tmpDirCreated != "" {
		return tmpDirCreated, nil
	}
	created, err := CreateTempDir()
	if err != nil {
		return "", err
	}
	tmpDirCreated = created
	return tmpDirCreated, nil
}
