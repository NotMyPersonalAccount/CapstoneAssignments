const body = document.querySelector("body");
const nav = document.createElement("nav");
nav.className = "navbar navbar-expand-lg navbar-dark bg-dark";
nav.innerHTML = `<div class="container-fluid">
  <a class="navbar-brand" href="3.02-Bootstrap.html">Research</a>
  <button
    class="navbar-toggler"
    type="button"
    data-bs-toggle="collapse"
    data-bs-target="#navbarNavDropdown"
    aria-controls="navbarNavDropdown"
    aria-expanded="false"
    aria-label="Toggle navigation"
  >
    <span class="navbar-toggler-icon"></span>
  </button>
  <div class="collapse navbar-collapse" id="navbarNavDropdown">
    <ul class="navbar-nav">
      <li class="nav-item">
        <a
          class="nav-link active"
          aria-current="page"
          href="Research.html"
          >Home</a
        >
      </li>
      <li class="nav-item">
        <a class="nav-link active" aria-current="page" href="#"
          >Initial Research</a
        >
      </li>
      <li class="nav-item">
        <a
          class="nav-link active"
          aria-current="page"
          href="RootCauses.html"
          >Root Causes</a
        >
      </li>
      <li class="nav-item">
        <a
          class="nav-link active"
          aria-current="page"
          href="CaseStudy.html"
          >Case Study</a
        >
      </li>
    </ul>
  </div>
</div>`;
body.prepend(nav);
