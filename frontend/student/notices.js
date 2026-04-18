requireAuth()

const API = CONFIG.API
const cls = localStorage.getItem('student_class')
const section = localStorage.getItem('student_section')

async function init() {
    const target = cls && section ? `${cls}-${section}` : ''
    const res = await authFetch(`${API}/notices?target=${target}`)
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
                <span class="info-date">${formatDate(n.created_at)}</span>
            </div>
            <div class="info-content" style="margin-top:8px;">${esc(n.body)}</div>
            <div class="info-by" style="margin-top:8px;">Posted by ${esc(n.posted_by || 'Admin')}</div>
        </div>
    `).join('')
}

function formatDate(dateStr) {
    if (!dateStr) return '—'
    return new Date(dateStr).toLocaleDateString('en-IN', { day: 'numeric', month: 'short', year: 'numeric' })
}

init()