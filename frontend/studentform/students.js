const API = 'http://localhost:8080/api'

const params = new URLSearchParams(window.location.search)
const classId = params.get('class')

let allStudents = []
let feeStatusMap = {}  // student_id → {status, total_paid}
let selectedRolls = new Set()
let selectionMode = false
let longPressTimer = null
let longPressFired = false

const now = new Date()
const currentMonth = now.toLocaleString('default', { month: 'long' })
const currentYear = now.getFullYear()

function showToast(msg, type = 'success') {
    const toast = document.getElementById('toast')
    toast.textContent = msg
    toast.className = `toast ${type} show`
    setTimeout(() => toast.classList.remove('show'), 3000)
}

// ── Selection Mode ──────────────────────────────────────────

function enterSelectionMode(triggerRollNo = null) {
    selectionMode = true
    selectedRolls.clear()

    if (triggerRollNo) selectedRolls.add(triggerRollNo)

    document.getElementById('normalHeader').classList.add('hidden')
    document.getElementById('selectHeader').classList.remove('hidden')
    document.body.classList.add('selection-active')

    renderStudents()
    updateSelectionBar()

    if (navigator.vibrate) navigator.vibrate(30)
}

function exitSelectionMode() {
    selectionMode = false
    selectedRolls.clear()

    document.getElementById('normalHeader').classList.remove('hidden')
    document.getElementById('selectHeader').classList.add('hidden')
    document.body.classList.remove('selection-active')

    renderStudents()
    updateSelectionBar()
}

function toggleSelect(rollNo) {
    if (selectedRolls.has(rollNo)) {
        selectedRolls.delete(rollNo)
    } else {
        selectedRolls.add(rollNo)
    }
    updateSelectionBar()
    updateRowCheckboxes()
}

function selectAll() {
    if (selectedRolls.size === allStudents.length) {
        selectedRolls.clear()
    } else {
        allStudents.forEach(s => selectedRolls.add(s.roll_no))
    }
    updateSelectionBar()
    updateRowCheckboxes()
}

function updateSelectionBar() {
    const count = selectedRolls.size
    const total = allStudents.length
    const countEl = document.getElementById('selectedCount')
    const dlBtn = document.getElementById('downloadBtn')
    const allBtn = document.getElementById('selectAllBtn')

    if (countEl) countEl.textContent = count === 0 ? 'Select students' : `${count} of ${total} selected`
    if (dlBtn) {
        dlBtn.disabled = count === 0
        dlBtn.classList.toggle('ready', count > 0)
    }

    if (allBtn) {
        allBtn.textContent = count === allStudents.length && total > 0 ? 'Deselect All' : 'Select All'
    }
}

function updateRowCheckboxes() {
    allStudents.forEach(s => {
        const cb = document.getElementById(`cb-${s.roll_no}`)
        const row = document.getElementById(`row-${s.roll_no}`)
        if (cb) cb.checked = selectedRolls.has(s.roll_no)
        if (row) row.classList.toggle('selected', selectedRolls.has(s.roll_no))
    })
}

// ── Long Press / Right Click ────────────────────────────────

function attachRowEvents(el, rollNo) {
    el.addEventListener('contextmenu', (e) => {
        e.preventDefault()
        if (!selectionMode) enterSelectionMode(rollNo)
    })

    el.addEventListener('touchstart', () => {
        longPressFired = false
        longPressTimer = setTimeout(() => {
            longPressFired = true
            if (!selectionMode) enterSelectionMode(rollNo)
        }, 500)
    }, { passive: true })

    el.addEventListener('touchend', () => clearTimeout(longPressTimer))
    el.addEventListener('touchmove', () => clearTimeout(longPressTimer))

    el.addEventListener('click', () => {
        if (longPressFired) return
        if (selectionMode) toggleSelect(rollNo)
    })
}

// ── Field Picker ─────────────────────────────────────────────

const ALL_FIELDS = [
    { key: 'roll_no', label: 'Roll No', default: true },
    { key: 'name', label: 'Name', default: true },
    { key: 'class', label: 'Class', default: true },
    { key: 'section', label: 'Section', default: true },
    { key: 'phone', label: 'Phone', default: true },
    { key: 'gender', label: 'Gender', default: true },
    { key: 'dob', label: 'Date of Birth', default: false },
    { key: 'aadhar_no', label: 'Aadhar No', default: false },
    { key: 'epunjab_id', label: 'ePunjab ID', default: false },
    { key: 'father_name', label: 'Father Name', default: false },
    { key: 'father_contact', label: 'Father Contact', default: false },
    { key: 'father_aadhar', label: 'Father Aadhar', default: false },
    { key: 'mother_name', label: 'Mother Name', default: false },
    { key: 'mother_contact', label: 'Mother Contact', default: false },
    { key: 'address', label: 'Address', default: false },
    { key: 'caste', label: 'Caste', default: false },
    { key: 'previous_school_details', label: 'Previous School', default: false },
]

function openFieldPicker() {
    const modal = document.getElementById('fieldPickerModal')
    const grid = document.getElementById('fieldGrid')

    grid.innerHTML = ALL_FIELDS.map(f => `
        <label class="field-toggle ${f.default ? 'checked' : ''}" id="ftl-${f.key}">
            <input type="checkbox" id="ft-${f.key}" ${f.default ? 'checked' : ''}
                   onchange="updateFieldToggle('${f.key}')">
            <span class="field-toggle-label">${f.label}</span>
            <span class="field-checkmark">✓</span>
        </label>
    `).join('')

    updateFieldPickerCount()
    modal.classList.add('show')
}

function closeFieldPicker() {
    document.getElementById('fieldPickerModal').classList.remove('show')
}

function updateFieldToggle(key) {
    const cb = document.getElementById(`ft-${key}`)
    const label = document.getElementById(`ftl-${key}`)
    label.classList.toggle('checked', cb.checked)
    updateFieldPickerCount()
}

function toggleAllFields() {
    const anyUnchecked = ALL_FIELDS.some(f => !document.getElementById(`ft-${f.key}`).checked)
    ALL_FIELDS.forEach(f => {
        const cb = document.getElementById(`ft-${f.key}`)
        const label = document.getElementById(`ftl-${f.key}`)
        cb.checked = anyUnchecked
        label.classList.toggle('checked', anyUnchecked)
    })
    updateFieldPickerCount()
}

function updateFieldPickerCount() {
    const count = ALL_FIELDS.filter(f => document.getElementById(`ft-${f.key}`).checked).length
    const btn = document.getElementById('fieldToggleAllBtn')
    const dlBtn = document.getElementById('confirmDownloadBtn')
    if (btn) btn.textContent = count === ALL_FIELDS.length ? 'Deselect All' : 'Select All'
    if (dlBtn) {
        dlBtn.disabled = count === 0
        dlBtn.textContent = count === 0 ? 'Select at least one field' : `Download ${selectedRolls.size} student${selectedRolls.size > 1 ? 's' : ''}`
    }
}

// ── CSV Export ──────────────────────────────────────────────

function downloadCSV() {
    if (selectedRolls.size === 0) return
    openFieldPicker()
}

function confirmDownload() {
    const selected = allStudents.filter(s => selectedRolls.has(s.roll_no))
    const activeFields = ALL_FIELDS.filter(f => document.getElementById(`ft-${f.key}`)?.checked)

    if (selected.length === 0 || activeFields.length === 0) return

    const headers = activeFields.map(f => f.label)
    const rows = selected.map(s =>
        activeFields.map(f => `"${(s[f.key] || '').toString().replace(/"/g, '""')}"`)
    )

    const csv = [headers.join(','), ...rows.map(r => r.join(','))].join('\n')
    const blob = new Blob([csv], { type: 'text/csv' })
    const url = URL.createObjectURL(blob)

    const a = document.createElement('a')
    a.href = url
    a.download = `students-class${classId}-${Date.now()}.csv`
    a.click()
    URL.revokeObjectURL(url)

    closeFieldPicker()
    showToast(`Downloaded ${selected.length} student${selected.length > 1 ? 's' : ''}`)
    exitSelectionMode()
}

// ── Fee Badge ─────────────────────────────────────────────────

function getFeeBadge(studentId) {
    const fee = feeStatusMap[studentId]
    if (!fee) {
        return `<span class="fee-badge-sm unpaid"
            onclick="event.stopPropagation(); window.location.href='fee-collect.html?student_id=${studentId}&month=${currentMonth}&year=${currentYear}'">
            Unpaid
        </span>`
    }
    return `<span class="fee-badge-sm ${fee.status}"
        onclick="event.stopPropagation(); window.location.href='fees.html'">
        ${fee.status === 'paid' ? `₹${fee.total_paid.toLocaleString()}` : fee.status}
    </span>`
}

// ── Render ───────────────────────────────────────────────────

function renderStudents() {
    const list = document.getElementById('studentList')

    if (!allStudents || allStudents.length === 0) {
        list.innerHTML = '<div class="empty-state">No students in this class yet.</div>'
        return
    }

    list.innerHTML = allStudents.map(s => `
        <div class="student-row ${selectedRolls.has(s.roll_no) ? 'selected' : ''}" id="row-${s.roll_no}">

            <div class="checkbox-wrap ${selectionMode ? 'visible' : ''}">
                <label class="custom-checkbox">
                    <input type="checkbox" id="cb-${s.roll_no}"
                           ${selectedRolls.has(s.roll_no) ? 'checked' : ''}
                           onchange="toggleSelect('${s.roll_no}')">
                    <span class="checkmark"></span>
                </label>
            </div>

            <div class="student-roll">${s.roll_no || '—'}</div>

            <div class="student-info">
                <div class="student-name">${s.name || '—'}</div>
                <div class="student-meta">${s.gender || ''} ${s.dob ? '· ' + s.dob : ''}</div>
            </div>

            ${!selectionMode ? getFeeBadge(s.ID) : ''}

            <a class="student-phone" href="tel:${s.phone}" onclick="event.stopPropagation()">
                ${s.phone || '—'}
            </a>

            <button class="btn-view ${selectionMode ? 'hidden' : ''}"
                    onclick="event.stopPropagation(); window.location.href='student-detail.html?roll_no=${s.roll_no}'">
                View
            </button>
        </div>
    `).join('')

    // attach events after render
    allStudents.forEach(s => {
        const row = document.getElementById(`row-${s.roll_no}`)
        if (row) attachRowEvents(row, s.roll_no)
    })
}

// ── Load ─────────────────────────────────────────────────────

async function loadFeeStatus() {
    try {
        const res = await fetch(`${API}/fees/class/${classId}/month/${currentMonth}/year/${currentYear}`)
        const data = await res.json()
        feeStatusMap = {}
        if (data) {
            data.forEach(s => {
                feeStatusMap[s.student_id] = {
                    status: s.has_paid ? (s.fees?.[0]?.status || 'paid') : 'unpaid',
                    total_paid: s.total_paid || 0
                }
            })
        }
    } catch (e) {
        console.log('fee status load failed', e)
    }
}

async function loadClassInfo() {
    const res = await fetch(`${API}/classes`)
    const classes = await res.json()
    const cls = classes.find(c => c.class == classId)
    if (!cls) return

    document.getElementById('classSummary').innerHTML = `
        <div class="cs-class">Class ${cls.class} — ${cls.section || '—'}</div>
        <div class="cs-divider">|</div>
        <div class="cs-teacher">
            👤 ${cls.teacher_name || 'No teacher'}
            ${cls.teacher_contact ? `· <a href="tel:${cls.teacher_contact}">${cls.teacher_contact}</a>` : ''}
        </div>
        <div class="cs-divider">|</div>
        <div class="cs-fees">📚 ₹${cls.tuition_fee?.toLocaleString() || '—'}/mo</div>
    `
}

async function loadStudents() {
    const res = await fetch(`${API}/students/class/${classId}`)
    const students = await res.json()
    allStudents = students || []

    const subtitle = document.getElementById('pageSubtitle')
    if (subtitle) subtitle.textContent =
        `${allStudents.length} student${allStudents.length !== 1 ? 's' : ''} · ${currentMonth} ${currentYear}`

    await loadFeeStatus()
    renderStudents()
}

loadClassInfo()
loadStudents()