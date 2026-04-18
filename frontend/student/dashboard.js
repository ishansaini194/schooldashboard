requireAuth()

const API = CONFIG.API
const studentID = localStorage.getItem('student_id')
const epunjabID = localStorage.getItem('epunjab_id')

async function init() {
    setTodayDate()
    await Promise.all([
        loadStudentInfo(),
        loadHomework(),
        loadNotices(),
        loadFeeStatus(),
    ])
}

function setTodayDate() {
    document.getElementById('todayDate').textContent =
        new Date().toLocaleDateString('en-IN', { weekday: 'long', day: 'numeric', month: 'long', year: 'numeric' })
}

async function loadStudentInfo() {
    const res = await authFetch(`${API}/students/epunjab/${epunjabID}`)
    if (!res || !res.ok) return
    const s = await res.json()

    document.getElementById('studentName').textContent = s.name || '—'
    document.getElementById('studentMeta').textContent = `Class ${s.class}${s.section ? ' — ' + s.section : ''} · Roll ${s.roll_no || '—'}`
    document.getElementById('studentAvatar').textContent = s.name ? s.name[0].toUpperCase() : '?'
}

async function loadFeeStatus() {
    const month = new Date().toLocaleString('en-IN', { month: 'long' })
    const year = new Date().getFullYear()
    const res = await authFetch(`${API}/fees/student/${studentID}/yearly?year=${year}`)
    if (!res || !res.ok) return
    const data = await res.json()

    const currentMonthFee = data.months?.find(m => m.month === month)
    const pill = document.getElementById('feePill')
    if (!currentMonthFee || currentMonthFee.status === 'unpaid') {
        pill.textContent = 'Fee Due'
        pill.className = 'fee-pill unpaid'
    } else if (currentMonthFee.status === 'partial') {
        const remaining = currentMonthFee.remaining || 0
        pill.textContent = remaining > 0 ? `₹${remaining.toLocaleString()} due` : 'Partial'
        pill.className = 'fee-pill partial'
    } else {
        const paid = currentMonthFee.paid_amount || 0
        pill.textContent = paid > 0 ? `Paid ₹${paid.toLocaleString()} ✓` : 'Fee Paid ✓'
        pill.className = 'fee-pill paid'
    }
}

async function loadHomework() {
    // need class/section from student profile — read from localStorage or fetch
    const cls = localStorage.getItem('student_class')
    const section = localStorage.getItem('student_section')
    if (!cls) { document.getElementById('hwList').innerHTML = '<div class="dash-empty">No homework found</div>'; return }

    const res = await authFetch(`${API}/homework/class/${cls}/section/${section || ''}`)
    if (!res || !res.ok) return
    const list = await res.json() || []

    document.getElementById('hwCount').textContent = `${list.length} today`

    if (list.length === 0) {
        document.getElementById('hwList').innerHTML = '<div class="dash-empty">No homework assigned</div>'
        return
    }

    document.getElementById('hwList').innerHTML = list.slice(0, 3).map(h => `
        <div class="info-card">
            <div class="info-card-top">
                <span class="subject-badge">${esc(h.subject)}</span>
                <span class="info-date">${formatDate(h.created_at)}</span>
            </div>
            <div class="info-content">${esc(h.content)}</div>
            <div class="info-by">— ${esc(h.created_by || 'Teacher')}</div>
        </div>
    `).join('')
}

async function loadNotices() {
    const cls = localStorage.getItem('student_class')
    const section = localStorage.getItem('student_section')
    const target = cls && section ? `${cls}-${section}` : ''

    const res = await authFetch(`${API}/notices?target=${target}`)
    if (!res || !res.ok) return
    const list = await res.json() || []

    document.getElementById('noticeCount').textContent = `${list.length} new`

    if (list.length === 0) {
        document.getElementById('noticeList').innerHTML = '<div class="dash-empty">No notices</div>'
        return
    }

    document.getElementById('noticeList').innerHTML = list.slice(0, 2).map(n => `
        <div class="info-card notice">
            <div class="info-card-top">
                <span class="notice-title">${esc(n.title)}</span>
                <span class="info-date">${formatDate(n.created_at)}</span>
            </div>
            <div class="info-content">${esc(n.body)}</div>
        </div>
    `).join('')
}

function formatDate(dateStr) {
    if (!dateStr) return '—'
    return new Date(dateStr).toLocaleDateString('en-IN', { day: 'numeric', month: 'short' })
}

function showToast(msg, type = 'success') {
    const toast = document.getElementById('toast')
    toast.textContent = msg
    toast.className = `toast ${type} show`
    setTimeout(() => toast.classList.remove('show'), 3000)
}

init()