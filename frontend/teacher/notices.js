requireAuth()
const API = CONFIG.API
const username = localStorage.getItem('username')

async function init() {
    await loadClasses()
    await loadNotices()
}

async function loadClasses() {
    const res = await authFetch(`${API}/classes`)
    if (!res || !res.ok) return
    const classes = await res.json() || []

    const options = classes.map(c =>
        `<option value="${c.class}-${c.section}">Class ${c.class}${c.section ? ' — ' + c.section : ''}</option>`
    ).join('')
    document.getElementById('noticeTarget').innerHTML =
        '<option value="all">All Students</option>' + options
}

async function loadNotices() {
    const res = await authFetch(`${API}/notices`)
    if (!res || !res.ok) return
    const list = await res.json() || []

    if (list.length === 0) {
        document.getElementById('noticeList').innerHTML = '<div class="dash-empty">No notices posted yet</div>'
        return
    }

    document.getElementById('noticeList').innerHTML = list.map(n => `
        <div class="info-card notice">
            <div class="info-card-top">
                <span class="notice-title">${esc(n.title)}</span>
                <div style="display:flex; gap:8px; align-items:center;">
                    <span class="info-date">${formatDate(n.created_at)}</span>
                    <button class="btn-delete" onclick="deleteNotice(${n.ID || n.id})">✕</button>
                </div>
            </div>
            <div class="info-content" style="margin-top:6px;">${esc(n.body)}</div>
            <div class="info-by">Target: ${n.target === 'all' ? 'All students' : 'Class ' + n.target} · ${esc(n.posted_by || 'Teacher')}</div>
        </div>
    `).join('')
}

async function postNotice() {
    const title = document.getElementById('noticeTitle').value.trim()
    const body = document.getElementById('noticeBody').value.trim()
    const target = document.getElementById('noticeTarget').value

    if (!title || !body) {
        showToast('Please fill all fields', 'error')
        return
    }

    const res = await authFetch(`${API}/notices`, {
        method: 'POST',
        body: JSON.stringify({
            title,
            body,
            target,
            posted_by: username || 'Teacher'
        })
    })

    if (res && res.ok) {
        showToast('Notice posted!')
        document.getElementById('noticeTitle').value = ''
        document.getElementById('noticeBody').value = ''
        await loadNotices()
    } else {
        showToast('Failed to post notice', 'error')
    }
}

async function deleteNotice(id) {
    const ok = await confirmDialog({ title: 'Delete this notice?', confirmText: 'Delete', danger: true })
    if (!ok) return
    const res = await authFetch(`${API}/notices/${id}`, { method: 'DELETE' })
    if (res && res.ok) {
        showToast('Deleted')
        await loadNotices()
    }
}

function formatDate(dateStr) {
    if (!dateStr) return '—'
    return new Date(dateStr).toLocaleDateString('en-IN', { day: 'numeric', month: 'short', year: 'numeric' })
}

function showToast(msg, type = 'success') {
    const toast = document.getElementById('toast')
    toast.textContent = msg
    toast.className = `toast ${type} show`
    setTimeout(() => toast.classList.remove('show'), 3000)
}

init()