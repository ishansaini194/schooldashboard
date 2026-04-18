requireAuth()
const API = CONFIG.API
const username = localStorage.getItem('username')

async function init() {
    await loadClasses()
    await loadPapers()
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

async function loadPapers() {
    const res = await authFetch(`${API}/papers`)
    if (!res || !res.ok) return
    const list = await res.json() || []

    if (list.length === 0) {
        document.getElementById('paperList').innerHTML = '<div class="dash-empty">No papers uploaded yet</div>'
        return
    }

    document.getElementById('paperList').innerHTML = list.map(p => `
        <div class="paper-row">
            <div class="paper-info">
                <div class="paper-title">Class ${p.class}${p.section ? '-' + p.section : ''} — ${esc(p.subject)}</div>
                <div class="paper-meta">${p.exam_type === 'midterm' ? 'Mid-term' : 'Final'} · ${p.year} · Uploaded by ${esc(p.uploaded_by || 'Teacher')}</div>
            </div>
            <div style="display:flex; gap:8px;">
                <a href="${p.drive_link}" target="_blank" class="download-btn">Open ↗</a>
                <button class="btn-delete" onclick="deletePaper(${p.ID || p.id})">✕</button>
            </div>
        </div>
    `).join('')
}

async function addPaper() {
    const cls = document.getElementById('paperClass').value
    const section = document.getElementById('paperSection').value
    const subject = document.getElementById('paperSubject').value
    const examType = document.getElementById('paperExamType').value
    const year = parseInt(document.getElementById('paperYear').value)
    const driveLink = document.getElementById('paperLink').value.trim()

    if (!cls || !subject || !driveLink) {
        showToast('Please fill all fields', 'error')
        return
    }

    const res = await authFetch(`${API}/papers`, {
        method: 'POST',
        body: JSON.stringify({
            class: cls,
            section,
            subject,
            exam_type: examType,
            year,
            drive_link: driveLink,
            uploaded_by: username || 'Teacher'
        })
    })

    if (res && res.ok) {
        showToast('Paper added!')
        document.getElementById('paperLink').value = ''
        await loadPapers()
    } else {
        showToast('Failed to add paper', 'error')
    }
}

async function deletePaper(id) {
    const ok = await confirmDialog({ title: 'Delete this paper?', confirmText: 'Delete', danger: true })
    if (!ok) return
    const res = await authFetch(`${API}/papers/${id}`, { method: 'DELETE' })
    if (res && res.ok) {
        showToast('Deleted')
        await loadPapers()
    }
}

function showToast(msg, type = 'success') {
    const toast = document.getElementById('toast')
    toast.textContent = msg
    toast.className = `toast ${type} show`
    setTimeout(() => toast.classList.remove('show'), 3000)
}

init()