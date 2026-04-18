requireAuth()

const API = CONFIG.API
const form = document.getElementById('studentForm')

// pre-fill class + section if passed via URL
const params = new URLSearchParams(window.location.search)
const classId = params.get('class')
const sectionId = params.get('section')
if (classId) {
    document.querySelector('input[name="class"]').value = classId
}
if (sectionId) {
    document.querySelector('input[name="section"]').value = sectionId
}

// uppercase the section input as user types
const sectionInput = document.querySelector('input[name="section"]')
if (sectionInput) {
    sectionInput.addEventListener('input', e => {
        e.target.value = e.target.value.toUpperCase().replace(/[^A-Z]/g, '')
    })
}

// strip non-digit chars from numeric fields as user types
const numericFields = ['phone', 'aadhar_no', 'father_contact', 'father_aadhar', 'mother_contact', 'roll_no']
numericFields.forEach(name => {
    const el = form.querySelector(`[name="${name}"]`)
    if (el) {
        el.addEventListener('input', e => {
            e.target.value = e.target.value.replace(/\D/g, '')
        })
    }
})

// real-time validation on blur
attachLiveValidation(form)

function showToast(msg, type = 'success') {
    const toast = document.getElementById('toast')
    toast.textContent = msg
    toast.className = `toast ${type} show`
    setTimeout(() => toast.classList.remove('show'), 3000)
}

form.addEventListener('submit', async function (e) {
    e.preventDefault()

    const formData = new FormData(this)
    const data = {}
    formData.forEach((value, key) => { data[key] = value })

    // validate
    const errors = validateStudent(data)
    if (Object.keys(errors).length > 0) {
        showFieldErrors(errors, this)
        const n = Object.keys(errors).length
        showToast(`Please fix ${n} error${n > 1 ? 's' : ''} below`, 'error')
        return
    }

    data.class_id = parseInt(data.class) || 0

    const submitBtn = this.querySelector('button[type="submit"]')
    submitBtn.disabled = true
    const origText = submitBtn.textContent
    submitBtn.textContent = 'Saving...'

    try {
        const response = await authFetch(`${API}/students`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(data)
        })

        const result = await response.json()

        if (response.ok) {
            showToast('Student created successfully!')
            this.reset()
            // clear any lingering error states
            this.querySelectorAll('.field-error').forEach(el => el.remove())
            this.querySelectorAll('.has-error').forEach(el => el.classList.remove('has-error'))
        } else {
            // backend may return field-level errors
            if (result.fields && typeof result.fields === 'object') {
                showFieldErrors(result.fields, this)
                showToast('Server rejected the form. Please fix the errors below.', 'error')
            } else {
                showToast('Error: ' + (result.error || 'Unknown error'), 'error')
            }
        }
    } finally {
        submitBtn.disabled = false
        submitBtn.textContent = origText
    }
})
