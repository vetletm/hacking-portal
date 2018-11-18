package templates

// Login HTML template
const Login = `
{{define "title"}}Login{{end}}

{{define "navigation"}}{{end}}

{{define "body"}}
<div class='container w-50'>
	<header class='m-5'>
		<h1>Hacking Portal</h1>
	</header>

	<div class='col-md-5'>
		<form id='form' onsubmit='javascript:;'>
			<p>Use your NTNU credentials</p>
			<div class='form-group'>
				<input class='form-control' id='username' placeholder='Username' name='username' autofocus>
			</div>
			<div class='form-group'>
				<input class='form-control' id='password' placeholder='Password' name='password' type='password'>
			</div>
			<button type='submit' class='btn btn-sm btn-success float-right'>Login</button>
		</form>
	</div>
</div>
{{end}}

{{define "scripts"}}
<script type="text/javascript">
	$('button[type="submit"]').click(function(e){
		e.preventDefault();

		var username = document.getElementById('username');
		var password = document.getElementById('password');

		if(username.value == ""){
			username.setCustomValidity("Use your NTNU username");
			username.reportValidity();
			return;
		} else
			username.setCustomValidity("");

		if(password.value == ""){
			password.setCustomValidity("This is required");
			password.reportValidity();
			return;
		} else
			password.setCustomValidity("");

		$.ajax({
			type: 'POST',
			url: '/login',
			data: JSON.stringify({username: username.value, password: password.value}),
			contentType: 'application/json; charset=UTF-8'
		}).done(function(){
			window.location.pathname = '/groups'
		}).fail(function(){
			password.setCustomValidity("Username or password was incorrect");
			password.reportValidity();
		});
	});
</script>
{{end}}
`
