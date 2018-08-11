package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/pascaldekloe/colfer"
)

// ANSI escape codes for markup
const (
	bold      = "\x1b[1m"
	italic    = "\x1b[3m"
	underline = "\x1b[4m"
	clear     = "\x1b[0m"
)

var (
	basedir = flag.String("b", ".", "Use a specific destination base `directory`.")
	prefix  = flag.String("p", "", "Adds a package `prefix`. Use slash as a separator when nesting.")
	format  = flag.Bool("f", false, "Normalizes the format of all input schemas on the fly.")
	verbose = flag.Bool("v", false, "Enables verbose reporting to "+italic+"standard error"+clear+".")

	sizeMax = flag.String("s", "16 * 1024 * 1024", "Sets the default upper limit for serial byte sizes. The\n`expression` is applied to the target language under the name\nColferSizeMax.")
	listMax = flag.String("l", "64 * 1024", "Sets the default upper limit for the number of elements in a\nlist. The `expression` is applied to the target language under\nthe name ColferListMax.")

	superClass = flag.String("x", "", "Makes all generated classes extend a super `class`. Use slash as\na package separator. Java only.")
)

var report = log.New(ioutil.Discard, "", 0)

func main() {
	flag.Parse()

	log.SetFlags(0)
	if *verbose {
		report.SetOutput(os.Stderr)
	}

	var files []string
	switch args := flag.Args(); len(args) {
	case 0:
		flag.Usage()
		os.Exit(2)
	case 1:
		files = []string{"."}
	default:
		files = args[1:]
	}

	// select language
	var gen func(string, colfer.Packages) error
	switch lang := flag.Arg(0); strings.ToLower(lang) {
	case "c":
		report.Println("Set up for C")
		gen = colfer.GenerateC
		if *superClass != "" {
			log.Fatal("colf: super class not supported with C")
		}

	case "go":
		report.Println("Set up for Go")
		gen = colfer.GenerateGo
		if *superClass != "" {
			log.Fatal("colf: super class not supported with Go")
		}

	case "java":
		report.Println("Set up for Java")
		gen = colfer.GenerateJava

	case "javascript", "js", "ecmascript":
		report.Println("Set up for ECMAScript")
		gen = colfer.GenerateECMA
		if *superClass != "" {
			log.Fatal("colf: super class not supported with ECMAScript")
		}

	default:
		log.Fatalf("colf: unsupported language %q", lang)
	}

	// resolve clean file set
	var writeIndex int
	for i := 0; i < len(files); i++ {
		f := files[i]

		info, err := os.Stat(f)
		if err != nil {
			log.Fatal(err)
		}
		if info.IsDir() {
			colfFiles, err := filepath.Glob(filepath.Join(f, "*.colf"))
			if err != nil {
				log.Fatal(err)
			}
			files = append(files, colfFiles...)
			continue
		}

		f = filepath.Clean(f)
		for j := 0; ; j++ {
			if j == writeIndex {
				files[writeIndex] = f
				writeIndex++
				break
			}
			if files[j] == f {
				report.Println("Duplicate inclusion of", f, "ignored")
				break
			}
		}
	}
	files = files[:writeIndex]
	report.Println("Found schema files", strings.Join(files, ", "))

	packages, err := colfer.ParseFiles(files)
	if err != nil {
		log.Fatal(err)
	}

	if *format {
		for _, file := range files {
			changed, err := colfer.Format(file)
			if err != nil {
				log.Fatal(err)
			}
			if changed {
				log.Println("colf: formatted", file)
			}
		}
	}

	if len(packages) == 0 {
		log.Fatal("colf: no struct definitons found")
	}

	for _, p := range packages {
		p.Name = path.Join(*prefix, p.Name)
		p.SizeMax = *sizeMax
		p.ListMax = *listMax
		p.SuperClass = *superClass
	}

	if err := gen(*basedir, packages); err != nil {
		log.Fatal(err)
	}
}

func init() {
	cmd := os.Args[0]

	help := bold + "NAME\n\t" + cmd + clear + " \u2014 compile Colfer schemas\n\n"
	help += bold + "SYNOPSIS\n\t" + cmd + clear
	help += " [ " + underline + "options" + clear + " ] " + underline + "language" + clear
	help += " [ " + underline + "file" + clear + " " + underline + "..." + clear + " ]\n\n"
	help += bold + "DESCRIPTION\n\t" + clear
	help += "Generates source code for a " + underline + "language" + clear + ". The options are: "
	help += bold + "C" + clear + ", " + bold + "Go" + clear + ",\n"
	help += "\t" + bold + "Java" + clear + " and " + bold + "JavaScript" + clear + ".\n"
	help += "\tThe " + underline + "file" + clear + " operands specify schema input. Directories are scanned\n"
	help += "\tfor files with the colf extension. When no files are given, then\n"
	help += "\tthe current " + italic + "working directory" + clear + " is used.\n"
	help += "\tA package definition may be spread over several schema files.\n"
	help += "\tThe directory hierarchy of the input is not relevant for the\n"
	help += "\tgenerated code.\n\n"
	help += bold + "OPTIONS\n" + clear

	tail := "\n" + bold + "EXIT STATUS" + clear + "\n"
	tail += "\tThe command exits 0 on succes, 1 on compilation failure and 2\n"
	tail += "\twhen invoked without arguments.\n"
	tail += "\n" + bold + "EXAMPLES" + clear + "\n"
	tail += "\tCompile ./io.colf with compact limits as C:\n\n"
	tail += "\t\t" + cmd + " -b src -s 2048 -l 96 C io.colf\n\n"
	tail += "\tCompile ./api/*.colf in package com.example as Java:\n\n"
	tail += "\t\t" + cmd + " -p com/example -x com/example/Parent Java api\n"
	tail += "\n" + bold + "BUGS" + clear + "\n"
	tail += "\tReport bugs at <https://github.com/pascaldekloe/colfer/issues>.\n\n"
	tail += "\tText validation is not part of the marshalling and unmarshalling\n"
	tail += "\tprocess. C and Go just pass any malformed UTF-8 characters. Java\n"
	tail += "\tand JavaScript replace unmappable content with the '?' character\n"
	tail += "\t(ASCII 63).\n\n"
	tail += bold + "SEE ALSO\n\t" + clear + "protoc(1), flatc(1)\n"

	flag.Usage = func() {
		os.Stderr.WriteString(help)
		flag.PrintDefaults()
		os.Stderr.WriteString(tail)
	}
}
