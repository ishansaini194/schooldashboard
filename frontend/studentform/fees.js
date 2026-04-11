const API = 'http://localhost:8080/api'

let allFeeData = []
let currentTab = 'all'

// ── Init ─────────────────────────────────────────────────────

function init() {
    const now = new Date()
    const currentMonth = now.toLocaleString('default', { month: 'long' })
    const currentYear = now.getFullYear()

    // Set year options
    const yearSel = document.getElementById('yearFilter')
    for (let y = currentYear; y >= currentYear - 3; y--) {
        const opt = document.createElement('option')
        opt.value = y
        opt.textContent = y
        yearSel.appendChild(opt)
    }

    // Set current month
    document.getElementById('monthFilter').value = currentMonth

    loadClasses()
}

async function loadClasses() {
    const res = await fetch(`${API}/classes`)
    const classes = await res.json()
    const sel = document.getElementById('classFilter')

    classes.forEach(c => {
        const opt = document.createElement('option')
        opt.value = c.class
        opt.textContent = `Class ${c.class} — ${c.section || ''}`
        sel.appendChild(opt)
    })
}

// ── Load Fee Status ───────────────────────────────────────────

async function loadFeeStatus() {
    const cls   = document.getElementById('classFilter').value
    const month = document.getElementById('monthFilter').value
    const year  = document.getElementById('yearFilter').value

    if (!cls) {
        document.getElementById('feeList').innerHTML =
            '<div class="empty-state">Select a class to view fee status</div>'
        return
    }

    document.getElementById('pageSubtitle').textContent = `Class ${cls} · ${month} ${year}`

    const res = await fetch(`${API}/fees/class/${cls}/month/${month}/year/${year}`)
    allFeeData = await res.json() || []

    updateStats()
    renderList()
}

function updateStats() {
    const total   = allFeeData.length
    const paid    = allFeeData.filter(s => s.has_paid && s.fees?.every(f => f.status === 'paid')).length
    const pending = total - paid
    const amount  = allFeeData.reduce((sum, s) => sum + (s.total_paid || 0), 0)

    document.getElementById('statTotal').textContent   = total
    document.getElementById('statPaid').textContent    = paid
    document.getElementById('statPending').textContent = pending
    document.getElementById('statAmount').textContent  = `₹${amount.toLocaleString()}`
}

function renderList() {
    const search = document.getElementById('searchInput').value.toLowerCase()
    const list   = document.getElementById('feeList')

    let data = allFeeData

    // filter by tab
    if (currentTab === 'paid') {
        data = data.filter(s => s.has_paid)
    } else if (currentTab === 'pending') {
        data = data.filter(s => !s.has_paid)
    }

    // filter by search
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

    const cls   = document.getElementById('classFilter').value
    const month = document.getElementById('monthFilter').value
    const year  = document.getElementById('yearFilter').value

    list.innerHTML = data.map(s => {
        const hasPaid   = s.has_paid
        const totalPaid = s.total_paid || 0
        const latestReceipt = s.fees?.length > 0 ? s.fees[0].receipt_no : null
        const status = s.fees?.length > 0 ? s.fees[0].status : 'unpaid'

        return `
        <div class="fee-row ${hasPaid ? 'paid-row' : 'pending-row'}">
            <div class="fee-roll">${s.roll_no || '—'}</div>
            <div class="fee-info">
                <div class="fee-name">${s.student_name}</div>
                <div class="fee-meta">${hasPaid ? `Paid ₹${totalPaid.toLocaleString()}` : 'Not paid'}</div>
            </div>
            <span class="fee-badge ${status}">${status}</span>
            ${hasPaid && latestReceipt
                ? `<button class="btn-view-receipt"
                    onclick="window.location.href='fee-receipt.html?receipt=${latestReceipt}'">
                    View Receipt
                   </button>`
                : `<button class="btn-collect"
                    onclick="window.location.href='fee-collect.html?student_id=${s.student_id}&month=${month}&year=${year}'">
                    Collect Fee
                   </button>`
            }
        </div>
        `
    }).join('')
}

function filterList() {
    renderList()
}

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

init()