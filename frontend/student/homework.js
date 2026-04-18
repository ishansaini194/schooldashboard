requireAuth()

const API = CONFIG.API
const cls = localStorage.getItem('student_class')
const section = localStorage.getItem('student_section')

async function init() {
    document.getElementById('classLabel').textContent =
        `Class ${cls}${section ? ' — ' + section : ''}`

    if (!cls) {
        document.getElementById('hwList').innerHTML = '<div class="dash-empty">Class info not found. Please login again.</div>'
        return
    }

    const res = await authFetch(`${API}/homework/class/${cls}/section/${section || ''}`)
    if (!res || !res.ok) return
    const list = await res.json() || []

    if (list.length === 0) {
        document.getElementById('hwList').innerHTML = '<div class="dash-empty">No homework assigned yet</div>'
        return
    }

    // group by date
    const grouped = {}
    list.forEach(h => {
        const date = formatDate(h.created_at)
        if (!grouped[date]) grouped[date] = []
        grouped[date].push(h)
    })

    document.getElementById('hwList').innerHTML = Object.entries(grouped).map(([date, items]) => `
        <div class="date-group-label">${date}</div>
        ${items.map(h => `
            <div class="info-card">
                <div class="info-card-top">
                    <span class="subject-badge">${esc(h.subject)}</span>
                </div>
                <div class="info-content">${esc(h.content)}</div>
                <div class="info-by">— ${esc(h.created_by || 'Teacher')}</div>
            </div>
        `).join('')}
    `).join('')
}

function formatDate(dateStr) {
    if (!dateStr) return '—'
    const d = new Date(dateStr)
    const today = new Date()
    const yesterday = new Date(today)
    yesterday.setDate(today.getDate() - 1)

    if (d.toDateString() === today.toDateString()) return 'Today'
    if (d.toDateString() === yesterday.toDateString()) return 'Yesterday'
    return d.toLocaleDateString('en-IN', { weekday: 'long', day: 'numeric', month: 'short' })
}

init()