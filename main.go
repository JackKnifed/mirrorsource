package main

import "html/template"

var templates = template.Must(template.ParseFiles("templates.html"))
