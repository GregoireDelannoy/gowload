/*
filedrag.js - HTML5 File Drag & Drop demonstration
Featured on SitePoint.com
Developed by Craig Buckler (@craigbuckler) of OptimalWorks.net
*/
(function() {
	// getElementById
	function $id(id) {
		return document.getElementById(id);
	}

	function humanFileSize(size) {
    var i = size == 0 ? 0 : Math.floor( Math.log(size) / Math.log(1024) );
    return ( size / Math.pow(1024, i) ).toFixed(2) * 1 + ' ' + ['B', 'kB', 'MB', 'GB', 'TB'][i];
	};

	// file selection
	function FileSelectHandler(e) {
		// fetch FileList object
		var files = e.target.files || e.dataTransfer.files;

		// process all File objects
		for (var i = 0, f; f = files[i]; i++) {
			UploadFile(f);
		}
	}

	function UploadFile(file) {
		var xhr = new XMLHttpRequest();
		if (xhr.upload) {
			// create progress bar
			var o = $id("progress");
			var progress = o.appendChild(document.createElement("p"));
			progress.appendChild(document.createTextNode(file.name + " (" + humanFileSize(file.size) + ")"));

			// progress bar
			xhr.upload.addEventListener("progress", function(e) {
				var pc = parseInt(100 - (e.loaded / e.total * 100));
				progress.style.backgroundPosition = pc + "% 0";
			}, false);

			xhr.upload.addEventListener("error", function(e){
				progress.style.backgroundColor = 'red';
				alert("upload error for " + file.name)
				console.dir(e)
			})

			// file received/failed
			xhr.onreadystatechange = function(e) {
				if (xhr.readyState == 4) {
					console.log("XHR finished: " + xhr.status)
					progress.className = (xhr.status == 200 ? "success" : "failed");
					if (xhr.status != 200) {
						alert("upload error for " + file.name + " : " + xhr.responseText)
					}
				}
			};

			// start upload
			xhr.open("POST", $id("upload").action, true);
			var formData = new FormData();
			formData.append("notaname", file);
			xhr.send(formData);
		}
	}

	// initialize
	function Init() {
		var fileselect = $id("fileselect"),
			submitbutton = $id("submitbutton");

		// file select
		fileselect.addEventListener("change", FileSelectHandler, false);

		// is XHR2 available?
		var xhr = new XMLHttpRequest();
		if (xhr.upload) {
			// remove submit button
			submitbutton.style.display = "none";
		}
	}

	// call initialization file
	if (window.File && window.FileList && window.FileReader) {
		Init();
	}
})();