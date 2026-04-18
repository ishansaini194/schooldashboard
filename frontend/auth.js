function getToken() {
    return localStorage.getItem('token')
}

function logout() {
    localStorage.removeItem('token')
    localStorage.removeItem('username')
    localStorage.removeItem('role')
    localStorage.removeItem('student_id')
    localStorage.removeItem('epunjab_id')
    localStorage.removeItem('student_class')
    localStorage.removeItem('student_section')
    window.location.href = '/frontend/login/login.html'
}

function requireAuth() {
    if (!getToken()) {
        window.location.href = '/frontend/login/login.html'
    }
}

function requireRole(role) {
    const userRole = localStorage.getItem('role')
    if (!getToken()) {
        window.location.href = '/frontend/login/login.html'
        return
    }
    if (userRole !== role) {
        window.location.href = '/frontend/login/login.html'
    }
}

async function authFetch(url, options = {}) {
    const token = getToken()
    const headers = {
        'Content-Type': 'application/json',
        ...options.headers,
    }
    if (token) headers['Authorization'] = `Bearer ${token}`

    const res = await fetch(url, { ...options, headers })

    if (res.status === 401) {
        logout()
        return
    }

    return res
}

// ── HTML escape helper — use for ALL user-controlled content in innerHTML ──
function escapeHtml(s) {
    if (s === null || s === undefined) return ''
    return String(s)
        .replace(/&/g, '&amp;')
        .replace(/</g, '&lt;')
        .replace(/>/g, '&gt;')
        .replace(/"/g, '&quot;')
        .replace(/'/g, '&#039;')
}
// short alias
const esc = escapeHtml