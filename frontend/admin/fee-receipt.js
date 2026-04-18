requireAuth()

const API = CONFIG.API

const params = new URLSearchParams(window.location.search)
const receiptNo = params.get('receipt')

async function loadReceipt() {
    if (!receiptNo) {
        document.getElementById('receipt').innerHTML = '<p style="text-align:center;padding:20px;">No receipt found</p>'
        return
    }

    const res = await authFetch(`${API}/fees/receipt/${receiptNo}`)
    if (!res.ok) {
        document.getElementById('receipt').innerHTML = '<p style="text-align:center;padding:20px;">Receipt not found</p>'
        return
    }

    const f = await res.json()

    // Meta
    document.getElementById('rReceiptNo').textContent = f.receipt_no
    document.getElementById('rDate').textContent = f.paid_at ? f.paid_at.split('T')[0] : '—'

    // Student
    document.getElementById('rStudentName').textContent = f.student_name || '—'
    document.getElementById('rClass').textContent = f.class ? `Class ${f.class}` : '—'
    document.getElementById('rRollNo').textContent = '—' // not stored in fee, shown as placeholder
    document.getElementById('rEpunjab').textContent = f.epunjab_id || '—'

    if (!f.epunjab_id) {
        document.getElementById('epunjabRow').style.display = 'none'
    }

    // Fee
    document.getElementById('rFeeType').textContent = f.fee_type === 'tuition' ? 'Tuition Fee' : 'Transport Fee'
    document.getElementById('rMonth').textContent = `${f.month} ${f.year}`
    document.getElementById('rBaseAmount').textContent = `₹${f.base_amount?.toLocaleString() || 0}`

    if (f.discount > 0) {
        document.getElementById('rDiscount').textContent =
            `₹${f.discount.toLocaleString()} (${f.discount_reason || ''})`
    } else {
        document.getElementById('discountRow').style.display = 'none'
    }

    document.getElementById('rFinalAmount').textContent = `₹${f.final_amount?.toLocaleString() || 0}`

    // Payment
    document.getElementById('rPaidAmount').textContent = `₹${f.paid_amount?.toLocaleString() || 0}`

    if (f.remaining > 0) {
        document.getElementById('rRemaining').textContent = `₹${f.remaining.toLocaleString()}`
        if (f.due_date) {
            document.getElementById('rDueDate').textContent = f.due_date
        } else {
            document.getElementById('dueDateRow').style.display = 'none'
        }
    } else {
        document.getElementById('remainingRow').style.display = 'none'
        document.getElementById('dueDateRow').style.display = 'none'
    }

    // Status
    const statusEl = document.getElementById('rStatus')
    statusEl.textContent = f.status?.toUpperCase() || 'PAID'
    statusEl.className = `receipt-status ${f.status || 'paid'}`

    // Auto print if coming fresh from payment
    const autoPrint = new URLSearchParams(window.location.search).get('print')
    if (autoPrint === '1') {
        setTimeout(() => window.print(), 800)
    }
}

loadReceipt()