document.addEventListener("click", function (event) {
  const sidebar = document.getElementById("sidebar");

  if (!sidebar) return;

  if (
    sidebar.classList.contains("open") &&
    !sidebar.contains(event.target) &&
    !event.target.closest(".sidebar-toggle")
  ) {
    sidebar.classList.remove("open");
  }
});