window.addEventListener( "load", function () {
  const form = document.getElementById("register");

  form.addEventListener("submit", (e) => {
    e.preventDefault()

    const xhr = new XMLHttpRequest();

    xhr.onreadystatechange = function() {
      if(xhr.readyState == 4) {
        console.log(xhr.response)
      }
    }
  

    xhr.open("POST", "/usrauth");
    xhr.send(new FormData(form));

    form.reset();
  })
  
});