requireAuth()
const API = CONFIG.API

const username = localStorage.getItem('username')

async function init() {
    setTodayDate()
    setTeacherInfo()
    await Promise.all([
        loadHomework(),
        loadNotices(),
        loadResults(),
        loadPapers(),
    ])
}

function setTodayDate() {
    document.getElementById('todayDate').textContent =
        new Date().toLocaleDateString('en-IN', { weekday: 'long', day: 'numeric', month: 'long', year: 'numeric' })
}

function setTeacherInfo() {
    const name = username || 'Teacher'
    document.getElementById('teacherName').textContent = name
    document.getElementById('teacherAvatar').textContent = name[0].toUpperCase()
}

async function loadHomework() {
    // fetch from all classes
    const classRes = await authFetch(`${API}/classes`)
    if (!classRes || !classRes.ok) { document.getElementById('statHomework').textContent = '0'; return }
    const classes = await classRes.json() || []

    // fetch homework for each class and combine
    const results = await Promise.all(
        classes.map(c =>
            authFetch(`${API}/homework/class/${c.class}/section/${c.section || ''}`)
                .then(r => r?.json())
                .then(data => data || [])
        )
    )

    const list = results.flat()
    const today = new Date().toDateString()
    const todayCount = list.filter(h => new Date(h.created_at).toDateString() === today).length
    document.getElementById('statHomework').textContent = todayCount

    if (list.length === 0) {
        document.getElementById('recentHomework').innerHTML = '<div class="dash-empty">No homework posted yet</div>'
        return
    }

    document.getElementById('recentHomework').innerHTML = list.slice(0, 3).map(h => `
        <div class="info-card">
            <div class="info-card-top">
                <span class="subject-badge">${esc(h.subject)}</span>
                <span class="class-badge">Class ${h.class}${h.section ? '-' + h.section : ''}</span>
                <span class="info-date">${formatDate(h.created_at)}</span>
            </div>
            <div class="info-content">${esc(h.content)}</div>
        </div>
    `).join('')
}

async function loadNotices() {
    const res = await authFetch(`${API}/notices`)
    if (!res || !res.ok) { document.getElementById('statNotices').textContent = '0'; return }
    const list = await res.json() || []
    document.getElementById('statNotices').textContent = list.length

    document.getElementById('recentNotices').innerHTML = list.slice(0, 3).map(n => `
        <div class="info-card notice">
            <div class="info-card-top">
                <span class="notice-title">${esc(n.title)}</span>
                <span class="info-date">${formatDate(n.created_at)}</span>
            </div>
            <div class="info-content" style="margin-top:6px;">${esc(n.body)}</div>
        </div>
    `).join('') || '<div class="dash-empty">No notices posted yet</div>'
}

async function loadResults() {
    try {
        const res = await authFetch(`${API}/results/mine`)
        if (!res || !res.ok) { document.getElementById('statResults').textContent = '0'; return }
        const list = await res.json() || []
        document.getElementById('statResults').textContent = list.length
    } catch (_) {
        document.getElementById('statResults').textContent = '0'
    }
}

async function loadPapers() {
    try {
        const res = await authFetch(`${API}/papers`)
        if (!res || !res.ok) { document.getElementById('statPapers').textContent = '0'; return }
        const list = await res.json() || []
        document.getElementById('statPapers').textContent = list.length
    } catch (_) {
        document.getElementById('statPapers').textContent = '0'
    }
}

function formatDate(dateStr) {
    if (!dateStr) return '—'
    const d = new Date(dateStr)
    const today = new Date()
    const yesterday = new Date(today)
    yesterday.setDate(today.getDate() - 1)
    if (d.toDateString() === today.toDateString()) return 'Today'
    if (d.toDateString() === yesterday.toDateString()) return 'Yesterday'
    return d.toLocaleDateString('en-IN', { day: 'numeric', month: 'short' })
}

function showToast(msg, type = 'success') {
    const toast = document.getElementById('toast')
    toast.textContent = msg
    toast.className = `toast ${type} show`
    setTimeout(() => toast.classList.remove('show'), 3000)
}

init()