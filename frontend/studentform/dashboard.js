// ─────────────────────────────────────────────────────────────
// dashboard.js — React-ready structure
// ─────────────────────────────────────────────────────────────

const API = 'http://localhost:8080/api'

// ── State ─────────────────────────────────────────────────────
const MONTHS = ['January', 'February', 'March', 'April', 'May', 'June',
    'July', 'August', 'September', 'October', 'November', 'December']

let currentMonth = new Date().getMonth()
let currentYear = new Date().getFullYear()
let allClasses = []
let allStudents = []
let pendingData = []

// ── Init ──────────────────────────────────────────────────────

async function init() {
    setTodayDate()
    updateMonthLabel()
    await Promise.all([
        loadClasses(),
        loadRecent(),
    ])
    await loadAllStudents()
    await Promise.all([
        loadSummary(),
        loadClassProgress(),
        loadPending(),
    ])
}

// ── Date Helpers ──────────────────────────────────────────────

function setTodayDate() {
    const now = new Date()
    document.getElementById('todayDate').textContent =
        now.toLocaleDateString('en-IN', { weekday: 'long', day: 'numeric', month: 'long', year: 'numeric' })
}

function updateMonthLabel() {
    document.getElementById('monthLabel').textContent = `${MONTHS[currentMonth]} ${currentYear}`
    document.getElementById('progressMonthLabel').textContent = `${MONTHS[currentMonth]} ${currentYear}`
}

function changeMonth(dir) {
    currentMonth += dir
    if (currentMonth > 11) { currentMonth = 0; currentYear++ }
    if (currentMonth < 0) { currentMonth = 11; currentYear-- }
    updateMonthLabel()
    loadSummary()
    loadClassProgress()
    loadPending()
}

// ── Data Fetching ─────────────────────────────────────────────

async function loadClasses() {
    const res = await fetch(`${API}/classes`)
    allClasses = await res.json() || []
}

async function loadAllStudents() {
    const results = await Promise.all(
        allClasses.map(c =>
            fetch(`${API}/students/class/${c.class}`)
                .then(r => r.json())
                .then(data => (data || []).map(s => ({ ...s, section: c.section })))
        )
    )
    allStudents = results.flat()
}

async function loadSummary() {
    const month = MONTHS[currentMonth]
    const res = await fetch(`${API}/dashboard/summary?month=${month}&year=${currentYear}`)
    const data = await res.json()
    renderSummary(data)
}

async function loadClassProgress() {
    const month = MONTHS[currentMonth]
    const results = await Promise.all(
        allClasses.map(c =>
            fetch(`${API}/fees/class/${c.class}/month/${month}/year/${currentYear}`)
                .then(r => r.json())
                .then(students => ({ cls: c, students: students || [] }))
        )
    )
    renderClassProgress(results)
}

async function loadPending() {
    const month = MONTHS[currentMonth]
    const res = await fetch(`${API}/fees/pending/all?month=${month}&year=${currentYear}`)
    const data = await res.json()
    pendingData = data.pending || []
    renderPending(pendingData)
}

async function loadRecent() {
    const res = await fetch(`${API}/fees/recent`)
    const fees = await res.json() || []
    renderRecent(fees)
}

// ── Render Functions ──────────────────────────────────────────

// → SummaryCards + Progress Bar + Alert
function renderSummary(data) {
    document.getElementById('statStudents').textContent = data.total_students || 0
    document.getElementById('statClasses').textContent = data.total_classes || 0
    document.getElementById('statCollected').textContent = `₹${(data.total_collected || 0).toLocaleString()}`
    document.getElementById('statPending').textContent = data.pending_count || 0

    // overdue alert — red banner
    if (data.overdue_count > 0) {
        const banner = document.getElementById('alertBanner')
        banner.classList.remove('hidden')
        banner.classList.add('overdue')
        document.getElementById('alertText').textContent =
            `${data.overdue_count} partial payment${data.overdue_count > 1 ? 's are' : ' is'} overdue — immediate attention needed`
    } else {
        document.getElementById('alertBanner').classList.add('hidden')
    }

    // collection progress bar
    const collected = data.total_collected || 0
    const expected = data.expected_total || 0
    const pct = expected > 0 ? Math.min(Math.round((collected / expected) * 100), 100) : 0
    const remaining = expected - collected

    document.getElementById('cpCollected').textContent = `₹${collected.toLocaleString()}`
    document.getElementById('cpExpected').textContent = `₹${expected.toLocaleString()}`
    document.getElementById('cpPct').textContent = `${pct}%`
    document.getElementById('cpBar').style.width = `${pct}%`
    document.getElementById('cpRemaining').textContent =
        remaining > 0 ? `₹${remaining.toLocaleString()} remaining` : '✓ Fully collected'

    // color bar based on %
    const bar = document.getElementById('cpBar')
    if (pct >= 80) bar.style.background = 'var(--accent)'
    else if (pct >= 50) bar.style.background = '#c07020'
    else bar.style.background = '#c0392b'
}

// → ClassProgressCard
function renderClassProgress(results) {
    const list = document.getElementById('classProgressList')

    if (!results || results.length === 0) {
        list.innerHTML = '<div class="dash-empty">No classes found</div>'
        return
    }

    list.innerHTML = results.map(({ cls, students }) => {
        const total = students.length
        const paid = students.filter(s => s.has_paid).length
        const pct = total > 0 ? Math.round((paid / total) * 100) : 0
        const month = MONTHS[currentMonth]

        return `
        <div class="progress-card" onclick="window.location.href='fees.html?class=${cls.class}&month=${month}&year=${currentYear}'" style="cursor:pointer">
            <div class="progress-top">
                <div class="progress-class">
                    <div class="progress-class-name">Class ${cls.class}${cls.section ? ' — ' + cls.section : ''}</div>
                    <div class="progress-teacher">${cls.teacher_name || 'No teacher'}</div>
                </div>
                <div class="progress-count">${paid}/${total} paid</div>
            </div>
            <div class="progress-bar-wrap">
                <div class="progress-bar" style="width:${pct}%"></div>
            </div>
            <div class="progress-footer">
                <span class="progress-pct">${pct}%</span>
                ${total - paid > 0
                ? `<span class="progress-pending">${total - paid} pending</span>`
                : '<span class="progress-done">✓ All paid</span>'}
            </div>
        </div>`
    }).join('')
}

// → PendingRow
function renderPending(pending) {
    const list = document.getElementById('pendingList')

    if (!pending || pending.length === 0) {
        list.innerHTML = '<div class="dash-empty">🎉 No pending fees this month</div>'
        return
    }

    const shown = pending.slice(0, 8)

    list.innerHTML = shown.map(s => `
        <div class="pending-row ${s.status}">
            <div class="pending-class">Class ${s.class}</div>
            <div class="pending-info">
                <div class="pending-name">${s.student_name}</div>
                <div class="pending-roll">Roll ${s.roll_no || '—'}</div>
            </div>
            <div class="pending-right">
                ${s.remaining > 0 ? `<div class="pending-amount">₹${s.remaining.toLocaleString()} left</div>` : ''}
                <span class="pending-badge ${s.status}">${s.status}</span>
            </div>
            <button class="btn-mini"
                onclick="window.location.href='fee-collect.html?student_id=${s.student_id}&month=${MONTHS[currentMonth]}&year=${currentYear}'">
                Pay
            </button>
        </div>
    `).join('')

    if (pending.length > 8) {
        list.innerHTML += `<div class="pending-more">+${pending.length - 8} more — <a href="fees.html?tab=pending">View all</a></div>`
    }
}

// → RecentPaymentRow
function renderRecent(fees) {
    const list = document.getElementById('recentList')

    if (!fees || fees.length === 0) {
        list.innerHTML = '<div class="dash-empty">No recent payments</div>'
        return
    }

    list.innerHTML = fees.map(f => {
        const timeAgo = getTimeAgo(f.paid_at || f.CreatedAt)
        return `
        <div class="recent-row">
            <div class="recent-avatar">${f.student_name ? f.student_name[0].toUpperCase() : '?'}</div>
            <div class="recent-info">
                <div class="recent-name">${f.student_name || '—'}</div>
                <div class="recent-meta">Class ${f.class} · ${f.month} · ${f.fee_type === 'transport' ? '🚌' : '📚'}</div>
            </div>
            <div class="recent-right">
                <div class="recent-amount">₹${(f.paid_amount || 0).toLocaleString()}</div>
                <div class="recent-time">${timeAgo}</div>
            </div>
            <button class="btn-mini secondary"
                onclick="window.location.href='fee-receipt.html?receipt=${f.receipt_no}'">
                Receipt
            </button>
        </div>`
    }).join('')
}

// ── Quick Search ──────────────────────────────────────────────

let searchTimeout = null

function quickSearch() {
    clearTimeout(searchTimeout)
    searchTimeout = setTimeout(() => {
        const q = document.getElementById('quickSearch').value.toLowerCase().trim()
        const results = document.getElementById('quickSearchResults')

        if (!q) {
            results.classList.add('hidden')
            return
        }

        const matches = allStudents.filter(s =>
            s.name?.toLowerCase().includes(q) ||
            s.roll_no?.toLowerCase().includes(q) ||
            s.epunjab_id?.toLowerCase().includes(q)
        ).slice(0, 6)

        if (matches.length === 0) {
            results.innerHTML = '<div class="qs-empty">No students found</div>'
        } else {
            results.innerHTML = matches.map(s => `
                <div class="qs-item" onclick="window.location.href='student-detail.html?roll_no=${s.roll_no}'">
                    <div class="qs-avatar">${s.name ? s.name[0].toUpperCase() : '?'}</div>
                    <div class="qs-info">
                        <div class="qs-name">${s.name}</div>
                        <div class="qs-meta">Class ${s.class}${s.section ? ' — ' + s.section : ''} · Roll ${s.roll_no || '—'}</div>
                    </div>
                    <div class="qs-arrow">→</div>
                </div>
            `).join('')
        }

        results.classList.remove('hidden')
    }, 200)
}

// close search on outside click
document.addEventListener('click', e => {
    if (!e.target.closest('.quick-search-wrap')) {
        document.getElementById('quickSearchResults').classList.add('hidden')
    }
})

// ── Export ────────────────────────────────────────────────────

function exportPendingCSV() {
    if (!pendingData || pendingData.length === 0) {
        showToast('No pending fees to export', 'error')
        return
    }

    const headers = ['Class', 'Roll No', 'Name', 'Phone', 'Status', 'Paid', 'Remaining', 'Due Date']
    const rows = pendingData.map(s => [
        s.class, s.roll_no, s.student_name, s.phone,
        s.status, s.paid_amount, s.remaining, s.due_date || ''
    ].map(v => `"${(v || '').toString().replace(/"/g, '""')}"`))

    const csv = [headers.join(','), ...rows.map(r => r.join(','))].join('\n')
    const blob = new Blob([csv], { type: 'text/csv' })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `pending-fees-${MONTHS[currentMonth]}-${currentYear}.csv`
    a.click()
    URL.revokeObjectURL(url)
    showToast('CSV downloaded!')
}

function exportPendingPDF() {
    if (!pendingData || pendingData.length === 0) {
        showToast('No pending fees to export', 'error')
        return
    }

    const printArea = document.getElementById('pdfPrintArea')
    printArea.innerHTML = `
        <div class="pdf-header">
            <div class="pdf-school">KRB School</div>
            <div class="pdf-title">Pending Fee Report — ${MONTHS[currentMonth]} ${currentYear}</div>
            <div class="pdf-date">Generated: ${new Date().toLocaleDateString('en-IN')}</div>
        </div>
        <table class="pdf-table">
            <thead>
                <tr>
                    <th>Class</th><th>Roll No</th><th>Name</th><th>Phone</th>
                    <th>Status</th><th>Paid</th><th>Remaining</th><th>Due Date</th>
                </tr>
            </thead>
            <tbody>
                ${pendingData.map(s => `
                    <tr>
                        <td>${s.class}</td>
                        <td>${s.roll_no || '—'}</td>
                        <td>${s.student_name}</td>
                        <td>${s.phone || '—'}</td>
                        <td class="${s.status}">${s.status}</td>
                        <td>₹${(s.paid_amount || 0).toLocaleString()}</td>
                        <td>₹${(s.remaining || 0).toLocaleString()}</td>
                        <td>${s.due_date || '—'}</td>
                    </tr>
                `).join('')}
            </tbody>
        </table>
        <div class="pdf-footer">Total: ${pendingData.length} students pending</div>
    `

    printArea.classList.remove('no-print')
    document.body.classList.add('print-pdf-mode')
    window.print()
    setTimeout(() => {
        printArea.classList.add('no-print')
        document.body.classList.remove('print-pdf-mode')
    }, 1000)
}

// ── Helpers ───────────────────────────────────────────────────

function getTimeAgo(dateStr) {
    if (!dateStr) return '—'
    const diff = Date.now() - new Date(dateStr).getTime()
    const mins = Math.floor(diff / 60000)
    if (mins < 60) return `${mins}m ago`
    const hrs = Math.floor(mins / 60)
    if (hrs < 24) return `${hrs}h ago`
    return `${Math.floor(hrs / 24)}d ago`
}

function showToast(msg, type = 'success') {
    const toast = document.getElementById('toast')
    toast.textContent = msg
    toast.className = `toast ${type} show`
    setTimeout(() => toast.classList.remove('show'), 3000)
}

init()