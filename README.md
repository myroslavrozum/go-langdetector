# Go Langdetector

This project is a Go-based rewrite of an older language detector, available at [github.com/myroslavrozum/langdetector](https://github.com/myroslavrozum/langdetector) to brush-up my rusty knowledge of Golang.

It uses a cosine distance algorithm to determine the language of a given text. The training data for language profiles is sourced from Wikipedia articles.

The frontend is built with the "GoTTH" stack: **Go** on the backend, with **T**empl for HTML templating, **T**ailwindCSS for styling, and **H**tmx for interactivity (coming soon!).

## Prerequisites

- **Go**: Version 1.21 or later.
- **[templ](https://templ.dev/docs/installation)**: For HTML templating.
- **[Tailwind CSS Standalone CLI](https://tailwindcss.com/docs/installation/standalone-cli)**: For CSS generation.
- **[BadgerDB CLI](https://dgraph.io/docs/badger/get-started/#installing-badger)**: For managing the database.

## Getting Started

### 1. Clone the Repository

```bash
git clone https://github.com/myroslavrozum/go-langdetector.git
cd go-langdetector
```

### 2. Install Dependencies

Install the necessary Go modules and tools.

```bash
# Download Go application and tool dependencies
go mod tidy
go install github.com/a-h/templ/cmd/templ@latest
go install github.com/dgraph-io/badger/v4/badger@latest

# Verify installations (optional)
templ --version
badger --version
```

For Tailwind CSS, follow the [official guide](https://tailwindcss.com/docs/installation/standalone-cli) for your operating system. For macOS with Homebrew, you can run:

```bash
brew install tailwindcss
```

### 3. Run the Application

To run the web application, you need to have `templ` and `tailwindcss` running in watch mode to automatically recompile your templates and styles as you make changes.

Open three separate terminal windows in the project root directory and run the following commands:

```bash
# Terminal 1: Run the templ watcher
templ generate --watch

# Terminal 2: Run the Tailwind CSS watcher
tailwindcss -i ./webapp/assets/css/styles_template.css -o ./webapp/assets/css/styles_compiled.css --watch

# Terminal 3: Run the Go web server
go run .
```

## Configuring gears necessary for UI

go get -u github.com/gorilla/websocket

### GoTTH

#### First 'T' for Templ

```bash
go get -tool github.com/a-h/templ/cmd/templ@latest
go tool templ generate
```

#### Second 'T' for Tailwind

Install "Tailwind CSS IntelliSense" plugin for VSCode
For standalone Tailwind, without Node.JS follwo steps on
https://github.com/Aureuma/homebrew-tailwindcss

```bash
brew install aureuma/tailwindcss/tailwindcss-standalone 
```

For generic installation follow:
https://tailwindcss.com/docs/installation/tailwind-cli
add styles_compiled.css to .gitignore

```bash
tailwindcss -i webapp/assets/css/styles_compiled.css -o webapp/assets/css/styles.css
```

To run it with VSCode open two additional termainal windows and in each one run:

```bash
go tool templ generate -watch
tailwindcss -i webapp/assets/css/styles_template.css -o webapp/assets/css/styles_compiled.css --watch
```

respectively. THen you can use VSCode's `Run->...`
