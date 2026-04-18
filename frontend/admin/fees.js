requireAuth()

const API = CONFIG.API

// ── State ─────────────────────────────────────────────────────
let allFeeData = []
let allClasses = []
let currentTab = 'all'
let activePanelId = null

// ── State Persistence ─────────────────────────────────────────

function saveState() {
    localStorage.setItem('fees_class', document.getElementById('classFilter').value)
    localStorage.setItem('fees_month', document.getElementById('monthFilter').value)
    localStorage.setItem('fees_year', document.getElementById('yearFilter').value)
}

function restoreState() {
    const urlParams = new URLSearchParams(window.location.search)
    const urlClass = urlParams.get('class')
    const urlSection = urlParams.get('section')
    const urlMonth = urlParams.get('month')
    const urlYear = urlParams.get('year')
    const urlStudent = urlParams.get('student_id')

    let cls = localStorage.getItem('fees_class')
    if (urlClass) {
        // build "class|section" — section from URL if present, else first matching class
        if (urlSection) {
            cls = `${urlClass}|${urlSection}`
        } else {
            const match = allClasses.find(c => String(c.class) === String(urlClass))
            cls = match ? `${match.class}|${match.section || ''}` : urlClass
        }
    }
    const month = urlMonth || localStorage.getItem('fees_month')
    const year = urlYear || localStorage.getItem('fees_year')

    if (cls) document.getElementById('classFilter').value = cls
    if (month) document.getElementById('monthFilter').value = month
    if (year) document.getElementById('yearFilter').value = year

    return urlStudent
}

// ── Init ──────────────────────────────────────────────────────

async function init() {
    const now = new Date()
    const currentMonth = now.toLocaleString('default', { month: 'long' })
    const currentYear = now.getFullYear()

    const yearSel = document.getElementById('yearFilter')
    for (let y = currentYear; y >= currentYear - 3; y--) {
        const opt = document.createElement('option')
        opt.value = y
        opt.textContent = y
        yearSel.appendChild(opt)
    }

    document.getElementById('monthFilter').value = currentMonth

    await loadClasses()

    const highlightStudentId = restoreState()

    const urlTab = new URLSearchParams(window.location.search).get('tab')
    if (urlTab) switchTab(urlTab)

    await loadFeeStatus()
    if (highlightStudentId) highlightStudent(highlightStudentId)
}

// ── Data Fetching ─────────────────────────────────────────────

async function loadClasses() {
    const res = await authFetch(`${API}/classes`)
    allClasses = await res.json() || []
    const sel = document.getElementById('classFilter')

    allClasses.forEach(c => {
        const opt = document.createElement('option')
        opt.value = `${c.class}|${c.section || ''}`
        opt.textContent = `Class ${c.class}${c.section ? ' — ' + c.section : ''}`
        sel.appendChild(opt)
    })
}

async function loadFeeStatus() {
    const clsRaw = document.getElementById('classFilter').value
    const month = document.getElementById('monthFilter').value
    const year = document.getElementById('yearFilter').value

    saveState()

    if (clsRaw) {
        // single class — split "6|A" into class + section
        const [cls, section] = clsRaw.split('|')
        const sectionLabel = section ? ` — ${section}` : ''
        document.getElementById('pageSubtitle').textContent = `Class ${cls}${sectionLabel} · ${month} ${year}`
        const sectionParam = section ? `?section=${encodeURIComponent(section)}` : ''
        const res = await authFetch(`${API}/fees/class/${cls}/month/${month}/year/${year}${sectionParam}`)
        allFeeData = await res.json() || []
    } else {
        // all classes — fetch each (with section) and merge
        document.getElementById('pageSubtitle').textContent = `All Classes · ${month} ${year}`
        const results = await Promise.all(
            allClasses.map(c => {
                const sectionParam = c.section ? `?section=${encodeURIComponent(c.section)}` : ''
                return authFetch(`${API}/fees/class/${c.class}/month/${month}/year/${year}${sectionParam}`)
                    .then(r => r.json())
                    .then(data => (data || []).map(s => ({ ...s, class: c.class, section: c.section })))
            })
        )
        allFeeData = results.flat()
    }

    updateStats()
    renderList()
}

async function fetchStudentHistory(studentId) {
    const year = new Date().getFullYear()
    const res = await authFetch(`${API}/fees/student/${studentId}/yearly?year=${year}`)
    const data = await res.json()
    return data
}

// ── Render Functions ──────────────────────────────────────────

function updateStats() {
    const total = allFeeData.length
    const paid = allFeeData.filter(s => s.has_paid).length
    const pending = total - paid
    const amount = allFeeData.reduce((sum, s) => sum + (s.total_paid || 0), 0)

    document.getElementById('statTotal').textContent = total
    document.getElementById('statPaid').textContent = paid
    document.getElementById('statPending').textContent = pending
    document.getElementById('statAmount').textContent = `₹${amount.toLocaleString()}`
}

function renderList() {
    const search = document.getElementById('searchInput').value.toLowerCase()
    const list = document.getElementById('feeList')
    const month = document.getElementById('monthFilter').value
    const year = document.getElementById('yearFilter').value
    const cls = document.getElementById('classFilter').value

    let data = [...allFeeData]

    // #2 fix — pending tab only shows unpaid/partial for selected class
    if (currentTab === 'paid') data = data.filter(s => s.has_paid)
    if (currentTab === 'pending') data = data.filter(s => !s.has_paid)

    if (search) {
        data = data.filter(s =>
            s.student_name?.toLowerCase().includes(search) ||
            s.roll_no?.toLowerCase().includes(search) ||
            s.fees?.some(f => f.epunjab_id?.toLowerCase().includes(search))
        )
    }

    if (data.length === 0) {
        list.innerHTML = '<div class="empty-state">No students found</div>'
        return
    }

    // show class column when viewing all classes
    const showClass = !cls

    list.innerHTML = data.map(s => renderFeeRow(s, month, year, showClass)).join('')
}

// → FeeRow component
function renderFeeRow(s, month, year, showClass = false) {
    const hasPaid = s.has_paid
    const totalPaid = s.total_paid || 0
    const latestReceipt = s.fees?.length > 0 ? s.fees[0].receipt_no : null
    const status = s.fees?.length > 0 ? s.fees[0].status : 'unpaid'

    let actionBtns = ''
    if (!hasPaid) {
        actionBtns = `<button class="btn-collect"
            onclick="event.stopPropagation(); window.location.href='fee-collect.html?student_id=${s.student_id}&month=${month}&year=${year}'">
            Collect Fee
        </button>`
    } else if (status === 'partial') {
        actionBtns = `
            <button class="btn-collect"
                onclick="event.stopPropagation(); window.location.href='fee-collect.html?student_id=${s.student_id}&month=${month}&year=${year}'">
                Pay Remaining
            </button>
            <button class="btn-view-receipt"
                onclick="event.stopPropagation(); window.location.href='fee-receipt.html?receipt=${latestReceipt}'">
                View Receipt
            </button>`
    } else {
        actionBtns = `<button class="btn-view-receipt"
            onclick="event.stopPropagation(); window.location.href='fee-receipt.html?receipt=${latestReceipt}'">
            View Receipt
        </button>`
    }

    return `
    <div class="fee-row ${hasPaid ? 'paid-row' : 'pending-row'} ${activePanelId == s.student_id ? 'panel-active' : ''}"
         id="feerow-${s.student_id}"
         onclick="openPanel(${s.student_id}, '${esc(s.student_name)}', '${s.roll_no}', '${month}', ${year})">
        <div class="fee-roll">${s.roll_no || '—'}</div>
        <div class="fee-info">
            <div class="fee-name">${esc(s.student_name)}${showClass ? ` <span class="fee-class-tag">Class ${s.class}${s.section ? '-' + s.section : ''}</span>` : ''}</div>
            <div class="fee-meta">${hasPaid ? `Paid ₹${totalPaid.toLocaleString()}` : 'Not paid'}</div>
        </div>
        <span class="fee-badge ${status}">${status}</span>
        <div class="fee-actions" onclick="event.stopPropagation()">${actionBtns}</div>
    </div>`
}

// ── Export ────────────────────────────────────────────────────

function getVisibleData() {
    const search = document.getElementById('searchInput').value.toLowerCase()
    let data = [...allFeeData]
    if (currentTab === 'paid') data = data.filter(s => s.has_paid)
    if (currentTab === 'pending') data = data.filter(s => !s.has_paid)
    if (search) {
        data = data.filter(s =>
            s.student_name?.toLowerCase().includes(search) ||
            s.roll_no?.toLowerCase().includes(search)
        )
    }
    return data
}

function exportCSV() {
    const data = getVisibleData()
    if (data.length === 0) { showToast('No data to export', 'error'); return }

    const month = document.getElementById('monthFilter').value
    const year = document.getElementById('yearFilter').value
    const cls = document.getElementById('classFilter').value || 'All'

    const headers = ['Roll No', 'Name', 'Class', 'Status', 'Paid Amount', 'Remaining', 'Receipt No']
    const rows = data.map(s => {
        const fee = s.fees?.[0] || {}
        return [
            s.roll_no, s.student_name, s.class || cls,
            fee.status || 'unpaid',
            s.total_paid || 0,
            fee.remaining || 0,
            fee.receipt_no || ''
        ].map(v => `"${(v || '').toString().replace(/"/g, '""')}"`)
    })

    const csv = [headers.join(','), ...rows.map(r => r.join(','))].join('\n')
    const blob = new Blob([csv], { type: 'text/csv' })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `fees-${cls}-${month}-${year}.csv`
    a.click()
    URL.revokeObjectURL(url)
    showToast('CSV downloaded!')
}

function exportPDF() {
    const data = getVisibleData()
    if (data.length === 0) { showToast('No data to export', 'error'); return }

    const month = document.getElementById('monthFilter').value
    const year = document.getElementById('yearFilter').value
    const cls = document.getElementById('classFilter').value || 'All Classes'
    const tab = currentTab === 'all' ? 'All Students' : currentTab === 'paid' ? 'Paid' : 'Pending'

    const printArea = document.getElementById('pdfPrintArea')
    printArea.innerHTML = `
        <div class="pdf-header">
            <div class="pdf-school">KRB School</div>
            <div class="pdf-title">Fee Report — ${cls === 'All Classes' ? 'All Classes' : 'Class ' + cls} · ${month} ${year}</div>
            <div class="pdf-subtitle">${tab} · ${data.length} students</div>
            <div class="pdf-date">Generated: ${new Date().toLocaleDateString('en-IN')}</div>
        </div>
        <table class="pdf-table">
            <thead>
                <tr>
                    <th>Roll No</th>
                    <th>Name</th>
                    <th>Class</th>
                    <th>Status</th>
                    <th>Paid</th>
                    <th>Remaining</th>
                    <th>Receipt</th>
                </tr>
            </thead>
            <tbody>
                ${data.map(s => {
        const fee = s.fees?.[0] || {}
        return `<tr>
                        <td>${s.roll_no || '—'}</td>
                        <td>${esc(s.student_name)}</td>
                        <td>${s.class || cls}</td>
                        <td class="${fee.status || 'unpaid'}">${fee.status || 'unpaid'}</td>
                        <td>₹${(s.total_paid || 0).toLocaleString()}</td>
                        <td>₹${(fee.remaining || 0).toLocaleString()}</td>
                        <td>${fee.receipt_no || '—'}</td>
                    </tr>`
    }).join('')}
            </tbody>
        </table>
        <div class="pdf-footer">
            Total Collected: ₹${data.reduce((sum, s) => sum + (s.total_paid || 0), 0).toLocaleString()}
            &nbsp;·&nbsp; ${data.length} students
        </div>
    `

    printArea.classList.remove('no-print')
    document.body.classList.add('print-pdf-mode')
    window.print()
    setTimeout(() => {
        printArea.classList.add('no-print')
        document.body.classList.remove('print-pdf-mode')
    }, 1000)
}

// ── Side Panel ────────────────────────────────────────────────

async function openPanel(studentId, name, rollNo, month, year) {
    activePanelId = studentId

    // name was passed through esc() in the row HTML, so decode entities for display
    const decoded = name ? new DOMParser().parseFromString(name, 'text/html').body.textContent : ''

    document.querySelectorAll('.fee-row').forEach(r => r.classList.remove('panel-active'))
    const activeRow = document.getElementById(`feerow-${studentId}`)
    if (activeRow) activeRow.classList.add('panel-active')

    document.getElementById('panelAvatar').textContent = decoded ? decoded[0].toUpperCase() : '?'
    document.getElementById('panelName').textContent = decoded || '—'
    document.getElementById('panelMeta').textContent = rollNo ? `Roll No: ${rollNo}` : '—'

    document.getElementById('panelCollectBtn').onclick = () => {
        window.location.href = `fee-collect.html?student_id=${studentId}&month=${month}&year=${year}`
    }

    document.getElementById('panelHistory').innerHTML = '<div class="panel-loading">Loading history...</div>'
    document.getElementById('panelTotalPaid').textContent = '—'
    document.getElementById('panelTotalPending').textContent = '—'
    document.getElementById('panelMonthsCount').textContent = '—'

    document.getElementById('sidePanel').classList.add('open')
    document.getElementById('panelOverlay').classList.add('show')

    const fees = await fetchStudentHistory(studentId)
    renderPanelHistory(fees)
}

function closePanel() {
    activePanelId = null
    document.getElementById('sidePanel').classList.remove('open')
    document.getElementById('panelOverlay').classList.remove('show')
    document.querySelectorAll('.fee-row').forEach(r => r.classList.remove('panel-active'))
}

function renderPanelHistory(data) {
    const historyEl = document.getElementById('panelHistory')

    if (!data || !data.months) {
        historyEl.innerHTML = '<div class="panel-empty">No history found</div>'
        document.getElementById('panelTotalPaid').textContent = '₹0'
        document.getElementById('panelTotalPending').textContent = '₹0'
        document.getElementById('panelMonthsCount').textContent = '0/0'
        return
    }

    const months = data.months
    const paidCount = data.paid_count || 0
    const totalMonths = data.total_months || 0

    // calculate totals
    const totalPaid = months.reduce((sum, m) => sum + (m.paid_amount || 0), 0)
    const totalPending = months.reduce((sum, m) => sum + (m.remaining || 0), 0)

    document.getElementById('panelTotalPaid').textContent = `₹${totalPaid.toLocaleString()}`
    document.getElementById('panelTotalPending').textContent = `₹${totalPending.toLocaleString()}`
    document.getElementById('panelMonthsCount').textContent = `${paidCount}/${totalMonths}`

    historyEl.innerHTML = months.map(m => renderHistoryItem(m)).join('')
}

function renderHistoryItem(m) {
    const statusClass = m.status || 'unpaid'
    const feeTypeLabel = m.fee_type === 'transport' ? '🚌' : '📚'
    const hasRecord = m.has_record

    let dueDateHtml = ''
    if (m.status === 'partial' && m.due_date) {
        const overdue = new Date(m.due_date) < new Date()
        dueDateHtml = `<div class="history-due ${overdue ? 'overdue' : 'upcoming'}">
            ${overdue ? '⚠ Overdue' : '📅 Due'}: ${m.due_date}
        </div>`
    }

    // clickable only for unpaid and partial
    const isClickable = statusClass === 'unpaid' || statusClass === 'partial'
    const studentId = activePanelId
    const month = m.month
    const year = new Date().getFullYear()

    const clickAttr = isClickable
        ? `onclick="window.location.href='fee-collect.html?student_id=${studentId}&month=${month}&year=${year}'" style="cursor:pointer"`
        : ''

    const clickHint = isClickable
        ? `<div class="history-pay-hint">Tap to pay →</div>`
        : ''

    return `
    <div class="history-item ${statusClass} ${isClickable ? 'history-clickable' : ''}" ${clickAttr}>
        <div class="history-left">
            <span class="history-type">${hasRecord ? feeTypeLabel : '—'}</span>
            <div>
                <div class="history-month">${m.month}</div>
                ${m.receipt_no ? `<div class="history-receipt">${m.receipt_no}</div>` : ''}
                ${clickHint}
            </div>
        </div>
        <div class="history-right">
            <div class="history-amount">${hasRecord ? `₹${(m.paid_amount || 0).toLocaleString()}` : '—'}</div>
            <span class="history-badge ${statusClass}">${m.status}</span>
            ${m.remaining > 0 ? `<div class="history-remaining">₹${m.remaining.toLocaleString()} left</div>` : ''}
            ${dueDateHtml}
        </div>
    </div>`
}

// ── Helpers ───────────────────────────────────────────────────

function highlightStudent(studentId) {
    setTimeout(() => {
        const row = document.getElementById(`feerow-${studentId}`)
        if (row) {
            row.classList.add('highlighted')
            row.scrollIntoView({ behavior: 'smooth', block: 'center' })
            setTimeout(() => row.classList.remove('highlighted'), 3000)
        }
    }, 300)
}

function filterList() { renderList() }

function switchTab(tab) {
    currentTab = tab
    document.querySelectorAll('.tab').forEach(t => t.classList.remove('active'))
    document.getElementById(`tab${tab.charAt(0).toUpperCase() + tab.slice(1)}`).classList.add('active')
    renderList()
}

function showToast(msg, type = 'success') {
    const toast = document.getElementById('toast')
    toast.textContent = msg
    toast.className = `toast ${type} show`
    setTimeout(() => toast.classList.remove('show'), 3000)
}

document.addEventListener('keydown', e => {
    if (e.key === 'Escape') closePanel()
})

init()