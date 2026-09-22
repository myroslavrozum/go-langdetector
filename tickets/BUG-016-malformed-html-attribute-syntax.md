# BUG-016: Malformed HTML Attribute Syntax in `renderSupportedLanguages`

- **Status:** Resolved
- **Severity:** Medium
- **Component:** `webapp`
- **Affects:** [`webapp/model.go:40`](file:///Users/myroslavrozum/src/go-langdetector/webapp/model.go#L40)

## Description
In `webapp/model.go`:
```go
m.SupportedLanguagesRendered += fmt.Sprintf(`<li id="langmarker_" + %s class="%s"><a href="%s">%s</a></li>`, shortName, class, trainingUrl, fullName)
```
The format string contains `" + %s` inside the `id` attribute value rather than formatting `%s` into the attribute name.

## Impact
This produces malformed HTML elements:
```html
<li id="langmarker_" + en class="text-[#3b5998] hover:underline text-sm"><a href="...">English</a></li>
```
In standard HTML parsers, this renders an element with `id="langmarker_"`, followed by boolean attributes `+` and `en`. Any JavaScript or CSS looking for `#langmarker_en` fails to find the element.

## Steps to Reproduce
Inspect the DOM elements inside `#detected` in the browser dev tools.

## Suggested Fix
Fix the format string:
```go
m.SupportedLanguagesRendered += fmt.Sprintf(`<li id="langmarker_%s" class="%s"><a href="%s">%s</a></li>`, shortName, class, trainingUrl, fullName)
```
