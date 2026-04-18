requireAuth()
const API = CONFIG.API

async function init() {
    await loadClasses()
    await loadPapers()
}

async function loadClasses() {
    const res = await authFetch(`${API}/classes`)
    if (!res || !res.ok) return
    const classes = await res.json() || []
    const unique = [...new Map(classes.map(c => [c.class, c])).values()]
    document.getElementById('filterClass').innerHTML =
        '<option value="">All Classes</option>' +
        unique.map(c => `<option value="${c.class}">Class ${c.class}</option>`).join('')
}

async function loadPapers() {
    const subject = document.getElementById('filterSubject').value
    const examType = document.getElementById('filterExam').value
    const cls = document.getElementById('filterClass').value

    let url = `${API}/papers?`
    if (cls) url += `class=${cls}&`
    if (subject) url += `subject=${subject}&`
    if (examType) url += `exam_type=${examType}&`

    document.getElementById('paperList').innerHTML = '<div class="dash-loading">Loading...</div>'

    const res = await authFetch(url)
    if (!res || !res.ok) return
    const list = await res.json() || []

    if (list.length === 0) {
        document.getElementById('paperList').innerHTML = '<div class="dash-empty">No papers uploaded yet</div>'
        return
    }

    document.getElementById('paperList').innerHTML = list.map(p => `
        <div class="paper-card">
            <div class="paper-info">
                <div class="paper-title">${esc(p.subject)} — ${p.exam_type === 'midterm' ? 'Mid-term' : 'Final'} ${p.year}</div>
                <div class="paper-meta">Class ${p.class}${p.section ? ' ' + p.section : ''} · Uploaded by ${esc(p.uploaded_by || 'Teacher')}</div>
            </div>
            <a href="${p.drive_link}" target="_blank" class="download-btn">Open ↗</a>
        </div>
    `).join('')
}

init()