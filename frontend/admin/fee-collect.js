requireAuth()

const API = CONFIG.API

let allStudents = []
let selectedStudent = null
let searchTimeout = null

// ── Init ─────────────────────────────────────────────────────

function init() {
    const now = new Date()
    const currentMonth = now.toLocaleString('default', { month: 'long' })
    const currentYear = now.getFullYear()

    // Year options
    const yearSel = document.getElementById('feeYear')
    for (let y = currentYear; y >= currentYear - 3; y--) {
        const opt = document.createElement('option')
        opt.value = y
        opt.textContent = y
        yearSel.appendChild(opt)
    }

    // Default month
    document.getElementById('feeMonth').value = currentMonth

    // Pre-fill from URL params
    const params = new URLSearchParams(window.location.search)
    const studentId = params.get('student_id')
    const month = params.get('month')
    const year = params.get('year')

    if (month) document.getElementById('feeMonth').value = month
    if (year) document.getElementById('feeYear').value = year

    loadAllStudents(studentId)
}

async function loadAllStudents(preSelectId = null) {
    const classRes = await authFetch(`${API}/classes`)
    const classes = await classRes.json()

    const promises = classes.map(c =>
        authFetch(`${API}/students/class/${c.class}?section=${encodeURIComponent(c.section || '')}`).then(r => r.json())
    )
    const results = await Promise.all(promises)
    allStudents = results.flat().filter(Boolean)

    // classMap keyed by "class|section" for exact lookup
    window.classMap = {}
    classes.forEach(c => {
        window.classMap[`${c.class}|${c.section || ''}`] = c
        window.classMap[c.class] = c // fallback for single-section classes
    })

    if (preSelectId) {
        const s = allStudents.find(s => s.id == preSelectId)
        if (s) {
            selectStudent(s)
            await checkExistingFee(preSelectId)
        }
    }
}

async function checkExistingFee(studentId) {
    const month = document.getElementById('feeMonth').value
    const year = document.getElementById('feeYear').value

    const res = await authFetch(`${API}/fees/student/${studentId}`)
    if (!res || !res.ok) return
    const fees = await res.json() || []

    // month is returned as name string ("April") from new handler
    const existing = fees.find(f =>
        f.month === month && String(f.year) === String(year)
    )

    if (!existing) return

    if (existing.status === 'paid') {
        showToast('Fee already fully paid for this month', 'error')
        document.getElementById('feeForm').querySelectorAll('input, select, button[type="submit"]')
            .forEach(el => el.disabled = true)
        document.getElementById('paymentStatus').className = 'payment-status paid'
        document.getElementById('paymentStatus').textContent = '✓ Fee already fully paid for this month'
        document.getElementById('paymentStatus').style.display = 'block'
        return
    }

    if (existing.status === 'partial') {
        document.getElementById('baseAmount').value = existing.final_amount
        document.getElementById('finalAmount').value = existing.remaining
        document.getElementById('paidAmount').value = existing.remaining
        document.getElementById('remaining').value = 0

        document.getElementById('paymentStatus').className = 'payment-status partial'
        document.getElementById('paymentStatus').textContent =
            `⚠ Partial payment exists — ₹${existing.remaining} remaining`
        document.getElementById('paymentStatus').style.display = 'block'

        window.existingFeeId = existing.id
    }
}

// ── Student Search ────────────────────────────────────────────

function searchStudents() {
    clearTimeout(searchTimeout)
    searchTimeout = setTimeout(() => {
        const q = document.getElementById('studentSearch').value.toLowerCase().trim()
        const dropdown = document.getElementById('searchDropdown')

        if (!q) {
            dropdown.classList.add('hidden')
            return
        }

        const matches = allStudents.filter(s =>
            s.name?.toLowerCase().includes(q) ||
            s.roll_no?.toLowerCase().includes(q) ||
            s.epunjab_id?.toLowerCase().includes(q)
        ).slice(0, 8)

        if (matches.length === 0) {
            dropdown.innerHTML = '<div class="dropdown-empty">No students found</div>'
        } else {
            dropdown.innerHTML = matches.map(s => `
                <div class="dropdown-item" onclick="selectStudent(${JSON.stringify(s).replace(/"/g, '&quot;')})">
                    <div class="dropdown-name">${esc(s.name)}</div>
                    <div class="dropdown-meta">Class ${s.class} · Roll ${s.roll_no || '—'} · ${s.epunjab_id || 'No ePunjab ID'}</div>
                </div>
            `).join('')
        }

        dropdown.classList.remove('hidden')
    }, 200)
}

function selectStudent(s) {
    selectedStudent = s

    // hide search, show card
    document.getElementById('studentSearch').value = ''
    document.getElementById('searchDropdown').classList.add('hidden')
    document.getElementById('selectedStudentCard').classList.remove('hidden')

    // fill card
    document.getElementById('ssAvatar').textContent = s.name ? s.name[0].toUpperCase() : '?'
    document.getElementById('ssName').textContent = s.name
    document.getElementById('ssMeta').textContent =
        `Class ${s.class}${s.section ? ' — ' + s.section : ''} · Roll ${s.roll_no || '—'} · ${s.epunjab_id || 'No ePunjab ID'}`

    // fill hidden fields
    document.getElementById("studentId").value = s.id
    document.getElementById('studentEpunjab').value = s.epunjab_id || ''
    document.getElementById('studentName').value = s.name
    document.getElementById('studentClass').value = s.class
    // store enrollment_id for new schema
    window.enrollmentId = s.enrollment_id || null

    // auto-fill base amount from class
    updateBaseAmount()
}

function clearStudent() {
    selectedStudent = null
    window.existingFeeId = null
    document.getElementById('selectedStudentCard').classList.add('hidden')
    document.getElementById('studentSearch').value = ''
    document.getElementById('studentId').value = ''
    document.getElementById('baseAmount').value = ''
    document.getElementById('finalAmount').value = ''
    document.getElementById('remaining').value = ''
    updatePaymentStatus()
}

function updateBaseAmount() {
    if (!selectedStudent) return
    const feeType = document.getElementById('feeType').value
    // try exact class|section key first, fall back to class only
    const key = `${selectedStudent.class}|${selectedStudent.section || ''}`
    const cls = window.classMap?.[key] || window.classMap?.[selectedStudent.class]
    if (!cls) return

    const base = feeType === 'tuition' ? cls.tuition_fee : cls.transport_fee
    document.getElementById('baseAmount').value = base || ''
    calcFinal()
}

// ── Calculations ──────────────────────────────────────────────

function calcFinal() {
    const base = parseInt(document.getElementById('baseAmount').value) || 0
    const discount = parseInt(document.getElementById('discount').value) || 0
    const final = base - discount
    document.getElementById('finalAmount').value = final > 0 ? final : 0
    calcRemaining()
}

function calcRemaining() {
    const final = parseInt(document.getElementById('finalAmount').value) || 0
    const paid = parseInt(document.getElementById('paidAmount').value) || 0
    const rem = final - paid
    document.getElementById('remaining').value = rem > 0 ? rem : 0
    updatePaymentStatus()
}

function updatePaymentStatus() {
    const final = parseInt(document.getElementById('finalAmount').value) || 0
    const paid = parseInt(document.getElementById('paidAmount').value) || 0
    const status = document.getElementById('paymentStatus')

    if (!final || !paid) {
        status.className = 'payment-status'
        status.textContent = ''
        return
    }

    if (paid >= final) {
        status.className = 'payment-status paid'
        status.textContent = '✓ Full payment — Receipt will be marked PAID'
    } else if (paid > 0) {
        status.className = 'payment-status partial'
        status.textContent = `⚠ Partial payment — ₹${final - paid} remaining`
    } else {
        status.className = 'payment-status unpaid'
        status.textContent = 'No amount entered'
    }
}

// ── Submit ────────────────────────────────────────────────────

document.getElementById('feeForm').addEventListener('submit', async function (e) {
    e.preventDefault()

    if (!selectedStudent) {
        showToast('Please select a student first', 'error')
        return
    }

    const paidAmount = parseInt(document.getElementById('paidAmount').value) || 0

    if (paidAmount <= 0) {
        showToast('Paid amount must be greater than 0', 'error')
        return
    }

    const submitBtn = this.querySelector('button[type="submit"]')
    submitBtn.disabled = true
    const origText = submitBtn.textContent
    submitBtn.textContent = 'Processing...'

    try {
        // if existing partial — update instead of create
        if (window.existingFeeId) {
            const res = await authFetch(`${API}/fees/${window.existingFeeId}/complete`, {
                method: 'PUT',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({
                    paid_amount: paidAmount,
                    due_date: document.getElementById('dueDate').value,
                })
            })

            if (res.ok) {
                const fee = await res.json()
                window.location.href = `fee-receipt.html?receipt=${fee.receipt_no}`
            } else {
                const err = await res.json()
                showToast('Error: ' + err.error, 'error')
            }
            return
        }

        // normal new payment flow
        const data = {
            student_id:     parseInt(document.getElementById('studentId').value),
            enrollment_id:  window.enrollmentId || null,
            epunjab_id:     document.getElementById('studentEpunjab').value,
            student_name:   document.getElementById('studentName').value,
            roll_no:        selectedStudent.roll_no || '',
            class:          document.getElementById('studentClass').value,
            section:        selectedStudent.section || '',
            month:          document.getElementById('feeMonth').value,
            year:           parseInt(document.getElementById('feeYear').value),
            fee_type:       document.getElementById('feeType').value,
            base_amount:    parseInt(document.getElementById('baseAmount').value) || 0,
            discount:       parseInt(document.getElementById('discount').value) || 0,
            discount_reason: document.getElementById('discountReason').value,
            paid_amount:    paidAmount,
            due_date:       document.getElementById('dueDate').value,
        }

        const res = await authFetch(`${API}/fees/pay`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(data)
        })

        if (res.ok) {
            const fee = await res.json()
            window.location.href = `fee-receipt.html?receipt=${fee.receipt_no}`
        } else {
            const err = await res.json()
            showToast('Error: ' + err.error, 'error')
        }
    } finally {
        // only re-enable on error path (success navigates away)
        submitBtn.disabled = false
        submitBtn.textContent = origText
    }
})

function showToast(msg, type = 'success') {
    const toast = document.getElementById('toast')
    toast.textContent = msg
    toast.className = `toast ${type} show`
    setTimeout(() => toast.classList.remove('show'), 3000)
}

// close dropdown on outside click
document.addEventListener('click', (e) => {
    if (!e.target.closest('.search-wrapper')) {
        document.getElementById('searchDropdown').classList.add('hidden')
    }
})

init()