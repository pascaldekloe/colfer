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
	bold   = "\x1b[1m"
	italic = "\x1b[3m"
	clear  = "\x1b[0m"
)

var (
	basedir = flag.String("b", ".", "Use a base `directory` for the generated code.")
	prefix  = flag.String("p", "", "Compile to a `package` prefix.")
	format  = flag.Bool("f", false, "Normalize the format of all input schemas on the fly.")
	verbose = flag.Bool("v", false, "Enable verbose reporting to "+italic+"standard error"+clear+".")

	sizeMax = flag.String("s", "16 * 1024 * 1024", "Set the default upper limit for serial byte sizes. The\n`expression` is applied to the target language under the name\nColferSizeMax.")
	listMax = flag.String("l", "64 * 1024", "Set the default upper limit for the number of elements in a\nlist. The `expression` is applied to the target language under\nthe name ColferListMax.")

	superClass  = flag.String("x", "", "Make all generated classes extend a super `class`.")
	interfaces  = flag.String("i", "", "Make all generated classes implement one or more `interfaces`.\nUse commas as a list separator.")
	snippetFile = flag.String("c", "", "Insert a code snippet from a `file`.")

	tagFiles = flag.String("t", "", "Supply custom tags with one or more `files`. Use commas as a list\nseparator. See the TAGS section for details.")
)

var name = os.Args[0]
var report = log.New(ioutil.Discard, os.Args[0]+": ", 0)

var (
	schemaPaths []string      // source files in use
	schemaInfos []os.FileInfo // corresponding descriptors
)

func main() {
	flag.Parse()

	log.SetFlags(0)
	if *verbose {
		report.SetOutput(os.Stderr)
	}

	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(2)
	}

	// select language
	var gen func(string, colfer.Packages) error
	switch lang := flag.Arg(0); strings.ToLower(lang) {
	case "c":
		report.Print("set-up for C")
		gen = colfer.GenerateC
		if *superClass != "" {
			log.Fatalf("%s: super class not supported with C", name)
		}
		if *interfaces != "" {
			log.Fatalf("%s: interfaces not supported with C", name)
		}
		if *snippetFile != "" {
			log.Fatalf("%s: snippet not supported with C", name)
		}
		if *tagFiles != "" {
			log.Fatalf("%s: tags not supported with C", name)
		}

	case "go":
		report.Print("set-up for Go")
		gen = colfer.GenerateGo
		if *superClass != "" {
			log.Fatalf("%s: super class not supported with Go", name)
		}
		if *interfaces != "" {
			log.Fatalf("%s: interfaces not supported with Go", name)
		}
		if *snippetFile != "" {
			log.Fatalf("%s: snippet not supported with Go", name)
		}

	case "java":
		report.Print("set-up for Java")
		gen = colfer.GenerateJava

	case "javascript", "js", "ecmascript":
		report.Print("set-up for ECMAScript")
		gen = colfer.GenerateECMA
		if *superClass != "" {
			log.Fatalf("%s: super class not supported with ECMAScript", name)
		}
		if *interfaces != "" {
			log.Fatalf("%s: interfaces not supported with ECMAScript", name)
		}
		if *snippetFile != "" {
			log.Fatalf("%s: snippet not supported with ECMAScript", name)
		}
		if *tagFiles != "" {
			log.Fatalf("%s: tags not supported with ECMAScript", name)
		}

	default:
		log.Fatalf("%s: unsupported language %q", name, lang)
	}

	if flag.NArg() > 1 {
		mustResolveSchemaFiles(flag.Args()[1:]...)
	} else {
		mustResolveSchemaFiles(".")
	}
	report.Print("using schema files: ", strings.Join(schemaPaths, ", "))

	packages, err := colfer.ParseFiles(schemaPaths...)
	if err != nil {
		log.Fatal(err)
	}

	if *tagFiles != "" {
		for _, path := range strings.Split(*tagFiles, ",") {
			if err = packages.ApplyTagFile(path); err != nil {
				log.Fatal(err)
			}
		}
	}

	report.Print("using schema files: ", strings.Join(schemaPaths, ", "))

	if *format {
		for _, path := range schemaPaths {
			changed, err := colfer.FormatFile(path)
			if err != nil {
				log.Fatal(err)
			}
			if changed {
				log.Printf("%s: formatted %s", name, path)
			}
		}
	}

	if len(packages) == 0 {
		log.Fatalf("%s: no struct definitons found", name)
	}

	for _, p := range packages {
		p.Name = path.Join(*prefix, p.Name)
		p.SizeMax = *sizeMax
		p.ListMax = *listMax
		p.SuperClass = *superClass
		if *interfaces != "" {
			p.Interfaces = strings.Split(*interfaces, ",")
		}
		if len(*snippetFile) > 0 {
			snippet, err := ioutil.ReadFile(*snippetFile)
			if err != nil {
				log.Fatal(err)
			}
			p.CodeSnippet = string(snippet)
		}
	}

	if err := gen(*basedir, packages); err != nil {
		log.Fatal(err)
	}
}

func mustResolveSchemaFiles(paths ...string) {
	for _, path := range paths {
		info, err := os.Stat(path)
		if err != nil {
			log.Fatal(err)
		}

		if !info.IsDir() {
			addSchemaFile(path, info)
			continue
		}

		children, err := filepath.Glob(filepath.Join(path, "*.colf"))
		if err != nil {
			log.Fatal(err)
		}
		for _, path = range children {
			info, err = os.Stat(path)
			if err != nil {
				log.Fatal(err)
			}
			if !info.IsDir() {
				addSchemaFile(path, info)
			}
		}
	}
}

func addSchemaFile(path string, info os.FileInfo) {
	for _, previous := range schemaInfos {
		if os.SameFile(info, previous) {
			report.Printf("duplicate inclusion of %q ignored", path)
			return
		}
	}

	schemaPaths = append(schemaPaths, path)
	schemaInfos = append(schemaInfos, info)
}

func init() {
	help := bold + "NAME\n\t" + name + clear + " \u2014 compile Colfer schemas\n\n"
	help += bold + "SYNOPSIS\n\t" + name + clear + " [" + bold + "-h" + clear + "]\n\t"

	help += bold + name + clear + " [" + bold + "-vf" + clear + "] ["
	help += bold + "-b" + clear + " directory] ["
	help += bold + "-p" + clear + " package] \\\n\t\t["
	help += bold + "-s" + clear + " expression] ["
	help += bold + "-l" + clear + " expression] C"
	help += " [file ...]\n\t"

	help += bold + name + clear + " [" + bold + "-vf" + clear + "] ["
	help += bold + "-b" + clear + " directory] ["
	help += bold + "-p" + clear + " package] ["
	help += bold + "-t" + clear + " files] \\\n\t\t["
	help += bold + "-s" + clear + " expression] ["
	help += bold + "-l" + clear + " expression] Go"
	help += " [file ...]\n\t"

	help += bold + name + clear + " [" + bold + "-vf" + clear + "] ["
	help += bold + "-b" + clear + " directory] ["
	help += bold + "-p" + clear + " package] ["
	help += bold + "-t" + clear + " files] \\\n\t\t["
	help += bold + "-x" + clear + " class] ["
	help += bold + "-i" + clear + " interfaces] ["
	help += bold + "-c" + clear + " file] \\\n\t\t["
	help += bold + "-s" + clear + " expression] ["
	help += bold + "-l" + clear + " expression] Java"
	help += " [file ...]\n\t"

	help += bold + name + clear + " [" + bold + "-vf" + clear + "] ["
	help += bold + "-b" + clear + " directory] ["
	help += bold + "-p" + clear + " package] \\\n\t\t["
	help += bold + "-s" + clear + " expression] ["
	help += bold + "-l" + clear + " expression] JavaScript"
	help += " [file ...]\n\n"

	help += bold + "DESCRIPTION" + clear + "\n"
	help += "\tGenerates source code from a model definition for one language.\n"
	help += "\tThe file operands specify schema input. Directories are scanned\n"
	help += "\tfor files with the colf extension. When no files are given, then\n"
	help += "\tthe current " + italic + "working directory" + clear + " is used.\n"
	help += "\tA package definition may be spread over several schema files.\n"
	help += "\tThe directory hierarchy of the input is not relevant for the\n"
	help += "\tgenerated code.\n\n"

	help += bold + "OPTIONS\n" + clear
	// … rendered with the flag package
	tail := "\n"

	tail += bold + "TAGS" + clear + "\n"
	tail += "\tTags, a.k.a. annotations, are source code additions for structs\n"
	tail += "\tand fields. The compiler uses a line-orient data format in the\n"
	tail += "\tform of: <package> '.' <structure> ('.' <field>) <space> <code>.\n"
	tail += "\tLines starting with a '#' are ignored (as comments).\n\n"

	tail += bold + "EXIT STATUS" + clear + "\n"
	tail += "\tThe command exits 0 on succes, 1 on compilation failure and 2\n"
	tail += "\twhen invoked without arguments.\n\n"

	tail += bold + "EXAMPLES" + clear + "\n"
	tail += "\tCompile ./io.colf with compact limits as C:\n\n"
	tail += "\t\t" + name + " -b src -s 2048 -l 96 C io.colf\n\n"
	tail += "\tCompile ./*.colf with a common parent as Java:\n\n"
	tail += "\t\t" + name + " -p com.example.model -x com.example.io.IOBean Java\n\n"

	tail += bold + "BUGS" + clear + "\n"
	tail += "\tReport bugs at <https://github.com/pascaldekloe/colfer/issues>.\n\n"
	tail += "\tText validation is not part of the marshalling and unmarshalling\n"
	tail += "\tprocess. C and Go just pass any malformed UTF-8 characters. Java\n"
	tail += "\tand JavaScript replace unmappable content with the '?' character\n"
	tail += "\t(ASCII 63).\n\n"

	tail += bold + "SEE ALSO" + clear + "\n\tprotoc(1), flatc(1)\n"

	flag.Usage = func() {
		os.Stderr.WriteString(help)
		flag.PrintDefaults()
		os.Stderr.WriteString(tail)
	}
}
