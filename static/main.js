window.addEventListener( "load", function () {
  const form = document.getElementById("register");

  form.addEventListener("submit", (e) => {
    e.preventDefault();

    const xhr = new XMLHttpRequest();

    xhr.onreadystatechange = function() {
      if(xhr.readyState === 4) {
        // console.log(xhr.response)
        
        const responeSpan = document.querySelector("#server-response");

        const classes = {
          "417" : "server-response-417",
          "201" : "server-response-201",
          "401" : "server-response-401",
          "409" : "server-response-409",
          "418" : "server-response-418",
          "default" : "server-response"
        };

        for (const [_, value] of Object.entries(classes)) {
          responeSpan.classList.remove(value);
        }

        console.log(xhr.responseText);
        console.log(xhr.status);

        switch (xhr.status) {
          case 417:
            responeSpan.classList.add(classes["417"])
            responeSpan.innerHTML = "Ten nick powinien już mieć permisje. Jeśli dalej nie działa pytaj admina na discord czy tam gdzieś indziej.";
          break;

          case 201:
            responeSpan.classList.add(classes["201"])
            responeSpan.innerHTML = "Zakutalizowano permisje!"
          break;

          case 401:
            responeSpan.classList.add(classes["401"])
            responeSpan.innerHTML = "Złe hasło!"
          break;

          case 409:
            responeSpan.classList.add(classes["409"])
          responeSpan.innerHTML = "Nie ma cie w bazie danych LuckPerms wejdź na serwer i spróbuj ponownie";
          break;
          
          case 418:
            responeSpan.classList.add(classes["418"])
            responeSpan.innerHTML = "Co ty robisz mi tu xD"
          break;
        }
      }
    }
  

    xhr.open("POST", "/usrauth");
    xhr.send(new FormData(form));

    form.reset();
  })
  
});