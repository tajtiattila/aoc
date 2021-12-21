package aoc

import (
	"flag"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/tajtiattila/aoc/input"
)

var Verbose bool

func Main(year int) {
	input.Year = year

	flag.BoolVar(&Verbose, "v", false, "verbose mode")
	clearCache := flag.Bool("cc", false, "clear cache")
	flag.BoolVar(&input.IgnoreCache, "ic", false, "ignore cache")
	flag.Parse()

	if *clearCache {
		if err := input.ClearCache(); err != nil {
			log.Fatalln("clear cache:", err)
		}
	}

	want := make(map[int]struct{})
	for _, a := range flag.Args() {
		if ok := parsearg(want, a); !ok {
			log.Fatalf("invalid argument %q", a)
		}
	}

	for i, f := range regfn {
		if f == nil {
			continue
		}
		if _, ok := want[i]; ok || len(want) == 0 {
			f()
		}
	}
}

func parsearg(m map[int]struct{}, arg string) bool {
	var a, b int
	var aerr, berr error
	if strings.HasSuffix(arg, "+") {
		a, aerr = strconv.Atoi(strings.TrimSuffix(arg, "+"))
		b = len(regfn)
	} else if i := strings.Index(arg, ".."); i >= 0 {
		a, aerr = strconv.Atoi(arg[:i])
		b, berr = strconv.Atoi(arg[i+2:])
	} else {
		a, aerr = strconv.Atoi(arg)
		b = a
	}
	if aerr != nil || berr != nil {
		return false
	}
	for i := a; i <= b; i++ {
		m[i] = struct{}{}
	}
	return true
}

func Log(args ...interface{}) {
	if Verbose {
		fmt.Print(args...)
	}
}

func Logln(args ...interface{}) {
	if Verbose {
		fmt.Println(args...)
	}
}

func Logf(format string, args ...interface{}) {
	if Verbose {
		fmt.Printf(format, args...)
	}
}

var regfn []func()

func Register(day int, fn func()) {
	if day >= len(regfn) {
		x := make([]func(), day+1)
		copy(x, regfn)
		regfn = x
	}

	if regfn[day] != nil {
		panic("Register: already in use")
	}

	regfn[day] = fn
}
