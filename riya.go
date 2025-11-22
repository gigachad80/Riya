package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

type Category struct {
	Description string   `yaml:"description"`
	Patterns    []string `yaml:"patterns"`
	Examples    []string `yaml:"examples"`
}

type PatternFile struct {
	Categories map[string]Category `yaml:"categories"`
}

var (
	flagHelp       = flag.Bool("h", false, "Show help/usage")
	flagList       = flag.Bool("list", false, "List all patterns for selected categories")
	flagNoColor    = flag.Bool("no-color", false, "Disable colored output")
	flagAll        = flag.Bool("a", false, "Match all sensitive files/paths (default if no flags given)")
	flagSQL        = flag.Bool("s", false, "Match SQL/database files")
	flagGraphQL    = flag.Bool("g", false, "Match GraphQL endpoints")
	flagPHP        = flag.Bool("p", false, "Match PHP source/config files")
	flagBackup     = flag.Bool("b", false, "Match backup and temporary files")
	flagConfig     = flag.Bool("c", false, "Match config and environment files")
	flagJSJSON     = flag.Bool("j", false, "Match JavaScript and JSON files")
	flagLogs       = flag.Bool("l", false, "Match log and text files")
	flagCerts      = flag.Bool("k", false, "Match certificate and key files")
	flagFramework  = flag.Bool("f", false, "Match framework-specific files")
	flagVCS        = flag.Bool("v", false, "Match version control files")
	flagArchive    = flag.Bool("x", false, "Match archive/compressed files")
	flagCloudInfra = flag.Bool("d", false, "Match cloud and infrastructure config")
	flagAuthAdmin  = flag.Bool("m", false, "Match admin and authentication paths")
	flagDirs       = flag.Bool("r", false, "Match sensitive directories")
	flagMisc       = flag.Bool("misc", false, "Match miscellaneous sensitive files")

	// New feature flags
	flagExclude    = flag.String("exc", "", "Exclude patterns (comma-separated extensions/patterns)")
	flagInclude    = flag.String("inc", "", "Include only these patterns (comma-separated)")
	flagOutputFile = flag.String("o", "", "Output results to file (e.g., results.txt)")
	flagUnique     = flag.Bool("u", false, "Remove duplicate URLs (show each URL only once)")
	flagStats      = flag.Bool("stats", false, "Show statistics summary instead of URLs")
)

var categoryColors = map[string]string{
	"s":    "\033[31m",
	"g":    "\033[35m",
	"p":    "\033[34m",
	"b":    "\033[33m",
	"c":    "\033[36m",
	"j":    "\033[38;5;220m",
	"l":    "\033[38;5;130m",
	"k":    "\033[38;5;205m",
	"f":    "\033[34m",
	"v":    "\033[38;5;208m",
	"x":    "\033[32m",
	"d":    "\033[38;5;51m",
	"m":    "\033[38;5;118m",
	"r":    "\033[38;5;183m",
	"misc": "\033[37m",
}

const colorReset = "\033[0m"

var categoryPriority = map[string]int{
	"x": 1, "v": 2, "f": 3, "j": 4, "l": 5, "b": 6, "p": 7, "s": 8,
	"c": 9, "g": 10, "d": 11, "m": 12, "r": 13, "misc": 14, "k": 15,
}

func printHelp() {
	flags := []struct {
		flag     string
		desc     string
		colorKey string
	}{
		{"-h", "Show this help message", ""},
		{"-list", "List all known patterns for selected categories", ""},
		{"-no-color", "Disable colored output", ""},
		{"-a", "Match all sensitive files and paths (default if no flags given)", ""},
		{"-s", "Match SQL/database backup files", "s"},
		{"-g", "Match GraphQL endpoints", "g"},
		{"-p", "Match PHP source and config files", "p"},
		{"-b", "Match backup and temporary files", "b"},
		{"-c", "Match config and environment files", "c"},
		{"-j", "Match JavaScript and JSON files", "j"},
		{"-l", "Match log and text files", "l"},
		{"-k", "Match certificate and key files", "k"},
		{"-f", "Match framework-specific files", "f"},
		{"-v", "Match version control files", "v"},
		{"-x", "Match archive/compressed files", "x"},
		{"-d", "Match cloud and infrastructure config files", "d"},
		{"-m", "Match admin and authentication paths", "m"},
		{"-r", "Match sensitive directories", "r"},
		{"-misc", "Match miscellaneous sensitive files", "misc"},
	}

	fmt.Println("Usage: riya [flags]")
	fmt.Println()
	fmt.Println("Category Flags:")

	for _, f := range flags {
		colorCode := ""
		if c, ok := categoryColors[f.colorKey]; ok && !*flagNoColor {
			colorCode = c
		}
		flagText := f.flag
		if colorCode != "" {
			flagText = colorCode + flagText + colorReset
		}
		fmt.Printf("  %-15s %s\n", flagText, f.desc)
	}

	fmt.Println()
	fmt.Println("Filter & Output Flags:")
	fmt.Println("  -exc <list>     Exclude patterns (comma-separated: js,json,png)")
	fmt.Println("  -inc <list>     Include only these patterns (comma-separated)")
	fmt.Println("  -u              Remove duplicate URLs")
	fmt.Println("  -stats          Show match statistics instead of URLs")
	fmt.Println("  -o <file>       Save output to file (e.g., results.txt)")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  cat urls.txt | riya -s -p")
	fmt.Println("  waybackurls target.com | riya -a -exc js,json,css -o results.txt")
	fmt.Println("  cat urls.txt | riya -s -u -o sensitive.txt")
	fmt.Println("  waybackurls target.com | riya -inc sql,env,config -stats")
	fmt.Println("  cat urls.txt | riya -a -o links.txt")
	fmt.Println("  riya -g -list  # lists all GraphQL patterns")
}

func loadPatterns(filename string) (PatternFile, error) {
	var pf PatternFile
	data, err := os.ReadFile(filename)
	if err != nil {
		return pf, err
	}
	err = yaml.Unmarshal(data, &pf)
	return pf, err
}

func selectedCategories() []string {
	var cats []string
	if *flagAll || (!*flagSQL && !*flagGraphQL && !*flagPHP && !*flagBackup &&
		!*flagConfig && !*flagJSJSON && !*flagLogs && !*flagCerts &&
		!*flagFramework && !*flagVCS && !*flagArchive && !*flagCloudInfra &&
		!*flagAuthAdmin && !*flagDirs && !*flagMisc) {
		return []string{"all"}
	}

	flags := map[string]*bool{
		"s": flagSQL, "g": flagGraphQL, "p": flagPHP, "b": flagBackup,
		"c": flagConfig, "j": flagJSJSON, "l": flagLogs, "k": flagCerts,
		"f": flagFramework, "v": flagVCS, "x": flagArchive, "d": flagCloudInfra,
		"m": flagAuthAdmin, "r": flagDirs, "misc": flagMisc,
	}

	for cat, f := range flags {
		if *f {
			cats = append(cats, cat)
		}
	}
	return cats
}

func buildCategoryRegexMap(pf PatternFile, cats []string) (map[*regexp.Regexp]string, error) {
	result := make(map[*regexp.Regexp]string)
	for _, cat := range cats {
		if cat == "all" {
			for k, catData := range pf.Categories {
				for _, pattern := range catData.Patterns {
					re, err := regexp.Compile("(?i)" + pattern)
					if err != nil {
						return nil, fmt.Errorf("failed to compile pattern %s for category %s: %v", pattern, k, err)
					}
					result[re] = k
				}
			}
		} else {
			catData, ok := pf.Categories[cat]
			if !ok {
				continue
			}
			for _, pattern := range catData.Patterns {
				re, err := regexp.Compile("(?i)" + pattern)
				if err != nil {
					return nil, fmt.Errorf("failed to compile pattern %s for category %s: %v", pattern, cat, err)
				}
				result[re] = cat
			}
		}
	}
	return result, nil
}

func parseFilterPatterns(patternStr string) ([]*regexp.Regexp, error) {
	if patternStr == "" {
		return nil, nil
	}

	patterns := strings.Split(patternStr, ",")
	var regexes []*regexp.Regexp

	for _, p := range patterns {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}

		if !strings.Contains(p, "\\") && !strings.Contains(p, "^") && !strings.Contains(p, "$") {
			p = "\\." + strings.TrimPrefix(p, ".") + "(\\?|#|$)"
		}

		re, err := regexp.Compile("(?i)" + p)
		if err != nil {
			return nil, fmt.Errorf("invalid pattern %s: %v", p, err)
		}
		regexes = append(regexes, re)
	}
	return regexes, nil
}

func matchesAnyPattern(line string, patterns []*regexp.Regexp) bool {
	for _, re := range patterns {
		if re.MatchString(line) {
			return true
		}
	}
	return false
}

func colorize(text, color string) string {
	if *flagNoColor || *flagOutputFile != "" {
		return text
	}
	return color + text + colorReset
}

func printLists(pf PatternFile, cats []string) {
	for _, cat := range cats {
		catData, ok := pf.Categories[cat]
		if !ok {
			fmt.Printf("Unknown category key: %s\n", cat)
			continue
		}
		color := categoryColors[cat]
		fmt.Printf("%s\n", colorize(fmt.Sprintf("Category [%s] - %s:", cat, catData.Description), color))
		fmt.Println("Patterns:")
		for _, p := range catData.Patterns {
			fmt.Printf("  • %s\n", p)
		}
		fmt.Println("Examples:")
		for _, e := range catData.Examples {
			fmt.Printf("  • %s\n", colorize(e, color))
		}
		fmt.Println()
	}
}

func main() {
	flag.Parse()

	if *flagHelp {
		printHelp()
		return
	}

	pf, err := loadPatterns("patterns.yml")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load pattern file: %v\n", err)
		os.Exit(1)
	}

	cats := selectedCategories()

	if *flagList {
		printLists(pf, cats)
		return
	}

	regexMap, err := buildCategoryRegexMap(pf, cats)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Regex compilation error: %v\n", err)
		os.Exit(1)
	}

	if len(regexMap) == 0 {
		fmt.Fprintln(os.Stderr, "No valid patterns found for the selected categories")
		os.Exit(1)
	}

	// Parse exclude and include patterns
	excludePatterns, err := parseFilterPatterns(*flagExclude)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing exclude patterns: %v\n", err)
		os.Exit(1)
	}

	includePatterns, err := parseFilterPatterns(*flagInclude)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing include patterns: %v\n", err)
		os.Exit(1)
	}

	// Setup output destination
	var output *os.File
	var writer *bufio.Writer

	if *flagOutputFile != "" {
		output, err = os.Create(*flagOutputFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating output file: %v\n", err)
			os.Exit(1)
		}
		defer output.Close()
		writer = bufio.NewWriter(output)
		defer writer.Flush()
	} else {
		output = os.Stdout
		writer = bufio.NewWriter(output)
		defer writer.Flush()
	}

	// Tracking for unique and stats modes
	seenURLs := make(map[string]bool)
	categoryStats := make(map[string]int)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()

		// Skip if already seen (when -u flag is set)
		if *flagUnique {
			if seenURLs[line] {
				continue
			}
			seenURLs[line] = true
		}

		// Check exclusion patterns
		if len(excludePatterns) > 0 && matchesAnyPattern(line, excludePatterns) {
			continue
		}

		// Check inclusion patterns (if specified)
		if len(includePatterns) > 0 && !matchesAnyPattern(line, includePatterns) {
			continue
		}

		// Match against category patterns
		bestCat := ""
		bestPrio := 9999
		for re, cat := range regexMap {
			if re.MatchString(line) {
				prio := categoryPriority[cat]
				if prio < bestPrio {
					bestPrio = prio
					bestCat = cat
				}
			}
		}

		if bestCat != "" {
			if *flagStats {
				categoryStats[bestCat]++
			} else {
				color := categoryColors[bestCat]
				fmt.Fprintln(writer, colorize(line, color))
			}
		}
	}

	// Print statistics if -stats flag is set
	if *flagStats {
		fmt.Fprintln(writer)
		fmt.Fprintln(writer, "=== Match Statistics ===")
		totalMatches := 0
		for cat, count := range categoryStats {
			catData := pf.Categories[cat]
			color := categoryColors[cat]
			fmt.Fprintf(writer, "%s [%s] %s: %d matches\n",
				colorize("●", color),
				cat,
				catData.Description,
				count)
			totalMatches += count
		}
		fmt.Fprintf(writer, "\nTotal matches: %d\n", totalMatches)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
	}
}

