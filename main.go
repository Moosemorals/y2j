package main

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"gopkg.in/yaml.v2"
)

func convert(i interface{}) interface{} {
	switch x := i.(type) {
	case map[interface{}]interface{}:
		m2 := map[string]interface{}{}
		for k, v := range x {
			switch y := k.(type) {
			case string:
				m2[y] = convert(v)
			case bool:
				m2[strconv.FormatBool(y)] = convert(v)
			case int:
				m2[strconv.Itoa(y)] = convert(v)
			default:
				log.Printf("Warn: Werid type %v", k)
			}
		}
		return m2
	case []interface{}:
		for i, v := range x {
			x[i] = convert(v)
		}
	}
	return i
}

func changeExtension(path, newExt string) string {
	dir, file := filepath.Split(path)
	return filepath.Join(dir, strings.TrimSuffix(file, filepath.Ext(file))+"."+newExt)
}

func convertFile(pathIn, pathOut string) error {
	inFile, err := os.Open(pathIn)
	if err != nil {
		return err
	}
	defer inFile.Close()

	var y interface{}
	if err = yaml.NewDecoder(inFile).Decode(&y); err != nil {
		return err
	}

	if err = os.MkdirAll(filepath.Dir(pathOut), 0755); err != nil {
		return err
	}
	outFile, err := os.Create(pathOut)
	if err != nil {
		return err
	}
	defer outFile.Close()
	if err = json.NewEncoder(outFile).Encode(convert(y)); err != nil {
		return err
	}
	return nil
}

func convertTree(base string) {
	filepath.Walk(base, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		newPath := filepath.Join("json", strings.TrimPrefix(path, base))
		if filepath.Ext(path) == ".yaml" {
			newPath = changeExtension(newPath, "json")
		} else {
			newPath += ".json"
		}

		log.Printf("%s => %s", path, newPath)

		if e := convertFile(path, newPath); e != nil {
			log.Print(err)
			return e
		}
		return nil
	})
}

func main() {
	/*
		if err := convertFile("sde/bsd/warCombatZones.yaml", "json"); err != nil {
			log.Fatal(err)
		}
	*/

	convertTree("sde")
}
