const API = 'http://localhost:8080/api'

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
        document.getElementById('classId').value          = cls.ID
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
    const res     = await fetch(`${API}/classes`)
    const classes = await res.json()
    const grid    = document.getElementById('classGrid')

    if (!classes || classes.length === 0) {
        grid.innerHTML = '<div class="empty-state">No classes yet. Add one!</div>'
        return
    }

    grid.innerHTML = classes.map(cls => `
        <div class="class-card" onclick="window.location.href='students.html?class=${cls.class}'" style="cursor:pointer">
            <div class="class-badge">Class ${cls.class} — ${cls.section || '—'}</div>
            <div class="class-teacher">
                <span class="teacher-icon">👤</span>
                <div>
                    <div class="teacher-name">${cls.teacher_name || 'No teacher assigned'}</div>
                    <div class="teacher-contact">${cls.teacher_contact || ''}</div>
                </div>
            </div>
            <div class="card-actions">
                <button class="btn-edit" onclick='event.stopPropagation(); openModal(${JSON.stringify(cls)})'>Edit</button>
                <button class="btn-delete" onclick='event.stopPropagation(); deleteClass(${cls.ID})'>Delete</button>
            </div>
        </div>
    `).join('')
}

document.getElementById('classForm').addEventListener('submit', async function(e) {
    e.preventDefault()

    const id   = document.getElementById('classId').value
    const data = {
        class:          parseInt(document.getElementById('classNo').value),
        section:        document.getElementById('section').value,
        teacher_name:   document.getElementById('teacherName').value,
        teacher_contact:document.getElementById('teacherContact').value,
        tuition_fee:    parseInt(document.getElementById('tuitionFee').value) || 0,
        transport_fee:  parseInt(document.getElementById('transportFee').value) || 0,
    }

    const url    = id ? `${API}/classes/${id}` : `${API}/classes`
    const method = id ? 'PUT' : 'POST'

    const res = await fetch(url, {
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
        showToast('Error: ' + err.error, 'error')
    }
})

async function deleteClass(id) {
    if (!confirm('Delete this class?')) return
    const res = await fetch(`${API}/classes/${id}`, { method: 'DELETE' })
    if (res.ok) {
        showToast('Class deleted!')
        loadClasses()
    } else {
        showToast('Failed to delete', 'error')
    }
}

loadClasses()