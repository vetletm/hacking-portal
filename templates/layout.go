package templates

// Layout HTML template
const Layout = `
{{define "layout"}}
<!DOCTYPE html>
<html>
	<head>
		<meta charset='utf-8'>
		<meta name='viewport' content='width=device-width, initial-scale=1, shrink-to-fit=no'>
		<title>{{template "title" .}}</title>

		<link rel='stylesheet' href='/static/css/libs/bootstrap.min.css'>
		<link rel='stylesheet' href='/static/css/libs/fontawesome.all.css'>
		<link rel='stylesheet' href='/static/css/tweaks.css'>
	</head>
	<body>
		{{template "navigation" .}}
		{{template "body" .}}

		<script src='/static/js/libs/jquery.min.js'></script>
		<script src='/static/js/libs/popper.min.js'></script>
		<script src='/static/js/libs/bootstrap.min.js'></script>
		{{template "scripts"}}
	</body>
</html>
{{end}}
`
