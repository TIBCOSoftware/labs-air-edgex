package driver

import (
	"fmt"
	"io/ioutil"
)

var currentSound = 0

// GetSoundData - get the data for specific sound
func GetSoundData(index int) []byte {

	filename := fmt.Sprintf("./data/0000000%d.json", index)
	// filename := fmt.Sprintf("./data/0000000%d.wav", index)

	fmt.Printf("Filename: %s \n", filename)

	dat, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println(err)
	}
	// fmt.Print(string(dat))

	return dat
}

// GetcurrentSoundIndex - returs the current index of the sound
func GetcurrentSoundIndex() int {
	value := currentSound

	currentSound++
	currentSound = currentSound % 9

	return value
}
