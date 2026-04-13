// ── Persistent Sidebar ───────────────────────────────────────
requireAuth()

document.body.insertAdjacentHTML('afterbegin', `
    <button id="menuBtn" title="Toggle menu">☰</button>
    <div id="sidebar" class="sidebar">
        <div class="sidebar-header">
            <h2>KRB School</h2>
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
            <li onclick="window.location.href='fees.html'">
                <span class="icon">₹</span> Fees
            </li>
            <li onclick="window.location.href='notices.html'">
                <span class="icon">📢</span> Notices
            </li>
            <li onclick="window.location.href='results.html'">
                <span class="icon">★</span> Results
            </li>
            <li onclick="window.location.href='papers.html'">
                <span class="icon">◫</span> Papers
            </li>
            <li onclick="window.location.href='reset-password.html'">
                <span class="icon">🔑</span> Reset Password
            </li>
            <li class="logout-item" onclick="logout()">
                <span class="icon">↩</span> Logout
            </li>
        </ul>
    </div>
    <div id="overlay"></div>
`)

const sidebar = document.getElementById('sidebar')
const menuBtn = document.getElementById('menuBtn')
const overlay = document.getElementById('overlay')

const isMobile = () => window.innerWidth <= 900

let isCollapsed = localStorage.getItem('sidebar_collapsed') === 'true'

function applyState() {
    if (isMobile()) {
        sidebar.classList.remove('collapsed')
        menuBtn.classList.remove('sidebar-collapsed')
        document.querySelectorAll('.page').forEach(p => p.classList.remove('sidebar-collapsed'))
        return
    }
    if (isCollapsed) {
        sidebar.classList.add('collapsed')
        sidebar.classList.remove('open')
        menuBtn.classList.add('sidebar-collapsed')
        document.querySelectorAll('.page').forEach(p => p.classList.add('sidebar-collapsed'))
    } else {
        sidebar.classList.remove('collapsed')
        menuBtn.classList.remove('sidebar-collapsed')
        document.querySelectorAll('.page').forEach(p => p.classList.remove('sidebar-collapsed'))
    }
}

function toggleSidebar() {
    if (isMobile()) {
        const isOpen = sidebar.classList.contains('open')
        sidebar.classList.toggle('open', !isOpen)
        overlay.classList.toggle('show', !isOpen)
    } else {
        isCollapsed = !isCollapsed
        localStorage.setItem('sidebar_collapsed', isCollapsed)
        applyState()
    }
}

function setActivePage() {
    const path = window.location.pathname
    const links = sidebar.querySelectorAll('li')
    links.forEach(li => li.classList.remove('active'))

    if (path.includes('dashboard')) links[0]?.classList.add('active')
    else if (path.includes('class')) links[1]?.classList.add('active')
    else if (path.includes('form') || path.includes('student')) links[2]?.classList.add('active')
    else if (path.includes('fee')) links[3]?.classList.add('active')
    else if (path.includes('notices')) links[4]?.classList.add('active')
    else if (path.includes('results')) links[5]?.classList.add('active')
    else if (path.includes('papers')) links[6]?.classList.add('active')
    else if (path.includes('reset-password')) links[7]?.classList.add('active')
}

menuBtn.addEventListener('click', toggleSidebar)
overlay.addEventListener('click', () => {
    sidebar.classList.remove('open')
    overlay.classList.remove('show')
})
window.addEventListener('resize', applyState)
applyState()
setActivePage()