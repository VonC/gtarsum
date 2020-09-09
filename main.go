package main

import (
	"archive/tar"
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strings"
	"sync"

	"github.com/VonC/gtarsum/version"

	"github.com/vbauerster/mpb/v5"

	"github.com/vbauerster/mpb/v5/decor"
)

func main() {
	l := len(os.Args) - 1
	if l < 1 {
		log.Fatalf("At least one filename is expected (instead of %d)", l)
	}

	var wg sync.WaitGroup
	envp := os.Getenv("progress")
	var p *mpb.Progress
	if envp != "" {
		p = mpb.New(mpb.WithWaitGroup(&wg))
	}

	for _, f := range os.Args[1:] {
		fl := strings.ToLower(f)
		if fl == "-v" || fl == "--version" || fl == "version" {
			fmt.Println(version.String())
			os.Exit(0)
		}
	}
	hbs := newHashables(os.Args[1:], p, envp)
	res := hbs.hash(&wg)
	if res.isPrintable() {
		fmt.Println(res)
	}
	os.Exit(res.status)
}

func (hbs hashables) hash(wg *sync.WaitGroup) *result {

	l := hbs.len()
	results := make(chan string, l)
	errors := make(chan error, l)
	for _, hb := range hbs.hbs {
		wg.Add(1)
		go func(hb *hashable) {
			defer wg.Done()
			h1h := hb.hash()
			results <- h1h
		}(hb)
	}

	wg.Wait()
	close(results)
	close(errors)

	if hbs.p != nil {
		hbs.p.Wait()
	}

	for err := range errors {
		println(err.Error())
		os.Exit(1)
	}

	currentHash := ""
	status := 0
	var differ bool
	i := 0
	for res := range results {
		i++
		differ = false
		if currentHash == "" {
			currentHash = res
		} else if currentHash != res {
			status = 1
			differ = true
		}
		if hbs.p != nil {
			if strings.HasSuffix(hbs.envp, ".hash") {
				fe := fmt.Sprintf("%s%d", hbs.envp, i)
				write(fe, res)
			}
			if differ {
				fmt.Printf("File '%s' hash '%s' differs\n", os.Args[i], res)
			} else {
				fmt.Printf("File '%s' hash='%s'\n", os.Args[i], res)
			}
		}
	}

	res := &result{hash: currentHash, status: status}
	return res
}

func write(fe, res string) {
	f, err := os.Create(fe)
	check(err)
	defer f.Close()
	_, err = f.WriteString(res)
	check(err)
}

type hashable struct {
	f       string
	entries entries
	p       *mpb.Progress
}

type hashables struct {
	hbs  []*hashable
	p    *mpb.Progress
	envp string
}

func newHashables(fs []string, p *mpb.Progress, envp string) *hashables {
	hbs := make([]*hashable, 0)
	for _, f := range os.Args[1:] {
		f = strings.Trim(f, `"`)
		if f == "" {
			log.Fatalf(`At least one filename is expected (instead of empty "" filename)`)
		}

		hb := newHashable(f, p)
		hbs = append(hbs, hb)
	}
	return &hashables{hbs: hbs, p: p, envp: envp}
}

func (hbs *hashables) len() int {
	return len(hbs.hbs)
}

type result struct {
	hash   string
	status int
}

func (r *result) isPrintable() bool {
	return r.status == 0
}

func (r *result) String() string {
	return r.hash
}

func newHashable(f string, p *mpb.Progress) *hashable {
	return &hashable{f: f, entries: make(map[string]string), p: p}
}

func check(err error) {
	if err != nil {
		panic(err)
	}
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

func (hb *hashable) hash() string {
	hbBntries := hb.gtarsum()
	return hbBntries.hash()
}

func (hb *hashable) gtarsum() entries {

	nbFiles := 0
	fnbFiles := func(tr *tar.Reader, th *tar.Header) {
		nbFiles = nbFiles + 1
	}
	readTarFiles(hb.f, fnbFiles)

	// fmt.Printf("%d files to process in '%s'\n", nbFiles, filename)

	entries := make(map[string]string)

	var bar *mpb.Bar
	if hb.p != nil {
		bar = hb.p.AddBar(int64(nbFiles), nil,
			mpb.PrependDecorators(
				decor.Name(fmt.Sprintf("File '%s' (%d): ", hb.f, nbFiles)),
				decor.NewPercentage("%d"),
			),
		)
	}

	fHashEntry := func(tr *tar.Reader, th *tar.Header) {
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
		if bar != nil {
			bar.Increment()
		}
	}

	readTarFiles(hb.f, fHashEntry)

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
