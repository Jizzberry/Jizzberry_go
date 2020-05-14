var Container = document.getElementById("navbarSupportedContent");
var btns = Container.getElementsByClassName("nav-item");
for (var i = 0; i < btns.length; i++) {
 btns[i].addEventListener("click", function() {
   var current = Container.getElementsByClassName("active");
   current[0].className = current[0].className.replace("active", "");
   this.className += " active";
 });
}

