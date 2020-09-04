package main

import (
	"archive/tar"
	"crypto/sha256"
	"fmt"
	"gtarsum/version"
	"io"
	"log"
	"os"
	"sort"
	"strings"
)

func main() {
	l := len(os.Args) - 1
	if l < 1 {
		log.Fatalf("At least one filename is expected (instead of %d)", l)
	}
	f := os.Args[1]
	f = strings.Trim(f, `"`)
	if f == "" {
		log.Fatalf(`One, and only one filename is expected (instead of empty "" filename)`)
	}
	fl := strings.ToLower(f)
	if fl == "-v" || fl == "--version" || fl == "version" {
		fmt.Println(version.String())
		os.Exit(0)
	}
	fmt.Printf("Tarsum for file '%s'\n", f)
	h1 := gtarsum(f)
	fmt.Printf("File '%s' hash='%s'\n", f, h1.hash())
}

type entries map[string]string

type tarFileVisitor func(tr *tar.Reader, th *tar.Header)

func readTarFiles(filename string, tfv tarFileVisitor) {

	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}

	defer func() {
		_ = f.Close()
	}()

	tr := tar.NewReader(f)

	for {
		// get the next file entry
		hd, err := tr.Next()

		if err != nil {

			if err != io.EOF {
				panic(err)
			}

			// fmt.Println("tar EOF")
			break
		}

		switch hd.Typeflag {
		case tar.TypeReg: // = regular file
			tfv(tr, hd)
		}
	}

}

func gtarsum(filename string) entries {

	nbFiles := 0
	f := func(tr *tar.Reader, th *tar.Header) {
		nbFiles = nbFiles + 1
	}
	readTarFiles(filename, f)

	fmt.Printf("%d files to process in '%s'\n", nbFiles, filename)

	entries := make(map[string]string)

	f = func(tr *tar.Reader, th *tar.Header) {
		name := th.Name
		h := sha256.New()
		for {
			buf := make([]byte, 1024*1024)

			bytesRead, err := tr.Read(buf)
			if err != nil {
				if err != io.EOF {
					panic(err)
				}
			}

			if bytesRead > 0 {
				_, err := h.Write(buf[:bytesRead])
				if err != nil {
					panic(err)
				}
			}

			if err == io.EOF {
				//fmt.Printf("tar entry '%s' EOF\n", name)
				break
			}
		}
		bs := h.Sum(nil)
		entries[name] = fmt.Sprintf("%x", bs)
	}

	readTarFiles(filename, f)

	return entries
}

func (es entries) hash() string {
	names := []string{}
	for name := range es {
		names = append(names, name)
	}
	sort.Strings(names)
	h := sha256.New()
	for _, name := range names {
		bs := es[name]
		_, err := h.Write([]byte(bs))
		if err != nil {
			panic(err)
		}
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}
