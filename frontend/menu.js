document.body.insertAdjacentHTML('afterbegin', `
    <button id="menuBtn">☰</button>
    <div id="sidebar" class="sidebar">
        <div class="sidebar-header">
            <h2>School Panel</h2>
            <p>Management System</p>
        </div>
        <ul>
            <li onclick="window.location.href='dashboard.html'">
                <span class="icon">⊞</span> Dashboard
            </li>
            <li onclick="window.location.href='class.html'">
                <span class="icon">◫</span> Classes
            </li>
            <li onclick="window.location.href='form.html'">
                <span class="icon">✦</span> Students
            </li>
        </ul>
    </div>
    <div id="overlay"></div>
`)

document.getElementById('menuBtn').addEventListener('click', () => {
    document.getElementById('sidebar').classList.toggle('open')
    document.getElementById('overlay').classList.toggle('show')
})

document.getElementById('overlay').addEventListener('click', () => {
    document.getElementById('sidebar').classList.remove('open')
    document.getElementById('overlay').classList.remove('show')
})