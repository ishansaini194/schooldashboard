const API = CONFIG.API

const roles = document.querySelectorAll('.role')
let selectedRole = 'admin'

roles.forEach(btn => {
    btn.addEventListener('click', () => {
        roles.forEach(b => b.classList.remove('active'))
        btn.classList.add('active')
        selectedRole = btn.dataset.role
        // hint to user about which login id format to use
        updatePlaceholder()
    })
})

function updatePlaceholder() {
    const input = document.querySelector('input[name="username"]')
    if (!input) return
    if (selectedRole === 'admin') input.placeholder = 'e.g. ADM001'
    else if (selectedRole === 'teacher') input.placeholder = 'e.g. TCH001'
    else input.placeholder = 'ePunjab ID (e.g. EP2024001)'
}

function showLoginError(msg) {
    let el = document.getElementById('loginError')
    if (!el) {
        el = document.createElement('div')
        el.id = 'loginError'
        el.className = 'login-error'
        document.getElementById('loginForm').insertBefore(el, document.querySelector('.remember'))
    }
    el.textContent = msg
    el.classList.add('show')
}

function clearLoginError() {
    const el = document.getElementById('loginError')
    if (el) el.classList.remove('show')
}

document.getElementById('loginForm').addEventListener('submit', async function (e) {
    e.preventDefault()
    clearLoginError()

    const submitBtn = this.querySelector('button[type="submit"]')
    submitBtn.disabled = true
    submitBtn.textContent = 'Logging in...'

    try {
        const data = {
            username: document.querySelector('input[name="username"]').value.trim(),
            password: document.querySelector('input[name="password"]').value,
        }

        const res = await fetch(`${API}/auth/login`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(data)
        })

        if (!res.ok) {
            const err = await res.json().catch(() => ({ error: 'Login failed' }))
            showLoginError(err.error || 'Invalid credentials')
            return
        }

        const result = await res.json()

        // verify the user's actual role matches the selected tab
        if (result.role !== selectedRole) {
            showLoginError(`This account is a ${result.role}. Please select the correct role tab above.`)
            return
        }

        localStorage.clear()
        localStorage.setItem('token', result.token)
        localStorage.setItem('username', result.username)
        localStorage.setItem('role', result.role)

        if (result.role === 'admin') {
            window.location.href = '../admin/dashboard.html'
        } else if (result.role === 'teacher') {
            window.location.href = '../teacher/dashboard.html'
        } else if (result.role === 'student') {
            localStorage.setItem('student_id', result.student_id)
            localStorage.setItem('epunjab_id', result.epunjab_id)

            try {
                const profileRes = await fetch(`${API}/students/epunjab/${result.epunjab_id}`, {
                    headers: { 'Authorization': `Bearer ${result.token}` }
                })
                if (profileRes.ok) {
                    const profile = await profileRes.json()
                    localStorage.setItem('student_class', profile.class)
                    localStorage.setItem('student_section', profile.section || '')
                }
            } catch (_) { /* non-fatal */ }

            window.location.href = '../student/dashboard.html'
        }
    } catch (err) {
        showLoginError('Network error. Please check your connection.')
    } finally {
        submitBtn.disabled = false
        submitBtn.textContent = 'Login'
    }
})

// auto-focus username on page load + set initial placeholder
document.querySelector('input[name="username"]')?.focus()
updatePlaceholder()
