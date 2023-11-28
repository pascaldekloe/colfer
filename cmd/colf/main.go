package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/pascaldekloe/colfer"
)

var (
	basedir = flag.String("b", ".", "Use a base `directory` for the generated code.")
	prefix  = flag.String("p", "", "Compile to a `package` prefix.")
	format  = flag.Bool("f", false, "Normalize the format of all schema input on the fly.")
	verbose = flag.Bool("v", false, "Enable verbose reporting to "+italic+"standard error"+clear+".")

	superClass  = flag.String("x", "", "Make all generated classes extend a super `class`.")
	interfaces  = flag.String("i", "", "Make all generated classes implement one or more `interfaces`.\nUse commas as a list separator.")
	tagFiles    = flag.String("t", "", "Supply custom tags with one or more `files`. Use commas as a list\nseparator. See the TAGS section for details.")
	snippetFile = flag.String("c", "", "Insert a code snippet from a `file`.")
)

func init() {
	flag.Bool("h", false, "Prints the manual to standard error.")
	flag.Usage = printManual
}

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
	var tagOptions colfer.TagOptions
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
		if *tagFiles != "" {
			log.Fatalf("%s: tags not supported with C", name)
		}
		if *snippetFile != "" {
			log.Fatalf("%s: snippet not supported with C", name)
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
		tagOptions.FieldAllow = colfer.TagSingle

	case "java":
		report.Print("set-up for Java")
		gen = colfer.GenerateJava
		tagOptions.StructAllow = colfer.TagMulti
		tagOptions.FieldAllow = colfer.TagMulti

	case "javascript", "js", "ecmascript":
		report.Print("set-up for ECMAScript")
		gen = colfer.GenerateECMA
		if *superClass != "" {
			log.Fatalf("%s: super class not supported with ECMAScript", name)
		}
		if *interfaces != "" {
			log.Fatalf("%s: interfaces not supported with ECMAScript", name)
		}
		if *tagFiles != "" {
			log.Fatalf("%s: tags not supported with ECMAScript", name)
		}
		if *snippetFile != "" {
			log.Fatalf("%s: snippet not supported with ECMAScript", name)
		}

	default:
		log.Fatalf("%s: unsupported language %q", name, lang)
	}

	if flag.NArg() > 1 {
		mustResolveSchemaFiles(flag.Args()[1:]...)
	} else {
		mustResolveSchemaFiles(".")
	}
	packages, err := colfer.ParseFiles(schemaPaths...)
	if err != nil {
		log.Fatal(err)
	}

	if *tagFiles != "" {
		for _, path := range strings.Split(*tagFiles, ",") {
			report.Print("using tag file: ", path)
			if err = packages.ApplyTagFile(path, tagOptions); err != nil {
				log.Fatal(err)
			}
		}
	}

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

	report.Print("using schema file: ", path)
	schemaPaths = append(schemaPaths, path)
	schemaInfos = append(schemaInfos, info)
}

// ANSI escape codes for markup
const (
	bold   = "\x1b[1m"
	italic = "\x1b[3m"
	clear  = "\x1b[0m"
)

func printManual() {
	nameSection := bold + "NAME\n\t" + name + clear + " \u2014 compile Colfer schemas\n"

	synopsisSection := bold + "SYNOPSIS\n\t" + name + clear + " [" + bold + "-h" + clear + "]\n\t" +
		bold + name + clear + " [" + bold + "-vf" + clear + "] [" +
		bold + "-b" + clear + " directory] [" +
		bold + "-p" + clear + " package] " + bold + "C" + clear +
		" [file ...]\n\t" +
		bold + name + clear + " [" + bold + "-vf" + clear + "] [" +
		bold + "-b" + clear + " directory] [" +
		bold + "-p" + clear + " package] [" +
		bold + "-t" + clear + " files] " + bold + "Go" + clear +
		" [file ...]\n\t" +
		bold + name + clear + " [" + bold + "-vf" + clear + "] [" +
		bold + "-b" + clear + " directory] [" +
		bold + "-p" + clear + " package] [" +
		bold + "-t" + clear + " files] \\\n\t\t[" +
		bold + "-x" + clear + " class] [" +
		bold + "-i" + clear + " interfaces] [" +
		bold + "-c" + clear + " file] " + bold + "Java" + clear +
		" [file ...]\n\t" +
		bold + name + clear + " [" + bold + "-vf" + clear + "] [" +
		bold + "-b" + clear + " directory] [" +
		bold + "-p" + clear + " package] " + bold + "JavaScript" + clear +
		" [file ...]\n"

	descriptionSection := bold + "DESCRIPTION" + clear + "\n" +
		"\tThe output is source code for either C, Go, Java or JavaScript.\n\n" +
		"\tFor each operand that names a file of a type other than\n" +
		"\tdirectory, " + bold + "colf" + clear + " reads the content as schema input. For each\n" +
		"\tnamed directory, " + bold + "colf" + clear + " reads all files with a .colf extension\n" +
		"\twithin that directory. If no operands are given, the contents of\n" +
		"\tthe current directory are used.\n\n" +
		"\tA package definition may be spread over several schema files.\n" +
		"\tThe directory hierarchy of the input is not relevant to the\n" +
		"\tgenerated code.\n"

	tagsSection := bold + "TAGS" + clear + "\n" +
		"\tTags, a.k.a. annotations, are source code additions for structs\n" +
		"\tand/or fields. Input for the compiler can be specified with the\n" +
		bold + "\t-t" + clear + " option. The data format is " + italic +
		"line-oriented" + clear + ".\n\n" +
		"\t\t<line> :≡ <qual> <space> <code> ;\n" +
		"\t\t<qual> :≡ <package> '.' <dest> ;\n" +
		"\t\t<dest> :≡ <struct> | <struct> '.' <field> ;\n\n" +
		"\tLines starting with a '#' are ignored (as comments). Java output\n" +
		"\tcan take multiple tag lines for the same struct or field. Each\n" +
		"\tcode line is applied in order of appearance.\n"

	exitStatusSection := bold + "EXIT STATUS" + clear + "\n" +
		"\tThe command exits 0 on success, 1 on error and 2 when invoked\n" +
		"\twithout arguments.\n"

	examplesSection := bold + "EXAMPLES" + clear + "\n" +
		"\tCompile ./io.colf into the src directory as C:\n\n" +
		"\t\t" + name + " -b src C io.colf\n\n" +
		"\tCompile ./*.colf with a common parent as Java:\n\n" +
		"\t\t" + name + " -p com.example.model -x com.example.io.IOBean Java\n"

	bugsSection := bold + "BUGS" + clear + "\n" +
		"\tReport bugs at <https://github.com/pascaldekloe/colfer/issues>.\n\n" +
		"\tText validation is not part of the marshalling and unmarshalling\n" +
		"\tprocess. C and Go just pass any malformed UTF-8 characters. Java\n" +
		"\tand JavaScript replace unmappable content with the '?' character\n" +
		"\t(ASCII 63).\n"

	seeAlsoSection := bold + "SEE ALSO" + clear + "\n\tprotoc(1), flatc(1)\n"

	w := flag.CommandLine.Output()
	fmt.Fprintln(w, nameSection)
	fmt.Fprintln(w, synopsisSection)
	fmt.Fprintln(w, descriptionSection)
	fmt.Fprintln(w, bold+"OPTIONS"+clear)
	flag.PrintDefaults()
	fmt.Fprintln(w)
	fmt.Fprintln(w, tagsSection)
	fmt.Fprintln(w, exitStatusSection)
	fmt.Fprintln(w, examplesSection)
	fmt.Fprintln(w, bugsSection)
	fmt.Fprint(w, seeAlsoSection)
}
