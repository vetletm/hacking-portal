package templates

// Groups HTML template
const Groups = `
{{define "title"}}Groups{{end}}

{{define "body"}}
<div class='container w-50'>
	<header class='m-5'>
		<h1>Groups</h1>
		<p>You are not currently in a group, to use the portal please join one.</p>
	</header>

	{{range .Groups}}
	<div class='card mb-3' data-group-id='{{.ID}}'>
		<div class='card-header'>
			<h5 class='float-left'>Group {{.ID}}</h5>
			{{if .Full}}
			<button type='button' class='btn btn-secondary btn-sm float-right join disabled'>Full</button>
			{{else}}
			<button type='button' class='btn btn-primary btn-sm float-right join'>Join</button>
			{{end}}
		</div>
		<div class='card-body p-2'>
			{{if .Members}}
				<h6 class='card-text m-1'>Members:</h6>
				<ul>
				{{range .Members}}
					{{if .Name}}
					<li class='card-text m-1'>{{.Name}}</li>
					{{else}}
					<li class='card-text m-1'>{{.ID}}</li>
					{{end}}
				{{end}}
				</ul>
			{{else}}
				<h6 class='card-text m-1'>No members</h6>
			{{end}}
		</div>
	</div>
	{{end}}
</div>
{{end}}

{{define "scripts"}}
<script type="text/javascript">
	$('button.join').click(function(){
		var button = $(this);
		var groupID = button.parent().parent().data('group-id');

		// prevent clicking other buttons
		if($(this).hasClass('disabled'))
			return;
		$('button.join').addClass('disabled');

		// reset button color
		button.removeClass('btn-danger');
		button.addClass('btn-secondary');

		// add some visual response that we're joining a group
		button.text('Joining');
		button.append('<i class="fa fa-spinner fa-spin ml-2"></i>');

		// attempt to join group
		$.ajax({
			type: 'POST',
			url: '/groups/join',
			data: JSON.stringify({groupID: groupID}),
			contentType: 'application/json; charset=UTF-8'
		}).done(function(){
			// redirect to group page
			setTimeout(function(){ // sleep for demo purposes
				window.location.pathname = '/group';
			}, 2500);
		}).fail(function(){
			// update the button text
			button.text('Join');

			// give visual indicator that it failed
			button.addClass('btn-danger');
			button.removeClass('btn-secondary');
		}).always(function(){
			// unlock the buttons
			$('.join:not(.btn-secondary)').removeClass('disabled');
		});
	});
</script>
{{end}}
`
