package main

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var tableHTML = `
{{define "table"}}
	<table id="res-list">
			<thead>
					<tr>
							<th>Name</th>
							<th>Active</th>
					</tr>
			</thead>
			<tbody>
					{{range .data}}
					<tr>
							<td>{{.Name}}</td>
							<td>{{.Active}}</td>
					</tr>
					{{end}}
			</tbody>
	</table>
{{end}}
{{if .tableOnly}}
	{{template "table" .}}
{{else}}
<html>
<head>
	<title>Example table</title>
</head>
<body>
	<form id="filters" action="/" method="GET" hx-get="/" hx-target="#res-list" hx-include="#filters" hx-push-url="true" hx-swap="outerHTML">
		<div>
			<input name="filter.name" type="text" placeholder="Text input" {{if .filter.Name}}value="{{.filter.Name}}"{{end}}>
		</div>
		<div>
			<label>Active</label>
			<label>
				<input type="radio" name="filter.active" value="true" {{if eq .filter.Active "true"}}checked{{end}}>
				Yes
			</label>
			<label>
				<input type="radio" name="filter.active" value="false" {{if eq .filter.Active "false"}}checked{{end}}>
				No
			</label>
			<label>
				<input type="radio" name="filter.active" value="" {{if eq .filter.Active "" "checked"}}checked{{end}}>
				Any
			</label>
		</div>
		<div>
			<a href="/">Reset filters</a>
			<button type="submit">Search</button>
		</div>
	</form>
	<p>Built: {{.built}}</p>
	{{template "table" .}}
	<script src="https://unpkg.com/htmx.org@1.0.2"></script>
</body>
</html>
{{end}}
`

type Data struct {
	Name   string
	Active bool
}

var data = []Data{
	{"Harry", false},
	{"Sue", true},
	{"Fred", true},
	{"Jane", false},
}

type ListFilter struct {
	Name   string
	Active string
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var l ListFilter
		t, err := template.New("T").Parse(tableHTML)
		if err != nil {
			panic(err)
		}

		l.Name = r.URL.Query().Get("filter.name")
		l.Active = r.URL.Query().Get("filter.active")

		payload := map[string]interface{}{
			"data":   filterData(data, l),
			"filter": l,
			"built":  time.Now().Format("2006-01-02 15:04:05"),
		}

		// Looks like a HTMX request and matches our trigger name, so let's proceed
		if r.Header.Get("HX-Request") == "true" && r.Header.Get("HX-Trigger") == "filters" {
			payload["tableOnly"] = true
		}

		err = t.ExecuteTemplate(w, "T", payload)
		if err != nil {
			panic(err)
		}
	})

	log.Printf("Starting on 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func filterData(data []Data, filter ListFilter) (res []Data) {
	active, _ := strconv.ParseBool(filter.Active)
	name := strings.ToLower(filter.Name)
	for _, datum := range data {
		if (len(name) == 0 || strings.Contains(strings.ToLower(datum.Name), name)) &&
			(len(filter.Active) == 0 || (datum.Active == active)) {
			res = append(res, datum)
		}
	}

	return
}
