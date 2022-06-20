# l'oGGo: Rich Terminal User Interface Logging App

## Introduction

*For the impatient, go to [Getting Started](#getting-started)*
<p align="center">
<img src="img/loggo_sm.png">
</p>
l'oGGo or Log & Go is a rich Terminal User Interface app written in [golang](https://go.dev/) that harness the
power of your terminal to digest JSON based logs. This project is a hobby project
and is by no means bulletproof, but should be stable enough for every-day
troubleshooting workflows.

It came to light as JSON based logs and applications slowly drifted 
to become the de-facto standard for logging across applications and platforms. Although JSON data
structure provided a sound and well-behaved data model, the lack of local tools
to aid streaming & rendering for realtime troubleshooting such verbosely-rich 
produced payloads motivated me to embark in this endevour as I was, for a little
while, no longer able to quickly cast eyes on logs and pinpoint hotspots.

<img src="img/compare.png">
<table>
<tr>
<td>
<p>Without l`oGGo</p>
<img src="mov/term.gif">
</td>
<td>
<p>With l`oGGo</p>
<img src="mov/loggo.gif">
</td>
</tr>
</table>

Loggo App leveraged [tview](https://github.com/rivo/tview/) and [tcell](https://github.com/gdamore/tcell) projects for rich Terminal User 
Interface (TUI).

## Getting Started

### macOS Systems:
The easiest way is to utilise [Homebrew](https://brew.sh/) package management system. Once 
installed simply issue the following command:

````
brew tap aurc/loggo
brew install aurc/loggo/loggo
````

### Build from Source:
Including **macOS**, build from source. 
Pre-Reqs:
- [Golang](https://go.dev/) v1.8+

````
go build -o loggo
````

## Using l'oGGo

Loggo can be used to stream parsed logs from a persisted file and from a 
piped input and also provides a tool for creating log templates.

### `help` Command

To gain fine grained insight of each `loggo` command params, use
the `help` command, e.g.:
````
loggo help
loggo stream help
loggo template help
````

### `stream` Command

![](img/loggo_log.png)

**From File:**
````
loggo stream --file <my file>
````
*With Template:*
````
loggo stream --file <my file> --template <my template yaml>
````

**From Pipe:**
````
tail -f <my file> | loggo stream
````
*With Template:*
````
tail -f <my file> | loggo stream --template <my template yaml>
````

Note that you can pipe to anything that produces an output to the `stdin`.

**Some Additional Features:**
- Drill down onto each log entry
![](img/log_entry.png)
- Configure Rendering Templates:
![](img/render_template.png)
- Fine Tune how columns are displayed (Template):
![](img/how_to_display.png)
    Note that single Value Matches are REGEX expressions.

### `template` Command
The template command opens up the template editor without the
need to stream logs. This is convenient if you want to craft
templates prior using the loggo command.

**Blank Canvas:**

![](mov/template.gif)
````
loggo template
````

**Edit Existing Template:**
````
loggo template --file <my template yaml>
````


## Current Limitations

Most of the items listed here are slated for development in the near future,
prior the first release.
- Search log entry.
- Filter log by json key(s).
- Copy single log entry to clipboard.
- Browse/Load new log templates on the fly.
- Create template with keys whose name contains `/` as it uses slashes to navigate to nested json branches.
