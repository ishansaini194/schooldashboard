requireAuth()

const API = CONFIG.API
const studentID = localStorage.getItem('student_id')
const year = new Date().getFullYear()

async function init() {
    document.getElementById('yearLabel').textContent = `Session ${year}-${year + 1}`

    const res = await authFetch(`${API}/fees/student/${studentID}/yearly?year=${year}`)
    if (!res || !res.ok) return
    const data = await res.json()

    document.getElementById('paidCount').textContent = data.paid_count || 0
    document.getElementById('totalMonths').textContent = data.total_months || 0

    const months = data.months || []
    if (months.length === 0) {
        document.getElementById('feeGrid').innerHTML = '<div class="dash-empty">No fee records found</div>'
        return
    }

    document.getElementById('feeGrid').innerHTML = months.map(m => `
        <div class="fee-month-card ${m.status}">
            <div class="fee-month-name">${m.month.slice(0, 3)}</div>
            <div class="fee-month-status-icon">${m.status === 'paid' ? '✓' : m.status === 'partial' ? '~' : '✗'}</div>
            ${m.has_record && m.paid_amount > 0
            ? `<div class="fee-month-amount">₹${m.paid_amount.toLocaleString()}</div>`
            : ''}
            ${m.remaining > 0
            ? `<div class="fee-month-remaining">₹${m.remaining} left</div>`
            : ''}
        </div>
    `).join('')
}

init()