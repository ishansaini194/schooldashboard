const API = 'http://localhost:8080/api'

// get class_id from URL
const params = new URLSearchParams(window.location.search)
const classId = params.get('class')

function showToast(msg, type = 'success') {
    const toast = document.getElementById('toast')
    toast.textContent = msg
    toast.className = `toast ${type} show`
    setTimeout(() => toast.classList.remove('show'), 3000)
}

async function loadStudents() {
    const res = await fetch(`${API}/students/class/${classId}`)
    const students = await res.json()
    const list = document.getElementById('studentList')

    document.getElementById('pageSubtitle').textContent = `${students.length || 0} students in this class`

    if (!students || students.length === 0) {
        list.innerHTML = '<div class="empty-state">No students in this class yet.</div>'
        return
    }

    list.innerHTML = students.map(s => `
    <div class="student-row">
        <div class="student-roll">${s.roll_no || '—'}</div>
        <div class="student-info">
            <div class="student-name">${s.name}</div>
            <div class="student-meta">${s.gender || ''} ${s.dob ? '· DOB: ' + s.dob : ''}</div>
        </div>
        <a class="student-phone" href="tel:${s.phone}">${s.phone || '—'}</a>
        <button class="btn-edit" onclick="window.location.href='student-detail.html?roll_no=${s.roll_no}'">View</button>
    </div>
`).join('')
}

loadStudents()