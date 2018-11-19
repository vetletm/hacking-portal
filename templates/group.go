package templates

// Group HTML template
const Group = `
{{define "title"}}Group {{.User.GroupID}}{{end}}

{{define "body"}}
{{$machines := .Machines}}
<div class='container w-50'>
	<header class='m-5'>
		<h1>Group {{.User.GroupID}}</h1>
	</header>

	<div class='d-flex flex-column'>
		<div class='row-md-4'>
			<h3>Kali Machines</h3>
			{{if .Machines}}
			<div class='list-group pb-2'>
			{{range $i, $machine := $machines}}
				<div class='list-group-item clearfix'>
					<span class='float-left'>Kali {{inc $i}} <a class='d-inline pl-3' href='#'>{{$machine.Address}}</a></span>
					<span class='float-right' data-kali-uuid='{{$machine.UUID}}'>
						<a class='btn btn-sm btn-outline-default border restart' href='#' data-toggle='tooltip' title='Restart'>
							<span class='fas fa-redo' aria-hidden='true'></span>
						</a>
					</span>
				</div>
			{{end}}
			</div>
			{{end}}
		</div>
	</div>
</div>
{{end}}

{{define "scripts"}}
<script type="text/javascript">
	$('[data-toggle="tooltip"]').tooltip() // enable tooltips
	$('a.restart').click(function(){
		if(confirm('This will forcefully restart the machine, losing all unsaved progress. Are you sure you want to do this?')){
			var button = $(this);
			var machineUUID = button.parent().data('kali-uuid');

			// prevent spamming
			if(button.hasClass('disabled'))
				return;
			button.addClass('disabled');

			// reset button color
			button.removeClass('btn-outline-danger');
			button.addClass('btn-outline-default');

			// attempt to restart the machine
			$.ajax({
				type: 'POST',
				url: '/group/restart/' + machineUUID
			}).fail(function(){
				// give visual indicator that it failed
				button.addClass('btn-outline-danger');
				button.removeClass('btn-outline-default');
			});

			// lock and spin the button for one minute
			button.children().addClass('fa-spin');
			setTimeout(function(){
				button.removeClass('disabled');
				button.children().removeClass('fa-spin');
			}, 60 * 1000);
		}
	});
</script>
{{end}}
`
