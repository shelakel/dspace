package main

import "text/template"

var tmpl = template.Must(template.New("default").Parse(`
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="utf-8">
	<title>Disk Space Visualizer</title>
	<style type="text/css">
		html, body { font: 14px/20px "Courier New", Courier, monospace }
		ul { list-style: none; padding: 0 0 5px 10px }
		li { cursor: pointer }
		span { display: block; padding: 3px }
		li:hover > span { background:#efefef; }
	</style>
</head>
<body>
<h1>Disk Space Analyzer</h1>
<ul id="root">
</ul>
<script src="https://ajax.googleapis.com/ajax/libs/jquery/1.11.2/jquery.min.js"></script>
<script id="data" type="application/json">
{{.}}
</script>
<script>
	var i = 0;
	function render(data) {
		var nodeBuilder = [
			'<li class="folder"><span title="',
			data.fullPath,
			'" data-for="#p', i++, '">&nbsp;',
			data.subDirs && data.subDirs.length > 0 ? ' <b>+</b>' : '&nbsp;&nbsp;',
			'&nbsp;',
			data.path, ' (', data.size, ')',
			'</span>'];
		if(data.subDirs && data.subDirs.length > 0) {
			nodeBuilder.push(['<ul id="p', i-1,'" style="display:none">'].join(''));
			for(var j = 0; j < data.subDirs.length; j++) {
				nodeBuilder.push(render(data.subDirs[j]));
			}
			nodeBuilder.push('</ul>');
		}
		nodeBuilder.push('</li>');
		return nodeBuilder.join('');
	}

	$(function() {
		data = JSON.parse(document.getElementById('data').innerHTML);
		$('#root').append($(render(data)));
		$('.folder > span').on('click', function() {
			var t = $(this);
			var b = t.find('b');
			b && b.text(b.text() == '+' ? '-' : '+');
			var f = $(t.data('for'));
			f && f.slideToggle('fast');
		});
	})
</script>
</body>
</html>`))
