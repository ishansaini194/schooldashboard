requireAuth()

const API = CONFIG.API

function showToast(msg, type = 'success') {
    const toast = document.getElementById('toast')
    toast.textContent = msg
    toast.className = `toast ${type} show`
    setTimeout(() => toast.classList.remove('show'), 3000)
}

function openModal(cls = null) {
    document.getElementById('modal').classList.add('show')
    if (cls) {
        document.getElementById('modalTitle').textContent  = 'Edit Class'
        document.getElementById('classId').value          = cls.id
        document.getElementById('classNo').value          = cls.class
        document.getElementById('section').value          = cls.section
        document.getElementById('teacherName').value      = cls.teacher_name
        document.getElementById('teacherContact').value   = cls.teacher_contact
        document.getElementById('tuitionFee').value       = cls.tuition_fee || 0
        document.getElementById('transportFee').value     = cls.transport_fee || 0
    } else {
        document.getElementById('modalTitle').textContent = 'Add Class'
        document.getElementById('classForm').reset()
        document.getElementById('classId').value = ''
    }
}

function closeModal() {
    document.getElementById('modal').classList.remove('show')
}

async function loadClasses() {
    const res     = await authFetch(`${API}/classes`)
    const classes = await res.json()
    const grid    = document.getElementById('classGrid')

    if (!classes || classes.length === 0) {
        grid.innerHTML = '<div class="empty-state">No classes yet. Add one!</div>'
        return
    }

    grid.innerHTML = classes.map(cls => `
        <div class="class-card" onclick="window.location.href='students.html?class=${cls.class}&section=${encodeURIComponent(cls.section || '')}'" style="cursor:pointer">
            <div class="class-badge">Class ${cls.class}${cls.section ? ' — ' + cls.section : ''}</div>
            <div class="class-teacher">
                <span class="teacher-icon">👤</span>
                <div>
                    <div class="teacher-name">${esc(cls.teacher_name || 'No teacher assigned')}</div>
                    <div class="teacher-contact">${cls.teacher_contact || ''}</div>
                </div>
            </div>
            <div class="card-actions">
                <button class="btn-edit" onclick='event.stopPropagation(); openModal(${JSON.stringify(cls)})'>Edit</button>
                <button class="btn-delete" onclick='event.stopPropagation(); deleteClass(${cls.id})'>Delete</button>
            </div>
        </div>
    `).join('')
}

function showFormError(field, msg) {
    const input = document.getElementById(field)
    if (!input) return
    const wrap = input.closest('.field') || input.parentElement
    wrap.classList.add('has-error')
    wrap.querySelectorAll('.field-error').forEach(el => el.remove())
    const errEl = document.createElement('div')
    errEl.className = 'field-error'
    errEl.textContent = msg
    wrap.appendChild(errEl)
}

function clearFormErrors(formEl) {
    formEl.querySelectorAll('.field-error').forEach(el => el.remove())
    formEl.querySelectorAll('.has-error').forEach(el => el.classList.remove('has-error'))
}

function validateClassForm() {
    const errs = {}
    const classNo = document.getElementById('classNo').value.trim()
    const section = document.getElementById('section').value.trim()
    const tName = document.getElementById('teacherName').value.trim()
    const tContact = document.getElementById('teacherContact').value.trim()
    const tuition = document.getElementById('tuitionFee').value.trim()
    const transport = document.getElementById('transportFee').value.trim()

    if (!classNo) errs.classNo = 'Required'
    else if (!/^([1-9]|1[0-2])$/.test(classNo)) errs.classNo = 'Must be 1-12'

    if (!section) errs.section = 'Required'
    else if (!/^[A-Za-z]$/.test(section)) errs.section = 'Single letter A-Z'

    if (tName && !/^[A-Za-z][A-Za-z\s.'-]{1,49}$/.test(tName)) errs.teacherName = 'Letters/spaces/dots only'
    if (tContact && !/^\d{10}$/.test(tContact)) errs.teacherContact = '10 digits required'
    if (tuition && (!/^\d+$/.test(tuition) || parseInt(tuition) < 0)) errs.tuitionFee = 'Must be ≥ 0'
    if (transport && (!/^\d+$/.test(transport) || parseInt(transport) < 0)) errs.transportFee = 'Must be ≥ 0'

    return errs
}

document.getElementById('classForm').addEventListener('submit', async function(e) {
    e.preventDefault()
    clearFormErrors(this)

    const errs = validateClassForm()
    if (Object.keys(errs).length > 0) {
        for (const [field, msg] of Object.entries(errs)) showFormError(field, msg)
        showToast(`Please fix ${Object.keys(errs).length} error${Object.keys(errs).length > 1 ? 's' : ''} below`, 'error')
        return
    }

    const submitBtn = this.querySelector('button[type="submit"]')
    submitBtn.disabled = true
    const origText = submitBtn.textContent
    submitBtn.textContent = 'Saving...'

    try {
        const id   = document.getElementById('classId').value
        const data = {
            class:          parseInt(document.getElementById('classNo').value),
            section:        document.getElementById('section').value.toUpperCase(),
            teacher_name:   document.getElementById('teacherName').value.trim(),
            teacher_contact:document.getElementById('teacherContact').value.trim(),
            tuition_fee:    parseInt(document.getElementById('tuitionFee').value) || 0,
            transport_fee:  parseInt(document.getElementById('transportFee').value) || 0,
        }

        const url    = id ? `${API}/classes/${id}` : `${API}/classes`
        const method = id ? 'PUT' : 'POST'

        const res = await authFetch(url, {
            method,
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(data)
        })

        if (res.ok) {
            showToast(id ? 'Class updated!' : 'Class created!')
            closeModal()
            loadClasses()
        } else {
            const err = await res.json()
            showToast('Error: ' + (err.error || 'Save failed'), 'error')
        }
    } finally {
        submitBtn.disabled = false
        submitBtn.textContent = origText
    }
})

// uppercase section & numeric-only filters
document.getElementById('section').addEventListener('input', e => {
    e.target.value = e.target.value.toUpperCase().replace(/[^A-Z]/g, '')
})
document.getElementById('teacherContact').addEventListener('input', e => {
    e.target.value = e.target.value.replace(/\D/g, '')
})

async function deleteClass(id) {
    const ok = await confirmDialog({
        title: 'Delete this class?',
        message: 'Students in this class will not be deleted, but their class assignment may break.',
        confirmText: 'Delete',
        danger: true,
    })
    if (!ok) return
    const res = await authFetch(`${API}/classes/${id}`, { method: 'DELETE' })
    if (res.ok) {
        showToast('Class deleted!')
        loadClasses()
    } else {
        showToast('Failed to delete', 'error')
    }
}

loadClasses()