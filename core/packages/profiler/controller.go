package profiler

import (
	"bytes"
	"flamingo/core/flamingo/web"
	"html/template"
	"net/http"
	"strconv"
	"time"
)

type (
	// ProfileController shows information about a requested profile
	ProfileController struct{}
)

const profileTemplate = `<!doctype html>
<html lang="en">
<head>
	<title>Profile {{.Context.ID}}</title>
	<link href="https://fonts.googleapis.com/css?family=Roboto" rel="stylesheet">
	<style>
	body {
		background: #475157;
		margin: 0;
		padding: 25px;
		font-family: 'Roboto', sans-serif;
		font-size: 14px;
	}

	.profiler {
		background: #fff;
		padding: 25px;
		max-width: 1000px;
		margin: 0 auto;
	}

	.profiler-summary {
		list-style: none;
		margin: 0;
		padding: 20px;
		background: #F79223;
		color: #fff;
		font-size: 18px;
	}

	.profiler-entries {
		list-style: none;
		margin: 0;
		padding: 0;
	}

	.profiler-entry {
		overflow: hidden;
		padding: 10px 0 5px 10px;
		margin-bottom: 5px;
		background: rgba(50, 50, 50, 0.04);
		border-bottom: 1px solid rgba(50, 50, 50, 0.05);
	}

	.profiler-entry .fnc {
		display: block;
		color: rgab(0, 0, 0, 0.5);
		margin: 0 0 5px;
	}

	.profiler-entry .file {
		overflow: hidden;
		transition: all 0.3s ease-in-out;
		max-height: 0;
	}

	.profiler-subentries {
		margin-top: 10px;
		padding-left: 0px;
	}

	.profiler-entry .duration {
		float: right;
		font-weight: bold;
		font-size: 16px;
	}

	.profiler-entry .duration-relative {
		float: right;
		width: 100px;
		border: 1px solid #ccc;
		margin: 5px 15px;
	}

	.profiler-entry .duration-relative .inner {
		background: #F79223;
		display: block;
		height: 12px;
		width: 0;
	}

	.profiler-entry .msg {
		margin: 0;
		font-size: 16px;
	}

	.profiler-entry .fnc {
		color: rgba(0, 0, 0, 0.5);
	}

	.profiler-entry .fnc.has-file {
		cursor: pointer;
	}

	.profiler-entry .fnc.has-file .icon {
		display: inline-block;
		position: relative;
		top: 0.1em;
		width: 0.75em;
		height: 0.75em;
		border-radius: 2px;
		margin-right: 0.25em;
		border: 1px solid rgba(0, 0, 0, 0.5);
		font-style: normal;
		pointer-events: none;
	}

	.profiler-entry .fnc.has-file .icon:after {
		content: "\002B";
		position: absolute;
		top: -5px;
		left: 2px;
		font-size: 13px;
	}

	.profiler-entry.is-open > .fnc.has-file .icon:after {
		content: "\2013";
		left: 1px;
	}

	.profiler-entry .file-meta,
	.profiler-entry .file-wrap {
		overflow: hidden;
		transition: all 0.3s ease-in-out;
		max-height: 0;
	}

	.profiler-entry.is-open > .file-meta,
	.profiler-entry.is-open > .file-wrap {
		max-height: 500px;
		overflow-y: auto;
	}

	.profiler-entry .file-meta {
		background: #F79223;
		color: white;
		padding: 0 20px 0 10px;
	}

	.profiler-entry.is-open > .file-meta {
		padding: 5px 20px 5px 10px
	}

	.profiler-entry .file-hint {
		margin: 0 0 10px;
		padding: 10px;
		font-size: 12px;
		line-height: 1.15em;
		background: #fff;
	}

	.profiler-subentries {
		margin-top: 10px;
		padding-left: 0px;
	}
	</style>
</head>
<body>
<div class="profiler">
	<header class="profiler-header">
		<h1>Profile {{.Context.ID}}</h1>
		<ul class="profiler-summary">
			<li class="duration-total" data-duration="{{printf "%d" .Duration }}">Time: {{.Duration}}</li>
			<li>Start: {{.Start}}</li>
		</ul>
	</header>
	<div class="profiler-content">
		<ul class="profiler-entries">
		{{ range $entry := .Childs }}
			{{ template "entry" $entry }}
		{{ end }}
		</ul>
	</div>
</div>
<script>
	document.addEventListener('click', function(e) {
		if(e.target.classList.contains('fnc')) {
		e.target.parentNode.classList.toggle('is-open');
		}
	});

	var totalDuration = document.querySelector('.duration-total').dataset.duration;

	Array.from(document.querySelectorAll('.duration-relative')).forEach(addRelativeDuration);
	function addRelativeDuration(element) {
		var duration = element.dataset.duration;
		var relativeDuration = Math.min(Math.round(100 / totalDuration * duration), 100);
		element.querySelector('.inner').style.width = relativeDuration + '%';

	}
</script>
</body>
</html>

{{ define "entry" }}
<li class="profiler-entry">
	<span class="duration-relative" data-duration="{{printf "%d" .Duration }}"><i class="inner"></i></span>
	<span class="duration">{{ .Duration }}</span>
	<h3 class="msg">{{ .Msg }}</h3>
	<span class="fnc {{if and .Startpos .Endpos}}has-file{{end}}"><i class="icon"></i>{{ .Fnc }}</span>

	{{if and .Startpos .Endpos}}
		<div class="file-meta">
			<span class="file-path">{{.File}}</span>
			<span class="file-lines">Line {{ .Startpos }} - {{ .Endpos }}</span>
		</div>
		<div class="file-wrap">
			<pre class="file-hint">{{ .Filehint }}</pre>
		</div>
	{{ end }}

	{{if .Childs}}
		<ul class="profiler-subentries">
			{{ range $entry := .Childs }}
			{{ template "entry" $entry }}
			{{ end }}
		</ul>
	{{ end }}
</li>
{{ end }}
`

// Get Response for Debug Info
func (dc *ProfileController) Get(ctx web.Context) web.Response {
	t, err := template.New("tpl").Parse(profileTemplate)
	if err != nil {
		panic(err)
	}
	var body = new(bytes.Buffer)

	t.ExecuteTemplate(body, "tpl", profilestorage[ctx.Param1("profile")])

	return &web.ContentResponse{
		ContentType: "text/html; charset=utf-8",
		Status:      http.StatusOK,
		Body:        body,
	}
}

func (dc *ProfileController) Post(ctx web.Context) web.Response {
	dur, _ := strconv.ParseFloat(ctx.Form1("duration"), 64)
	profilestorage[ctx.Param1("profile")].ProfileOffline(ctx.Form1("key"), ctx.Form1("message"), time.Duration(dur*1000*1000))

	return &web.JSONResponse{}
}
