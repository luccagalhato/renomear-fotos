package main

import (
	"fmt"
	"image/jpeg"
	"io/ioutil"
	"log"
	"os"
	"strings"
	//	"os"
)

const path = ("./1-Felipe")

func main() {

	sqlVar = &sqlStruct{}

	if err := sqlVar.init(); err != nil {
		log.Fatal(err)
	}
	defer sqlVar.Close()
	// var code = "011104030004"
	// var cor = "013"
	// filename := sqlVar.GetGtinCode(code, cor)

	dirroots, _ := ioutil.ReadDir(path)
	for _, rootFile := range dirroots {
		if !rootFile.IsDir() {
			continue
		}
		folders, _ := ioutil.ReadDir(path + rootFile.Name())
		CodeandColor := strings.Split(rootFile.Name(), "-")
		filenames, err := sqlVar.GetGtinCode(CodeandColor[0], CodeandColor[1])
		if err != nil {
			log.Println(err)
			continue
		}
		if len(filenames) == 0 {
			fmt.Println("não encontrado")
			continue
		}
		// if len(filenames) == 0 {
		// 	err = os.Mkdir("./nao encontrado/"+rootFile.Name(), 0755)
		// 	if err != nil {
		// 		log.Fatal(err)
		// 	}
		// 	for _, photoFile := range folders {
		// 		if photoFile.IsDir() || photoFile.Name() == "Thumbs.db" || photoFile.Name() == ".DS_Store" {
		// 			continue
		// 		}
		// 		file, err := os.Open(path + rootFile.Name() + "/" + photoFile.Name())
		// 		if err != nil {
		// 			log.Println(err)
		// 			continue
		// 		}
		// 		img, err := jpeg.Decode(file)
		// 		if err != nil {
		// 			log.Println(err)
		// 			continue
		// 		}
		// 		for _, name := range photoFile.Name() {
		// 			fmt.Println(name)
		// 			dest2file, err := os.OpenFile("./nao encontrado/"+rootFile.Name()+"/"+photoFile.Name(), os.O_WRONLY, 0766)
		// 			if err != nil {
		// 				log.Println(err)
		// 				continue
		// 			}
		// 			jpeg.Encode(dest2file, img, &jpeg.Options{})
		// 		}
		// 	}
		// }
		err = os.Mkdir("./encontrados/"+rootFile.Name(), 0755)
		if err != nil {
			log.Fatal(err)
		}
		for _, photoFile := range folders {
			if photoFile.IsDir() || photoFile.Name() == "Thumbs.db" || photoFile.Name() == ".DS_Store" {
				continue
			}
			// for _, gtin := range filenames {
			file, err := os.Open(path + rootFile.Name() + "/" + photoFile.Name())
			if err != nil {
				log.Println(err)
				continue
			}
			img, err := jpeg.Decode(file)
			if err != nil {
				log.Println(err)
				continue
			}

			for _, gtin := range filenames {

				destFile, err := os.OpenFile("./encontrados/"+rootFile.Name()+"/"+gtin+"-"+photoFile.Name(), os.O_WRONLY|os.O_CREATE, 0766)
				if err != nil {
					log.Println(err)
					continue
				}
				jpeg.Encode(destFile, img, &jpeg.Options{})
			}
			// }
			// format := strings.Split(photoFile.Name(), ".")[1]
			// newFileName := fmt.Sprintf("%s-%s-%d.%s", code, color, index, format)

		}
	}
}
