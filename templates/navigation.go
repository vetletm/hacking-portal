package templates

// Navigation HTML template
const Navigation = `
{{define "navigation"}}
<nav class='navbar navbar-fixed-top bg-secondary'>
	<div class='container-fluid'>
		<div class='navbar-header w-100'>
			<img class='float-left' src='/static/images/ntnu-logo.svg'/>
			<h4 class='float-left pl-3 text-white'>Hacking Portal</h4>
			<div class='float-right pr-3 dropdown' style='cursor: pointer;'>
				<a class='dropdown-toggle' role='button' id='user-menu' data-toggle='dropdown' aria-haspopup='true' aria-expanded='false'>
					<img class='rounded' src='/static/images/stock-avatar.png'/>
				</a>
				<div class='dropdown-menu' aria-labelledby='user-menu'>
					<h5 class='dropdown-header'>{{.User.Name}}</h5>
					{{if ne .User.GroupID 0}}
					<a class='dropdown-item leave' href='/group/leave'>Leave Group</a>
					{{end}}
					<a class='dropdown-item' href='/logout'>Log out</a>
				</div>
			</div>
		</div>
	</div>
</nav>
{{end}}
`
