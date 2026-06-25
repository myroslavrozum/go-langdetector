go get -u github.com/gorilla/websocket

# GoTTH
## First 'T' for Templ
go get -tool github.com/a-h/templ/cmd/templ@latest

```
go tool templ generate
```
## Second 'T' for Tailwind
Install "Tailwind CSS IntelliSense" plugin for VSCode
For standalone Tailwind, without Node.JS follwo steps on
https://github.com/Aureuma/homebrew-tailwindcss
```
brew install aureuma/tailwindcss/tailwindcss-standalone 
```
For generic installation follow:
https://tailwindcss.com/docs/installation/tailwind-cli

```
add styles_compiled.css to .gitignore
tailwindcss -i webapp/assets/css/styles_compiled.css -o webapp/assets/css/styles.css


To run it with VSCode open two additional termainal windows and in each one run:
```
go tool templ generate -watch
tailwindcss -i webapp/assets/css/styles_template.css -o webapp/assets/css/styles_compiled.css --watch
```
respectively. THen you can use VSCode's `Run->Stert...`
