# BUG-017: "Train" Navigation Link Vanishes After Language Detection

- **Status:** Resolved
- **Severity:** Low
- **Component:** `webapp`
- **Affects:** [`webapp/index.templ:71`](file:///Users/myroslavrozum/src/go-langdetector/webapp/index.templ#L71), [`webapp/webhandlers.go:66`](file:///Users/myroslavrozum/src/go-langdetector/webapp/webhandlers.go#L66)

## Description
In `webapp/index.templ`, the `#detected` list initially renders the list of languages and a trailing "Train" action link:
```templ
<ul id="detected" class="flex flex-col gap-2.5">
	if len(model.SupportedLanguages) > 0 {	
		@templ.Raw(model.SupportedLanguagesRendered)
	} else {
		@errorMsg("Empty model")
	}
	<li><a href="#" class="text-[#3b5998] hover:underline text-sm">Train</a></li>
</ul>
```
However, the detect form specifies:
```templ
<form hx-post="/detect" hx-target="#detected" hx-swap="innerHTML">
```
When `/detect` finishes, it returns `[]byte(m.SupportedLanguagesRendered)`, which contains only the language `<li>` elements. HTMX replaces the inner HTML of `#detected` with this payload.

## Impact
Submitting any detection request completely removes the "Train" link from the sidebar.

## Steps to Reproduce
1. Open the homepage; observe the "Train" link at the bottom of the language list.
2. Enter text and click "Detect Language".
3. Notice that the "Train" link has vanished from the DOM.

## Suggested Fix
Include the "Train" link inside `SupportedLanguagesRendered`, or adjust the HTMX swap target to update only the language list elements without overwriting sibling controls.
