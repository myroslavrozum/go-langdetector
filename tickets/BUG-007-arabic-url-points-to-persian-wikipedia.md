# BUG-007: Arabic Training URL Points to Persian (Farsi) Wikipedia

- **Status:** Resolved
- **Severity:** High
- **Component:** `constants`
- **Affects:** [`constants/constants.go:13`](file:///Users/myroslavrozum/src/go-langdetector/constants/constants.go#L13)

## Description
In `constants/constants.go`, the entry for Arabic (`ar`) specifies a URL hosted on the Persian Wikipedia (`fa.wikipedia.org`):
```go
`ar`: {`Arabic`, `http://fa.wikipedia.org/wiki/%D9%88%DB%8C%DA%98%D9%87:%D8%B5%D9%81%D8%AD%D9%87_%D8%AA%D8%B5%D8%A7%D8%AF%D9%81%DB%8C`},
```
The URL path `%D9%88%DB%8C%DA%98%D9%87:...` translates to Persian "ویژه:صفحه_تصادفی" (Special:Random page in Persian).

## Impact
Arabic language profiles are trained on Persian text rather than Arabic. While Persian uses the Perso-Arabic script, Persian is an Indo-European language with distinct vocabulary, grammar, and extra characters (پ, چ, گ, etc.). Consequently, the detector will misclassify Arabic text and confuse it with Persian.

## Steps to Reproduce
1. Crawl `constants.UrlDictionary["ar"][1]`.
2. Inspect the retrieved text: it is written in Persian, not Arabic.

## Suggested Fix
Update the URL to point to Arabic Wikipedia:
```go
`ar`: {`Arabic`, `https://ar.wikipedia.org/wiki/Special:Random`},
```
(or `https://ar.wikipedia.org/wiki/%D8%AE%D8%A7%D8%B5:%D8%B9%D8%B4%D9%88%D8%A7%D8%A6%D9%8A%D8%A9`).
