package example

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	jsr "github.com/teamdunno/go-jsr-registry"
	jsr_tools "github.com/teamdunno/go-jsr-registry/tools"
)

func runProgram() error {
	// create new client
	client, err := jsr.NewClient()
	// if it fails, return error
	if err != nil {
		return err
	}
	// get the package meta
	pkg, err := client.GetPackageMeta(jsr.PackageMetaOption{Scope: "dunno", Name: "object"})
	// if it fails, return error
	if err != nil {
		return errors.New("failed to get package")
	}
	// stringify all versions
	resultBefore, err := json.MarshalIndent(pkg.Versions, "  ", "    ")
	// if it fails, return error
	if err != nil {
		return errors.New("failed to get package")
	}
	// print stringified all versions
	fmt.Print("before: " + string(resultBefore))
	// stringify unyanked versions
	resultAfter, err := json.MarshalIndent(jsr_tools.GetUnyankedVersionsFromPackageMeta(pkg.Versions), "  ", "    ")
	// if it fails, return error
	if err != nil {
		return errors.New("failed to get package")
	}
	// print stringified unyanked versions
	fmt.Print("after: " + string(resultAfter))
	// return nil since we dosent have any error again
	return nil
}

func main() {
	// run the program
	var res = runProgram()
	// if it has an error, print to the console and exit as fatal
	// or else, just exit normally
	if res != nil {
		log.Fatal(res)
	}
}
