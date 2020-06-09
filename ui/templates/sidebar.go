package templates

const Sidebar = `
{{ define "sidebar" }}
<nav class="nav flex-column nav-pills shadow sidebar bg-light border" aria-orientation="vertical">
    <a class="nav-link active" href="/dashboard" aria-controls="v-pills-dashboard" aria-selected="true">Dashboard</a>
    <a class="nav-link" href="#v-pills-profile" aria-controls="v-pills-profile" aria-selected="false">Profile</a>
    <a class="nav-link" href="#v-pills-messages" aria-controls="v-pills-messages" aria-selected="false">Messages</a>
    <a class="nav-link" href="#v-pills-settings" aria-controls="v-pills-settings" aria-selected="false">Settings</a>
</nav>
{{ end }}
`