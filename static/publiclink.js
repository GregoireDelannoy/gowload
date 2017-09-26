(function() {
  var submitHandler = function(event){
  event.preventDefault();
    // actual logic, e.g. validate the form
    console.log('Form submission cancelled.');
    console.log(event)
    var xhr = new XMLHttpRequest();
    xhr.onreadystatechange = function(e) {
      if (xhr.readyState == 4) {
        console.log("XHR finished: " + xhr.status)
        if (xhr.status != 200) {
          alert("Error getting link " + event.target.action + " : " + xhr.responseText)
        } else {
          event.target.parentElement.innerHTML = window.location.protocol + "//" + window.location.host + xhr.responseText
        }
      }
    };

    // start upload
    xhr.open("POST", event.target.action, true);
    var formData = new FormData(event.target);
    xhr.setRequestHeader("Content-type","application/x-www-form-urlencoded");
    xhr.send("actionType=link");
  }

  var forms = document.getElementsByClassName('publicLinkForm');
  for (var i = 0; i < forms.length; i++) {
    forms[i].addEventListener('submit', submitHandler)
  }
})();