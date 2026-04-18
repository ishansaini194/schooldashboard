requireAuth()
const API = CONFIG.API
const username = localStorage.getItem('username')
let allClasses = []

async function init() {
    await loadClasses()
    await loadHomework()
}

async function loadClasses() {
    const res = await authFetch(`${API}/classes`)
    if (!res || !res.ok) return
    const classes = await res.json() || []

    // deduplicate by class number
    const unique = [...new Map(classes.map(c => [c.class, c])).values()]

    const options = '<option value="">Select class</option>' +
        unique.map(c => `<option value="${c.class}">${c.class}</option>`).join('')

    // set whichever selects exist on the page
    const selects = ['hwClass', 'filterClass', 'resClass', 'paperClass', 'feeClass']
    selects.forEach(id => {
        const el = document.getElementById(id)
        if (el) el.innerHTML = options
    })
}

async function loadHomework() {
    const filterClass = document.getElementById('filterClass').value
    let url = filterClass
        ? `${API}/homework/class/${filterClass}/section/`
        : `${API}/homework/class/all/section/all`

    const res = await authFetch(url)
    if (!res || !res.ok) {
        document.getElementById('hwList').innerHTML = '<div class="dash-empty">No homework found</div>'
        return
    }
    const list = await res.json() || []

    if (list.length === 0) {
        document.getElementById('hwList').innerHTML = '<div class="dash-empty">No homework posted yet</div>'
        return
    }

    document.getElementById('hwList').innerHTML = list.map(h => `
        <div class="info-card">
            <div class="info-card-top">
                <div style="display:flex; gap:8px; align-items:center;">
                    <span class="subject-badge">${esc(h.subject)}</span>
                    <span class="class-badge">Class ${h.class}${h.section ? '-' + h.section : ''}</span>
                </div>
                <div style="display:flex; gap:8px; align-items:center;">
                    <span class="info-date">${formatDate(h.created_at)}</span>
                    <button class="btn-delete" onclick="deleteHomework(${h.ID || h.id})">✕</button>
                </div>
            </div>
            <div class="info-content">${esc(h.content)}</div>
            <div class="info-by">— ${esc(h.created_by || 'Teacher')}</div>
        </div>
    `).join('')
}

async function postHomework() {
    const cls = document.getElementById('hwClass').value
    const section = document.getElementById('hwSection').value
    const subject = document.getElementById('hwSubject').value
    const content = document.getElementById('hwContent').value.trim()

    if (!cls || !subject || !content) {
        showToast('Please fill all fields', 'error')
        return
    }

    const res = await authFetch(`${API}/homework`, {
        method: 'POST',
        body: JSON.stringify({
            class: cls,
            section: section,
            subject: subject,
            content: content,
            created_by: username || 'Teacher'
        })
    })

    if (res && res.ok) {
        showToast('Homework posted!')
        document.getElementById('hwContent').value = ''
        await loadHomework()
    } else {
        showToast('Failed to post homework', 'error')
    }
}

async function deleteHomework(id) {
    const ok = await confirmDialog({ title: 'Delete this homework?', confirmText: 'Delete', danger: true })
    if (!ok) return
    const res = await authFetch(`${API}/homework/${id}`, { method: 'DELETE' })
    if (res && res.ok) {
        showToast('Deleted')
        await loadHomework()
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