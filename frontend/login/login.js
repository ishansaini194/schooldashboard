const roles = document.querySelectorAll(".role")
let selectedRole = "admin"

// Role selection
roles.forEach(btn => {
    btn.addEventListener("click", () => {
        roles.forEach(b => b.classList.remove("active"))
        btn.classList.add("active")
        selectedRole = btn.dataset.role
    })
})

// Form submit
document.getElementById('loginForm').addEventListener('submit', async function (e) {
    e.preventDefault()

    const data = {
        username: document.querySelector('input[name="username"]').value,
        password: document.querySelector('input[name="password"]').value,
        role: selectedRole
    }

    const res = await fetch('http://localhost:8080/api/auth/login', {
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
            window.location.href = '../studentform/dashboard.html'
        } else if (result.role === 'student') {
            window.location.href = '../studentform/dashboard.html'
        }
    } else {
        const err = await res.json()
        alert(err.error)
    }
})