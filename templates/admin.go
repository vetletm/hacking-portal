package templates

// Admin HTML template
const Admin = `
{{define "title"}}Administration{{end}}

{{define "body"}}
{{$machines := .Machines}}
{{$groups := .Groups}}
<div class='container w-50'>
	<header class='m-5'>
		<h1>Administration</h1>
	</header>

	<div class='d-flex flex-column'>
		<div class='row-md-4'>
			<h3>Assign Kali machines to groups</h3>
			{{if .Machines}}
			<div class='list-group pb-2'>
				{{range $machine := $machines}}
				<div class='list-group-item clearfix' id='{{$machine.UUID}}'>
					<span class='float-left'>{{$machine.Name}}</span>
					<a class='float-left pl-3' href='#'>{{$machine.Address}}</a>
					<span class='float-left pl-3'>-></span>
					<div class='float-left pl-4 dropdown'>
						<button class='btn btn-sm btn-light dropdown-toggle' data-toggle='dropdown' aria-haspopup='true' aria-expanded='false'>
							{{if eq $machine.GroupID 0}}
							None
							{{else}}
							Group {{$machine.GroupID}}
							{{end}}
						</button>
						<div class='dropdown-menu groups'>
							{{if eq $machine.GroupID 0}}
							<a class='dropdown-item disabled' href='#' data-id='0'>None</a>
							{{else}}
							<a class='dropdown-item' href='#' data-id='0'>None</a>
							{{end}}
							{{range $group := $groups}}
							{{if eq $machine.GroupID $group.ID}}
							<a class='dropdown-item disabled' href='#' data-id='{{$group.ID}}'>Group {{$group.ID}}</a>
							{{else}}
							<a class='dropdown-item' href='#' data-id='{{$group.ID}}'>Group {{$group.ID}}</a>
							{{end}}
							{{end}}
						</div>
					</div>
					<span class='float-right' data-kali-index='{{.GroupIndex}}'>
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
	$('.dropdown-menu.groups a').click(function(){
		var uuid = $(this).parent().parent().parent().attr('id');
		var button = $('#' + uuid + ' .dropdown-toggle');
		var groupID = $(this).data('id');
		var groupName = $(this).text();
		var groupNameOld = button.text();

		// prevent spamming
		if(button.hasClass('disabled'))
			return;
		button.addClass('disabled');

		// reset button color
		button.removeClass('btn-outline-danger');
		button.addClass('btn-light');

		// replace button text with spinner
		button.html('<i class="fa fa-spinner fa-spin"></i>');

		// attempt to store the change
		$.ajax({
			type: 'POST',
			url: '/admin/assign',
			data: JSON.stringify({groupID: groupID, machineUUID: uuid}),
			contentType: 'application/json; charset=UTF-8'
		}).done(function(){
			// update the button text
			button.text(groupName);

			// enable the menu buttons
			$(this).siblings().removeClass('disabled');
			$(this).addClass('disabled');
		}).fail(function(){
			// update the button text
			button.text(groupNameOld);

			// give visual indicator that it failed
			button.addClass('btn-outline-danger');
			button.removeClass('btn-light');
		}).always(function(){
			// unlock the button
			button.removeClass('disabled');
		});
	});
	$('a.restart').click(function(){
		if(confirm('This will forcefully restart the machine, losing all unsaved progress. Are you sure you want to do this?')){
			var button = $(this);
			var uuid = $(this).parent().parent().attr('id');

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
				url: '/admin/restart/' + uuid
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
