// ── Shared Form Validators ───────────────────────────────────
// Usage:
//   const errs = validateStudent(data)
//   if (Object.keys(errs).length > 0) { showFieldErrors(errs); return }

const VALIDATORS = {
    name: {
        required: true,
        test: v => /^[A-Za-z][A-Za-z\s.'-]{1,49}$/.test(v.trim()),
        message: 'Only letters, spaces, dots allowed (2-50 chars)',
    },
    class: {
        required: true,
        test: v => /^([1-9]|1[0-2])$/.test(String(v).trim()),
        message: 'Must be a number between 1 and 12',
    },
    section: {
        required: true,
        test: v => /^[A-Za-z]$/.test(v.trim()),
        message: 'Must be a single letter (A-Z)',
    },
    roll_no: {
        required: true,
        test: v => /^\d{1,4}$/.test(v.trim()),
        message: 'Must be 1-4 digits',
    },
    epunjab_id: {
        required: true,
        test: v => /^[A-Za-z0-9]{4,20}$/.test(v.trim()),
        message: 'Must be 4-20 letters/digits',
    },
    phone: {
        required: true,
        test: v => /^\d{10}$/.test(v.trim()),
        message: 'Must be exactly 10 digits',
    },
    aadhar_no: {
        required: false,
        test: v => v.trim() === '' || /^\d{12}$/.test(v.trim()),
        message: 'Must be exactly 12 digits',
    },
    father_name: {
        required: false,
        test: v => v.trim() === '' || /^[A-Za-z][A-Za-z\s.'-]{1,49}$/.test(v.trim()),
        message: 'Only letters, spaces, dots allowed',
    },
    father_contact: {
        required: false,
        test: v => v.trim() === '' || /^\d{10}$/.test(v.trim()),
        message: 'Must be exactly 10 digits',
    },
    father_aadhar: {
        required: false,
        test: v => v.trim() === '' || /^\d{12}$/.test(v.trim()),
        message: 'Must be exactly 12 digits',
    },
    mother_name: {
        required: false,
        test: v => v.trim() === '' || /^[A-Za-z][A-Za-z\s.'-]{1,49}$/.test(v.trim()),
        message: 'Only letters, spaces, dots allowed',
    },
    mother_contact: {
        required: false,
        test: v => v.trim() === '' || /^\d{10}$/.test(v.trim()),
        message: 'Must be exactly 10 digits',
    },
    caste: {
        required: false,
        test: v => v.trim() === '' || /^[A-Za-z\s]{2,30}$/.test(v.trim()),
        message: 'Only letters and spaces (2-30 chars)',
    },
    gender: {
        required: false,
        test: v => v === '' || ['male', 'female'].includes(v),
        message: 'Select male or female',
    },
    dob: {
        required: false,
        test: v => {
            if (!v) return true
            const d = new Date(v)
            if (isNaN(d)) return false
            const age = (Date.now() - d.getTime()) / (365.25 * 24 * 3600 * 1000)
            return age >= 3 && age <= 25
        },
        message: 'Age must be between 3 and 25 years',
    },
    address: {
        required: false,
        test: v => v.trim() === '' || v.trim().length >= 5,
        message: 'Min 5 characters',
    },
}

// returns map of { fieldName: errorMessage }
function validateStudent(data) {
    const errors = {}
    for (const [field, rule] of Object.entries(VALIDATORS)) {
        const value = (data[field] ?? '').toString()
        if (rule.required && value.trim() === '') {
            errors[field] = 'Required'
            continue
        }
        if (!rule.test(value)) {
            errors[field] = rule.message
        }
    }
    return errors
}

// show errors inline under each field
function showFieldErrors(errors, formEl) {
    // clear previous errors
    formEl.querySelectorAll('.field-error').forEach(el => el.remove())
    formEl.querySelectorAll('.has-error').forEach(el => el.classList.remove('has-error'))

    for (const [field, msg] of Object.entries(errors)) {
        const input = formEl.querySelector(`[name="${field}"]`)
        if (!input) continue
        const wrap = input.closest('.field') || input.parentElement
        wrap.classList.add('has-error')
        const errEl = document.createElement('div')
        errEl.className = 'field-error'
        errEl.textContent = msg
        wrap.appendChild(errEl)
    }

    // scroll to first error
    const firstErr = formEl.querySelector('.has-error')
    if (firstErr) firstErr.scrollIntoView({ behavior: 'smooth', block: 'center' })
}

// real-time validation on blur
function attachLiveValidation(formEl) {
    formEl.querySelectorAll('input[name], select[name]').forEach(input => {
        input.addEventListener('blur', () => {
            const field = input.name
            const rule = VALIDATORS[field]
            if (!rule) return
            const value = input.value
            const wrap = input.closest('.field') || input.parentElement

            // remove existing error for this field
            wrap.classList.remove('has-error')
            wrap.querySelectorAll('.field-error').forEach(el => el.remove())

            if (rule.required && value.trim() === '') {
                showSingleError(wrap, 'Required')
            } else if (value.trim() !== '' && !rule.test(value)) {
                showSingleError(wrap, rule.message)
            }
        })
    })
}

function showSingleError(wrap, msg) {
    wrap.classList.add('has-error')
    const errEl = document.createElement('div')
    errEl.className = 'field-error'
    errEl.textContent = msg
    wrap.appendChild(errEl)
}
