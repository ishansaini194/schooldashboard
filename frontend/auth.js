// frontend/auth.js
function getToken() {
    return localStorage.getItem('token')
}

function logout() {
    localStorage.removeItem('token')
    localStorage.removeItem('username')
    window.location.href = '../login/login.html'
}

// redirect to login if no token
function requireAuth() {
    if (!getToken()) {
        window.location.href = '../login/login.html'
    }
}

function requireRole(role) {
    const userRole = localStorage.getItem('role')
    if (!getToken()) {
        window.location.href = '../login/login.html'
        return
    }
    if (userRole !== role) {
        window.location.href = '../login/login.html'
    }
}

// add token to all fetch calls
async function authFetch(url, options = {}) {
    const token = getToken()
    const headers = {
        'Content-Type': 'application/json',
        ...options.headers,
    }
    if (token) headers['Authorization'] = `Bearer ${token}`

    const res = await fetch(url, { ...options, headers })

    // if 401 — token expired, redirect to login
    if (res.status === 401) {
        logout()
        return
    }

    return res
}