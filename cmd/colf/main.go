package main

import (
	"flag"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/pascaldekloe/colfer"
)

var basedir = flag.String("b", ".", "Use a specific destination base directory.")
var prefix = flag.String("p", "", "Adds a package prefix. Use slash as a separator when nesting.")

func main() {
	log.SetFlags(0)
	flag.Parse()

	var files []string
	switch args := flag.Args(); len(args) {
	case 0:
		flag.Usage()
		os.Exit(2)
	case 1:
		colfFiles, err := filepath.Glob("*.colf")
		if err != nil {
			log.Fatal(err)
		}
		files = colfFiles
	default:
		files = args[1:]
	}

	var gen func(string, []*colfer.Package) error
	switch lang := flag.Arg(0); strings.ToLower(lang) {
	case "go":
		gen = colfer.Generate
	case "java":
		gen = colfer.GenerateJava
	case "ecmascript", "javascript", "js":
		gen = colfer.GenerateECMA
	default:
		log.Fatalf("colf: unsupported language %q", lang)
	}

	packages, err := colfer.ReadDefs(files)
	if err != nil {
		log.Fatal(err)
	}
	if len(packages) == 0 {
		log.Fatal("colfer: no struct definitons found")
	}

	for _, p := range packages {
		p.Name = path.Join(*prefix, p.Name)
	}

	if err := gen(*basedir, packages); err != nil {
		log.Fatal(err)
	}
}

// ANSI escape codes for markup
const (
	bold      = "\x1b[1m"
	underline = "\x1b[4m"
	clear     = "\x1b[0m"
)

func init() {
	cmd := os.Args[0]

	help := bold + "NAME\n\t" + cmd + clear + " \u2014 compile Colfer schemas\n\n"
	help += bold + "SYNOPSIS\n\t" + cmd + clear
	help += " [" + bold + "-b" + clear + " <" + underline + "dir" + clear + ">]"
	help += " [" + bold + "-p" + clear + " <" + underline + "path" + clear + ">]"
	help += " <" + underline + "language" + clear
	help += "> [<" + underline + "file" + clear + "> " + underline + "..." + clear + "]\n\n"
	help += bold + "DESCRIPTION\n\t" + clear
	help += "Generates source code for the given " + underline + "language" + clear
	help += ". The options are: " + bold + "Go" + clear + ",\n"
	help += "\t" + bold + "Java" + clear + " and " + bold + "ECMAScript" + clear + ".\n"
	help += "\tThe " + underline + "file" + clear + " operands are processed in command-line order. If " + underline + "file" + clear + " is\n"
	help += "\tabsent, " + cmd + " reads all \".colf\" files in the working directory.\n\n"

	tail := "\n" + bold + "BUGS" + clear
	tail += "\n\tReport bugs at https://github.com/pascaldekloe/colfer/issues\n\n"
	tail += bold + "SEE ALSO\n\t" + clear + "protoc(1)\n"

	flag.Usage = func() {
		os.Stderr.WriteString(help)
		flag.PrintDefaults()
		os.Stderr.WriteString(tail)
	}
}
