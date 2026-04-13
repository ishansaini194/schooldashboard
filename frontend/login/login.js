const API = CONFIG.API

const roles = document.querySelectorAll(".role")
let selectedRole = "admin"

roles.forEach(btn => {
    btn.addEventListener("click", () => {
        roles.forEach(b => b.classList.remove("active"))
        btn.classList.add("active")
        selectedRole = btn.dataset.role
    })
})

document.getElementById('loginForm').addEventListener('submit', async function (e) {
    e.preventDefault()

    const data = {
        username: document.querySelector('input[name="username"]').value,
        password: document.querySelector('input[name="password"]').value,
    }

    const res = await fetch(`${API}/auth/login`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(data)
    })

    if (res.ok) {
        const result = await res.json()
        localStorage.setItem('token', result.token)
        localStorage.setItem('username', result.username)
        localStorage.setItem('role', result.role)

        if (result.role === 'admin') {
            window.location.href = '../studentform/dashboard.html'
        } else if (result.role === 'teacher') {
            window.location.href = '../teacher/dashboard.html'
        } else if (result.role === 'student') {
            localStorage.setItem('student_id', result.student_id)
            localStorage.setItem('epunjab_id', result.epunjab_id)

            const profileRes = await fetch(`${API}/students/epunjab/${result.epunjab_id}`, {
                headers: { 'Authorization': `Bearer ${result.token}` }
            })
            if (profileRes.ok) {
                const profile = await profileRes.json()
                localStorage.setItem('student_class', profile.class)
                localStorage.setItem('student_section', profile.section || '')
            }

            window.location.href = '../student/dashboard.html'
        }
    } else {
        const err = await res.json()
        alert(err.error)
    }
})