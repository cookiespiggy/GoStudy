package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/cavaliercoder/go-rpm"
	"github.com/sassoftware/go-rpmutils"
)

func ExampleMD5Check(rpm_file string) {
	// open a rpm package for reading
	f, err := os.Open(rpm_file)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	// validate md5 checksum
	if err := rpm.MD5Check(f); err == nil {
		fmt.Printf("Package passed checksum validation\n")
	} else if err == rpm.ErrMD5ValidationFailed {
		fmt.Printf("Package failed checksum validation\n")
	} else {
		panic(err)
	}

	// Output: Package passed checksum validation
}

func ExamplePackageFile_Files(rpm_file string) {
	// open a package file
	pkg, err := rpm.OpenPackageFile(rpm_file)
	if err != nil {
		panic(err)
	}

	// list each file
	files := pkg.Files()
	fmt.Printf("total %v\n", len(files))
	for _, fi := range files {
		fmt.Printf("%v %v %v %5v %v %v %v\n",
			fi.Mode().Perm(),
			fi.Owner(),
			fi.Group(),
			fi.Size(),
			fi.ModTime().UTC().Format("Jan 02 15:04"),
			fi.Name(),
			fi.Digest())
	}
}

func main() {
	rpm_file := "net-tools-1.60-78.el5.x86_64.rpm"
	ExampleMD5Check(rpm_file)
	ExamplePackageFile_Files(rpm_file)

	f, err := os.Open(rpm_file)
	if err != nil {
		panic(err)
	}
	rpm, err := rpmutils.ReadRpm(f)
	if err != nil {
		panic(err)
	}
	// // Getting metadata
	// nevra, err := rpm.Header.GetNEVRA()
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(nevra)
	// provides, err := rpm.Header.GetStrings(rpmutils.PROVIDENAME)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println("Provides:")
	// for _, p := range provides {
	// 	fmt.Println(p)
	// }
	// Extracting payload
	// if err := rpm.ExpandPayload("destdir"); err != nil {
	// 	panic(err)
	// }
	reader, err := rpm.PayloadReader()
	if err != nil {
		panic(err)
	}

	for {
		header, err := reader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}

		fmt.Println(header.Filename(), header.Filesize64())
		if header.Mode()&0111 != 0 && header.Filesize64() != 0 {
			if strings.HasPrefix(header.Filename(), "/usr/share/doc") {
				continue
			}
			buf := make([]byte, header.Filesize64())
			_, err := reader.Read(buf)
			if err != nil {
				panic(err)
			}
			has := md5.Sum(buf)
			md5str := fmt.Sprintf("%x", has)
			fmt.Println(header.Filename(), md5str)
		}
	}
}
