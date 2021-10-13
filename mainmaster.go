package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	//	"os"
)

const path = "./renomear-auto"

func main() {

	sqlVar = &sqlStruct{}

	if err := sqlVar.init(); err != nil {
		log.Fatal(err)
	}
	defer sqlVar.Close()
	// var code = "011104030004"
	// var cor = "013"
	// filename := sqlVar.GetGtinCode(code, cor)

	if _, err := os.Stat("./encontrados/"); os.IsNotExist(err) {
		err := os.Mkdir("./encontrados/", 0755)
		if err != nil {
			log.Fatal(err)
		}
	}

	if _, err := os.Stat("./nao encontrados/"); os.IsNotExist(err) {
		err := os.Mkdir("./nao encontrados/", 0755)
		if err != nil {
			log.Fatal(err)
		}
	}

	dirroots, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}
	for ind, rootFile := range dirroots {
		if !rootFile.IsDir() {
			continue
		}
		folders, err := ioutil.ReadDir(path + "/" + rootFile.Name())
		if err != nil {
			log.Fatal(err)
		}
		CodeandColor := strings.Split(rootFile.Name(), "-")
		if len(CodeandColor) < 2 {
			dir(path+"/"+rootFile.Name(), "./nao encontrados/"+rootFile.Name())

			if err2 := ioutil.WriteFile("./nao encontrados/"+rootFile.Name()+"/error.txt", []byte("sem cor definida"), 0766); err2 != nil {
				log.Println(err2)
			}
			continue
		}

		filenames, err := sqlVar.GetGtinCode(CodeandColor[0], CodeandColor[1])
		if err != nil {
			dir(path+"/"+rootFile.Name(), "./nao encontrados/"+rootFile.Name())

			if err2 := ioutil.WriteFile("./nao encontrados/"+rootFile.Name()+"/error.txt", []byte(err.Error()), 0766); err2 != nil {
				log.Println(err2)
			}
			continue
		}

		err = os.Mkdir("./encontrados/"+rootFile.Name(), 0755)
		if err != nil {
			log.Println(err)
			continue
		}
		for index, photoFile := range folders {

			pfS := strings.Split(photoFile.Name(), ".")
			term := pfS[len(pfS)-1]

			if photoFile.IsDir() || photoFile.Name() == "Thumbs.db" || photoFile.Name() == ".DS_Store" {
				continue
			}

			file, err := os.Open(path + "/" + rootFile.Name() + "/" + photoFile.Name())
			if err != nil {
				log.Println("error opening photo", err)
				continue
			}
			filebytes, _ := ioutil.ReadAll(file)
			// img, err := jpeg.Decode(file)
			// if err != nil {
			// 	log.Println("error decoding photo", err)
			// 	continue
			// }

			for _, gtin := range filenames {
				destPath := fmt.Sprintf("%s/%s/%s-%d.%s", "./encontrados", rootFile.Name(), gtin, index+1, term)
				destFile, err := os.OpenFile(destPath, os.O_WRONLY|os.O_CREATE, 0666)
				if err != nil {
					log.Println("error opening destination file", destPath, err)
					continue
				}
				// if err := jpeg.Encode(destFile, img, &jpeg.Options{Quality: 100}); err != nil {
				// 	log.Println("error encoding destination file", destPath, err)
				// }
				destFile.Write(filebytes)
				// ioutil.Write(destPath, filebytes, 0666)
			}
		}
		fmt.Printf("%d/%d\n", ind+1, len(dirroots))
	}
}

func dir(src string, dst string) error {
	var err error
	var fds []os.FileInfo
	var srcinfo os.FileInfo

	if srcinfo, err = os.Stat(src); err != nil {
		return err
	}

	if err = os.MkdirAll(dst, srcinfo.Mode()); err != nil {
		return err
	}

	if fds, err = ioutil.ReadDir(src); err != nil {
		return err
	}
	for _, fd := range fds {
		srcfp := strings.Join([]string{src, fd.Name()}, "/")
		dstfp := strings.Join([]string{dst, fd.Name()}, "/")

		if fd.IsDir() {
			if err = dir(srcfp, dstfp); err != nil {
				fmt.Println(err)
			}
		} else {
			if err = file(srcfp, dstfp); err != nil {
				fmt.Println(err)
			}
		}
	}
	return nil
}

func file(src, dst string) error {
	var err error
	var srcfd *os.File
	var dstfd *os.File
	var srcinfo os.FileInfo

	if srcfd, err = os.Open(src); err != nil {
		return err
	}
	defer srcfd.Close()

	if dstfd, err = os.Create(dst); err != nil {
		return err
	}
	defer dstfd.Close()

	if _, err = io.Copy(dstfd, srcfd); err != nil {
		return err
	}
	if srcinfo, err = os.Stat(src); err != nil {
		return err
	}
	return os.Chmod(dst, srcinfo.Mode())
}
