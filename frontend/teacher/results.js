requireAuth()
const API = CONFIG.API
const username = localStorage.getItem('username')
let students = []

async function init() {
    await loadClasses()
}

async function loadClasses() {
    const res = await authFetch(`${API}/classes`)
    if (!res || !res.ok) return
    const classes = await res.json() || []
    document.getElementById('resClass').innerHTML =
        '<option value="">Select class</option>' +
        classes.map(c => `<option value="${c.class}">${c.class}</option>`).join('')
}

async function loadStudents() {
    const cls = document.getElementById('resClass').value
    const section = document.getElementById('resSection').value
    if (!cls) return

    const subject = document.getElementById('resSubject').value
    const examType = document.getElementById('resExamType').value
    const year = document.getElementById('resYear').value

    const res = await authFetch(`${API}/students/class/${cls}`)
    if (!res || !res.ok) return
    let allStudents = await res.json() || []

    if (section) allStudents = allStudents.filter(s => s.section === section)
    allStudents = [...new Map(allStudents.map(s => [s.ID || s.id, s])).values()]
    students = allStudents

    document.getElementById('studentMarksSection').style.display = 'block'
    document.getElementById('marksList').innerHTML = students.map((s, index) => `
        <div style="display:flex; align-items:center; gap:16px; background:var(--surface); border:1px solid var(--border); border-radius:10px; padding:12px 16px; margin-bottom:8px;">
            <div style="flex:1;">
                <div style="font-size:14px; font-weight:500; color:var(--text);">${s.name}</div>
                <div style="font-size:12px; color:var(--muted);">Roll ${s.roll_no || '—'}</div>
            </div>
            <input
                type="number"
                class="marks-input"
                id="marks_${s.ID || s.id}"
                data-index="${index}"
                placeholder="Marks"
                min="0"
                onkeydown="handleMarksKey(event, ${index})"
            >
        </div>
    `).join('')

    // load existing marks if subject and exam type selected
    if (subject && examType && year) {
        await loadExistingMarks(cls, section, subject, examType, year)
    }
}

async function loadExistingMarks(cls, section, subject, examType, year) {
    const res = await authFetch(
        `${API}/results/class/${cls}/section/${section || ''}?subject=${subject}&exam_type=${examType}&year=${year}`
    )
    if (!res || !res.ok) return
    const results = await res.json() || []

    // pre-fill inputs with existing marks
    results.forEach(r => {
        const input = document.getElementById(`marks_${r.student_id}`)
        if (input) input.value = r.marks
    })

    if (results.length > 0) {
        showToast(`${results.length} existing marks loaded`, 'success')
    }
}

async function saveAllMarks() {
    const cls = document.getElementById('resClass').value
    const section = document.getElementById('resSection').value
    const subject = document.getElementById('resSubject').value
    const examType = document.getElementById('resExamType').value
    const maxMarks = parseInt(document.getElementById('resMaxMarks').value)
    const year = parseInt(document.getElementById('resYear').value)

    if (!cls || !subject) {
        showToast('Please select class and subject', 'error')
        return
    }

    let saved = 0
    let skipped = 0

    for (const s of students) {
        const id = s.ID || s.id
        const marksInput = document.getElementById(`marks_${id}`)
        const marksVal = marksInput?.value.trim()
        if (!marksVal) { skipped++; continue }

        const res = await authFetch(`${API}/results`, {
            method: 'POST',
            body: JSON.stringify({
                student_id: id,
                subject,
                exam_type: examType,
                marks: parseInt(marksVal),
                max_marks: maxMarks,
                year,
                class: cls,
                section: section || '',
                entered_by: username || 'Teacher'
            })
        })
        if (res && res.ok) saved++
        else skipped++
    }

    showToast(`Saved ${saved} results${skipped > 0 ? ', ' + skipped + ' skipped' : ''}`)
}

function showToast(msg, type = 'success') {
    const toast = document.getElementById('toast')
    toast.textContent = msg
    toast.className = `toast ${type} show`
    setTimeout(() => toast.classList.remove('show'), 3000)
}

function handleMarksKey(event, index) {
    if (event.key === 'Enter') {
        event.preventDefault()
        const next = document.querySelector(`input[data-index="${index + 1}"]`)
        if (next) {
            next.focus()
        } else {
            saveAllMarks()
        }
    }
}

init()