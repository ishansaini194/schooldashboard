requireAuth()

const API = CONFIG.API

const params = new URLSearchParams(window.location.search)
const rollNo = params.get('roll_no')
const classId = params.get('class') || ''
const sectionId = params.get('section') || ''

// query params for student lookup — disambiguates roll_no across classes
const lookupQS = (() => {
    const qs = new URLSearchParams()
    if (classId) qs.set('class', classId)
    if (sectionId) qs.set('section', sectionId)
    const s = qs.toString()
    return s ? `?${s}` : ''
})()

let currentStudent = null

function showToast(msg, type = 'success') {
    const toast = document.getElementById('toast')
    toast.textContent = msg
    toast.className = `toast ${type} show`
    setTimeout(() => toast.classList.remove('show'), 3000)
}

async function loadStudent() {
    const res = await authFetch(`${API}/students/${rollNo}${lookupQS}`)
    if (!res.ok) {
        document.getElementById('studentName').textContent = 'Student not found'
        return
    }

    const s = await res.json()
    currentStudent = s

    // Hero
    document.getElementById('avatarInitial').textContent = s.name ? s.name[0].toUpperCase() : '?'
    document.getElementById('studentName').textContent = s.name || '—'
    document.getElementById('studentBadges').innerHTML = `
        ${s.class ? `<span class="badge">Class ${s.class}</span>` : ''}
        ${s.section ? `<span class="badge">${s.section}</span>` : ''}
        ${s.roll_no ? `<span class="badge neutral">Roll ${s.roll_no}</span>` : ''}
        ${s.gender ? `<span class="badge neutral">${s.gender}</span>` : ''}
    `

    // Basic Info
    document.getElementById('d-phone').innerHTML   = s.phone ? `<a href="tel:${s.phone}">${s.phone}</a>` : '—'
    document.getElementById('d-gender').textContent    = s.gender || '—'
    document.getElementById('d-dob').textContent       = s.dob || '—'
    document.getElementById('d-address').textContent   = s.address || '—'
    document.getElementById('d-caste').textContent     = s.caste || '—'

    // ID Details
    document.getElementById('d-aadhar').textContent    = s.aadhar_no || '—'
    document.getElementById('d-epunjab').textContent   = s.epunjab_id || '—'
    document.getElementById('d-prevschool').textContent= s.previous_school_details || '—'

    // Father
    document.getElementById('d-fname').textContent     = s.father_name || '—'
    document.getElementById('d-fcontact').textContent  = s.father_contact || '—'
    document.getElementById('d-fcontact').href         = s.father_contact ? `tel:${s.father_contact}` : '#'
    document.getElementById('d-faadhar').textContent   = s.father_aadhar || '—'

    // Mother
    document.getElementById('d-mname').textContent     = s.mother_name || '—'
    document.getElementById('d-mcontact').textContent  = s.mother_contact || '—'
    document.getElementById('d-mcontact').href         = s.mother_contact ? `tel:${s.mother_contact}` : '#'
}

function openEditModal() {
    if (!currentStudent) return
    const s = currentStudent

    document.getElementById('e-name').value      = s.name || ''
    document.getElementById('e-class').value     = s.class || ''
    document.getElementById('e-section').value   = s.section || ''
    document.getElementById('e-rollno').value    = s.roll_no || ''
    document.getElementById('e-gender').value    = s.gender || ''
    document.getElementById('e-dob').value       = s.dob || ''
    document.getElementById('e-phone').value     = s.phone || ''
    document.getElementById('e-aadhar').value    = s.aadhar_no || ''
    document.getElementById('e-epunjab').value   = s.epunjab_id || ''
    document.getElementById('e-fname').value     = s.father_name || ''
    document.getElementById('e-fcontact').value  = s.father_contact || ''
    document.getElementById('e-faadhar').value   = s.father_aadhar || ''
    document.getElementById('e-mname').value     = s.mother_name || ''
    document.getElementById('e-mcontact').value  = s.mother_contact || ''
    document.getElementById('e-caste').value     = s.caste || ''
    document.getElementById('e-address').value   = s.address || ''
    document.getElementById('e-prevschool').value= s.previous_school_details || ''

    document.getElementById('editModal').classList.add('show')
}

function closeEditModal() {
    document.getElementById('editModal').classList.remove('show')
}

document.getElementById('editForm').addEventListener('submit', async function(e) {
    e.preventDefault()

    const data = {
        name:                    document.getElementById('e-name').value,
        class:                   document.getElementById('e-class').value,
        section:                 document.getElementById('e-section').value,
        roll_no:                 document.getElementById('e-rollno').value,
        gender:                  document.getElementById('e-gender').value,
        dob:                     document.getElementById('e-dob').value,
        phone:                   document.getElementById('e-phone').value,
        aadhar_no:               document.getElementById('e-aadhar').value,
        epunjab_id:              document.getElementById('e-epunjab').value,
        father_name:             document.getElementById('e-fname').value,
        father_contact:          document.getElementById('e-fcontact').value,
        father_aadhar:           document.getElementById('e-faadhar').value,
        mother_name:             document.getElementById('e-mname').value,
        mother_contact:          document.getElementById('e-mcontact').value,
        caste:                   document.getElementById('e-caste').value,
        address:                 document.getElementById('e-address').value,
        previous_school_details: document.getElementById('e-prevschool').value,
    }

    // validate
    const errors = validateStudent(data)
    if (Object.keys(errors).length > 0) {
        showFieldErrors(errors, this)
        const n = Object.keys(errors).length
        showToast(`Please fix ${n} error${n > 1 ? 's' : ''} below`, 'error')
        return
    }

    const submitBtn = this.querySelector('button[type="submit"]')
    submitBtn.disabled = true
    const origText = submitBtn.textContent
    submitBtn.textContent = 'Saving...'

    try {
        const res = await authFetch(`${API}/students/${rollNo}${lookupQS}`, {
            method: 'PUT',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(data)
        })

        if (res.ok) {
            showToast('Student updated!')
            closeEditModal()
            loadStudent()
        } else {
            const err = await res.json()
            if (err.fields && typeof err.fields === 'object') {
                showFieldErrors(err.fields, this)
                showToast('Server rejected the form. Please fix the errors below.', 'error')
            } else {
                showToast('Error: ' + (err.error || 'Unknown error'), 'error')
            }
        }
    } finally {
        submitBtn.disabled = false
        submitBtn.textContent = origText
    }
})

// attach live validation + input sanitization to edit form (once)
;(function setupEditFormValidation() {
    const editForm = document.getElementById('editForm')
    if (!editForm) return

    // uppercase section
    const sec = document.getElementById('e-section')
    if (sec) sec.addEventListener('input', e => {
        e.target.value = e.target.value.toUpperCase().replace(/[^A-Z]/g, '')
    })

    // numeric-only fields
    const numericIds = ['e-phone', 'e-aadhar', 'e-fcontact', 'e-faadhar', 'e-mcontact', 'e-rollno']
    numericIds.forEach(id => {
        const el = document.getElementById(id)
        if (el) el.addEventListener('input', e => {
            e.target.value = e.target.value.replace(/\D/g, '')
        })
    })

    attachLiveValidation(editForm)
})()

async function deleteStudent() {
    const ok = await confirmDialog({
        title: 'Delete student?',
        message: `${currentStudent?.name || 'This student'} will be permanently removed. Their fee history will also be lost.`,
        confirmText: 'Delete',
        danger: true,
    })
    if (!ok) return

    const res = await authFetch(`${API}/students/${rollNo}${lookupQS}`, { method: 'DELETE' })
    if (res.ok) {
        showToast('Student deleted!')
        setTimeout(() => history.back(), 1000)
    } else {
        showToast('Failed to delete', 'error')
    }
}

loadStudent()