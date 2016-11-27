package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/pascaldekloe/colfer"
)

var basedir = flag.String("b", ".", "Use a specific destination base directory.")
var prefix = flag.String("p", "", "Adds a package prefix. Use slash as a separator when nesting.")
var verbose = flag.Bool("v", false, "Enables verbose reporting to the standard output.")

func main() {
	flag.Parse()

	log.SetFlags(0)

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
	var gen func(string, []*colfer.Package) error
	switch lang := flag.Arg(0); strings.ToLower(lang) {
	case "ecmascript", "javascript", "js":
		if *verbose {
			fmt.Println("Set up for ECMAScript")
		}
		gen = colfer.GenerateECMA
	case "go":
		if *verbose {
			fmt.Println("Set up for Go")
		}
		gen = colfer.GenerateGo
	case "java":
		if *verbose {
			fmt.Println("Set up for Java")
		}
		gen = colfer.GenerateJava
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
				if *verbose {
					fmt.Println("Dupe schema file", f)
				}
				break
			}
		}
	}
	files = files[:writeIndex]
	if *verbose {
		fmt.Println("Found schema files", strings.Join(files, ", "))
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
	help += " [" + bold + "-v" + clear + "]"
	help += " <" + underline + "language" + clear
	help += "> [<" + underline + "file" + clear + "> " + underline + "..." + clear + "]\n\n"
	help += bold + "DESCRIPTION\n\t" + clear
	help += "Generates source code for the given " + underline + "language" + clear
	help += ". The options are: " + bold + "Go" + clear + ",\n"
	help += "\t" + bold + "Java" + clear + " and " + bold + "ECMAScript" + clear + ".\n"
	help += "\tThe " + underline + "file" + clear + " operands specify the input. Directories are scanned for\n"
	help += "\tfiles with the colf extension. If " + underline + "file" + clear + " is absent, " + cmd + " includes\n"
	help += "\tthe working directory.\n"
	help += "\tA package can have multiple schema files.\n\n"

	tail := "\n" + bold + "EXIT STATUS" + clear + "\n"
	tail += "\tThe command exits 0 on succes, 1 on compilation failure and 2 when\n"
	tail += "\tinvoked without arguments.\n"
	tail += "\n" + bold + "EXAMPLES" + clear + "\n"
	tail += "\tCompile ./src/main/colfer/*.colf into ./target/ as Java:\n\n"
	tail += "\t\t" + cmd + " -p com/example -b target java src/main/colfer\n"
	tail += "\n" + bold + "BUGS" + clear + "\n"
	tail += "\tReport bugs at https://github.com/pascaldekloe/colfer/issues\n\n"
	tail += bold + "SEE ALSO\n\t" + clear + "protoc(1)\n"

	flag.Usage = func() {
		os.Stderr.WriteString(help)
		flag.PrintDefaults()
		os.Stderr.WriteString(tail)
	}
}
