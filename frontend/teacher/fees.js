requireAuth()
const API = CONFIG.API

async function init() {
    await loadClasses()
    setCurrentMonth()
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

function setCurrentMonth() {
    const month = new Date().toLocaleString('en-IN', { month: 'long' })
    const select = document.getElementById('feeMonth')
    for (const opt of select.options) {
        if (opt.value === month) { opt.selected = true; break }
    }
}

async function loadFees() {
    const cls = document.getElementById('feeClass').value
    if (!cls) return

    const month = document.getElementById('feeMonth').value
    const year = document.getElementById('feeYear').value

    const res = await authFetch(`${API}/fees/class/${cls}/month/${month}/year/${year}`)
    if (!res || !res.ok) return
    const students = await res.json() || []

    if (students.length === 0) {
        document.getElementById('feeList').innerHTML = '<div class="dash-empty">No students found</div>'
        return
    }

    const paid = students.filter(s => s.has_paid).length
    const total = students.length
    const pct = total > 0 ? Math.round((paid / total) * 100) : 0

    document.getElementById('feeList').innerHTML = `
        <div class="fee-progress-card">
            <div class="fee-progress-header">
                <span>${paid}/${total} paid</span>
                <span>${pct}%</span>
            </div>
            <div class="cp-bar-wrap">
                <div class="cp-bar" style="width:${pct}%; background:${pct >= 80 ? 'var(--accent)' : pct >= 50 ? '#c07020' : '#c0392b'}"></div>
            </div>
        </div>
        ${students.map(s => {
        const fee = s.fees?.[0]
        const status = fee?.status || 'unpaid'
        return `
            <div class="pending-row ${status}">
                <div class="pending-class">Roll ${s.roll_no || '—'}</div>
                <div class="pending-info">
                    <div class="pending-name">${esc(s.student_name)}</div>
                </div>
                <div class="pending-right">
                    ${fee?.paid_amount > 0 ? `<div class="pending-amount">₹${fee.paid_amount.toLocaleString()}</div>` : ''}
                    <span class="pending-badge ${status}">${status}</span>
                </div>
            </div>`
    }).join('')}
    `
}

function showToast(msg, type = 'success') {
    const toast = document.getElementById('toast')
    toast.textContent = msg
    toast.className = `toast ${type} show`
    setTimeout(() => toast.classList.remove('show'), 3000)
}

init()