requireAuth()
const API = CONFIG.API
let foundUserId = null

async function searchUser() {
    const searchId = document.getElementById('searchId').value.trim()
    if (!searchId) { showToast('Please enter an ID', 'error'); return }

    // try finding by epunjab_id via student endpoint
    const res = await authFetch(`${API}/students/epunjab/${searchId}`)

    if (res && res.ok) {
        const student = await res.json()

        // now find their user account
        const userRes = await authFetch(`${API}/users/epunjab/${searchId}`)
        if (userRes && userRes.ok) {
            const user = await userRes.json()
            foundUserId = user.id || user.ID
            document.getElementById('userInfo').innerHTML = `
                <div style="background:var(--accent-light); border:1px solid #c8dece; border-radius:8px; padding:12px 16px;">
                    <div style="font-size:14px; font-weight:500; color:var(--text);">${student.name}</div>
                    <div style="font-size:12px; color:var(--muted);">Class ${student.class} · Roll ${student.roll_no} · ${searchId}</div>
                </div>
            `
            document.getElementById('userFound').style.display = 'block'
        } else {
            showToast('No user account found for this ID', 'error')
        }
    } else {
        showToast('Student not found', 'error')
    }
}

async function resetPassword() {
    if (!foundUserId) { showToast('Please search for a user first', 'error'); return }
    const newPassword = document.getElementById('newPassword').value.trim()
    if (!newPassword) { showToast('Please enter a new password', 'error'); return }

    const res = await authFetch(`${API}/auth/reset-password/${foundUserId}`, {
        method: 'PUT',
        body: JSON.stringify({ new_password: newPassword })
    })

    if (res && res.ok) {
        showToast('Password reset successfully!')
        document.getElementById('newPassword').value = ''
        document.getElementById('userFound').style.display = 'none'
        document.getElementById('searchId').value = ''
        foundUserId = null
    } else {
        showToast('Failed to reset password', 'error')
    }
}

function showToast(msg, type = 'success') {
    const toast = document.getElementById('toast')
    toast.textContent = msg
    toast.className = `toast ${type} show`
    setTimeout(() => toast.classList.remove('show'), 3000)
}