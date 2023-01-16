# MDP

## Abstract
A basic markdown preview tool that converts markdown
source into HTML that can be viewed in a browser.

## High level sequence of steps
1. Read the contents of the input Markdown file.
2. Use some Go external libraries ot parse Markdown and generate a 
valid HTML block.
3. Wrap the results with an HTML  header and footer.
4. Save the buffer to an HTML file that you can view in a browser.

## Features
- Support for links

## How to install
```
go get github.com/gregidonut/mpd
```