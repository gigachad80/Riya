🚀 Project Name : Riya
===============

### Riya : A grep like tool for filtering sensitive files, paths, and endpoints from URL lists with pattern-based categorization   

 
![Maintenance](https://img.shields.io/badge/Maintained%3F-yes-purple.svg)
<a href="https://github.com/gigachad80/Riya/issues"><img src="https://img.shields.io/badge/contributions-welcome-brightgreen.svg?style=flat"></a>

## Table of Contents

* [📌 Overview](#-overview)
* [✨ Features](#-features)
* [🎯 Before & After Riya](#-before--after-riya)
* [📚 Requirements & Dependencies](#-requirements--dependencies)
* [📥 Installation Guide](#-installation-guide)
* [🚀 Usage](#-usage)
  - [Basic Usage](#basic-usage)
  - [Advanced Filtering](#advanced-filtering)
* [🔧 Technical Details](#-technical-details)
* [🤔 Why This Name?](#-why-this-name)
* [⌚ Development Time](#-development-time)
* [🙃 Why I Created This](#-why-i-created-this)
* [🙏 Credits & Inspiration](#-credits--inspiration)
* [📞 Contact](#-contact)
* [📄 License](#-license)

### 📌 Overview

**Riya** is a blazing-fast CLI tool designed for bug bounty hunters, security researchers, and penetration testers to efficiently filter and categorize sensitive files, paths, and endpoints from large URL lists. Instead of running dozens of grep commands with complex regex patterns, Riya provides a clean, intuitive interface with pre-configured pattern matching for common security-sensitive files.

**Key Capabilities:**
* 15+ pre-configured security-focused categories (SQL dumps, config files, API endpoints, etc.)
* Color-coded output for instant visual categorization
* Advanced filtering with include/exclude patterns
* Duplicate URL removal and statistics mode
* File output support for pipeline integration
* Pattern listing for transparency and customization

### ✨ Features

### 🎯 Smart Pattern Matching
- **15+ Categories** - SQL, GraphQL, PHP, backups, configs, logs, certificates, and more
- **Priority-Based Matching** - Automatically categorizes URLs by importance
- **Case-Insensitive** - Works with any URL format
- **YAML-Powered** - Easy-to-edit pattern configuration file
- **Regex Engine** - Robust pattern matching with Go's regexp library
- **Progress Friendly** - Works seamlessly with pipelines and redirects

### ⚡ Performance Features
- **Stream Processing** - Memory-efficient handling of large URL lists
- **Duplicate Removal** (`-u`) - Show each URL only once
- **Statistics Mode** (`-stats`) - Get match counts per category
- **File Output** (`-o`) - Save results to file for later analysis
- **Fast Filtering** - Process thousands of URLs in seconds

### 🎛️ Advanced Filtering
- **Exclude Patterns** (`-exc`) - Remove unwanted file types (js, css, png, etc.)
- **Include Patterns** (`-inc`) - Focus on specific patterns only
- **URL Query Support** - Properly handles URLs with parameters and fragments
- **Flexible Syntax** - Simple extensions or complex regex patterns

### 🎯 Before & After Riya

| Aspect | 😫 Before Riya (The Grep Nightmare) | ✨ After Riya (Clean & Simple) |
|--------|-------------------------------------|--------------------------------|
| **Commands** | 😤 Running 10+ separate grep patterns for different file types | 🚀 One simple command: `riya -s -c -k` |
| **Syntax** | 🤯 Memorizing cryptic patterns like `grep -E '\.(sql\|db\|sqlite)(\?\|#\|$)'` | 😎 Human-readable flags: `-s` for SQL, `-c` for configs |
| **Workflow** | ⏰ Manually filtering, deduplicating, and organizing results | ⚡ Instant color-coded, organized, deduplicated results |
| **Reliability** | ❌ Error-prone - easy to miss patterns or make regex mistakes | ✅ Pre-tested patterns covering edge cases and URL variations |


### 💡 Why This Matters

When analyzing thousands of URLs from tools like `waybackurls` or `gau`, manually grepping for sensitive files becomes tedious and error-prone. Riya transforms this workflow from:

```bash
# The old way 😓
cat urls.txt | grep -iE '\.sql$' > sql.txt
cat urls.txt | grep -iE '\.env$' > env.txt
cat urls.txt | grep -iE 'config\.(php|js|json)' > configs.txt
cat urls.txt | grep -iE '\.(key|pem|crt)$' > certs.txt
# ... repeat 15+ times ... 😵
```

To this:

```bash
# The Riya way 🎯
cat urls.txt | riya -a -o results.txt
```

### 📚 Requirements & Dependencies

* **Go 1.19+** - For building from source
* **patterns.yml** - Pattern configuration file (included in repository)

### 📥 Installation Guide

### ⚡ Quick Install

**Method 1: Build from Source**
```bash
git clone https://github.com/gigachad80/riya
cd riya
go build -o riya riya.go
```

OR 

**Method 2: Download Binary**
Download the latest binary from the [releases page](https://github.com/gigachad80/riya/releases) and add it to your PATH.

### 📂 Required Files
Make sure `patterns.yml` is in the same directory as the `riya` binary, or in your current working directory.

### 🚀 Usage

### Basic Usage

```bash
# Show all sensitive files (default behavior)
cat urls.txt | riya
cat urls.txt | riya -a

# Filter specific categories
cat urls.txt | riya -s -c -k  # SQL + configs + certificates
waybackurls target.com | riya -p -b  # PHP + backups

# View available patterns
riya -g -list  # Show all GraphQL patterns
riya -s -list  # Show all SQL patterns
```

### Advanced Filtering

```bash
# Exclude common noise (JS, CSS, images)
waybackurls target.com | riya -a -exc js,css,png,jpg,gif

# Include only specific patterns
cat urls.txt | riya -inc sql,env,config

# Remove duplicates and save to file
cat urls.txt | riya -s -p -u -o sensitive.txt

# Get statistics instead of URLs
waybackurls target.com | riya -a -stats

# Combine multiple filters
waybackurls target.com | riya -s -c -k -exc js,json -u -o critical.txt
```

### Real-World Examples

```bash
# Bug bounty recon pipeline
subfinder -d target.com | waybackurls | riya -a -exc js,css,woff -o findings.txt

# Focus on high-value targets
cat wayback.txt | riya -s -c -k -u > high_priority.txt

# Quick overview of what's exposed
echo "https://target.com" | waybackurls | riya -a -stats
```


### 🔧 Technical Details

### Architecture
- **Stream-Based Processing** - Memory-efficient reading from stdin
- **Regex Compilation** - Pre-compiled patterns for optimal performance
- **Priority System** - Categorizes URLs by security importance
- **YAML Configuration** - Easy-to-modify pattern definitions
- **Buffer Management** - Efficient I/O with buffered writers

### Pattern Matching
- **Case-Insensitive** - All patterns use `(?i)` flag
- **URL-Aware** - Handles query parameters (`?`) and fragments (`#`)
- **Flexible Syntax** - Supports simple extensions and complex regex
- **Deduplication** - Optional unique URL filtering with `-u` flag


### 🤔 Why This Name?

Okay so like... 😅 this name randomly popped into my head and the name reminded me of this friend from my teenage days (like 7-9 years back 💭 ). Her name was Riya and we used to play together back then 🎮✨.
The games we played? Yeah... can't exactly disclose that in an open source README 💀😂 (let's just say it's better left to imagination 🙈). But yeah, good times! So I just named this tool after her as a little throwback to those days 🕰️💫.
Fun fact: She probably has no idea a security tool is named after her now lmaooo 😭🤣
 ╰(*°▽°*)╯

### ⌚ Development Time

From initial concept to feature-complete implementation, including pattern research, testing, and documentation, the development took approximately **1 hr 31 min 58 sec** across multiple sessions.

### 🙃 Why I Created This

I was tired of writing complex grep patterns every single time I analyzed URLs. Instead of memorizing regex syntax like `grep -iE '\.(sql|db|sqlite3?)(\?|#|$)'`, I wanted something with **clean, simple flags** like `-s` for SQL files.

The goal was simple: replace the mess of multiple grep commands with one intuitive tool that just works. No more regex headaches, no more forgotten patterns, just straightforward filtering with human-readable options.

### 🙏 Credits & Inspiration


This tool was inspired by a **security researcher on Linkedin** (whose handle I unfortunately can't remember - will update this once I find them! 🔍). 

I kept seeing them tweet screenshots of their bug bounty workflow where they were running these insanely complex grep commands like:

```bash
grep -iE '\.(sql|db|sqlite3?)(\?|#|$)' urls.txt
grep -iE '\.(env|config|settings)\.(js|json|php|yml)' urls.txt  
grep -iE '(backup|old|temp).*\.(zip|tar|gz|sql)' urls.txt
grep -iE '\.(key|pem|crt|p12|pfx)(\?|#|$)' urls.txt
# ... and like 10+ more patterns
```

Watching them juggle all these complex regex patterns in their threads made me think: **"Bruh, why not just make ONE tool with simple flags instead of this grep hell?"** 

Instead of memorizing and typing out `grep -iE '\.(sql|db|sqlite3?)(\?|#|$)'` every single time, just do `riya -s`. That's it. That's the whole vibe.

**Huge thanks to that researcher!** Your daily struggle with grep inspired this tool 🙏. If anyone recognizes this workflow or knows who I'm talking about, please hit me up so I can give proper credit! 

I noticed they were constantly running multiple grep commands in their bug bounty workflow:
- `grep` for SQL files
- Another `grep` for API keys
- Yet another for config files
- Separate greps for backups, logs, certificates...

Watching them juggle 10+ different grep patterns made me think: **"Why not create a single tool that does all of this with simple flags?"**

Instead of:
```bash
grep -iE '\.sql'
grep -iE '\.env'
grep -iE '\.(key|pem)'
# ... and so on
```

Just do:
```bash
riya -s -c -k
```

**Thank you, mystery researcher!** Your workflow chaos inspired something (hopefully) useful. If anyone knows who this might be, please let me know so I can give proper credit! 🙏

### 📞 Contact

📧 Email: **pookielinuxuser@tutamail.com**

### 📄 License

**MIT License** 

First Published: November 22, 2025
Last Updated : Nov 22nd , 2025

**Made with ❤️ in Go** - Because security researchers deserve better than grep hell.

### 🌟 Star History

If you find this tool useful, please consider giving it a star! It helps others discover the project.


