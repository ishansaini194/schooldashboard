requireAuth()

const API = CONFIG.API
const studentID = localStorage.getItem('student_id')
const year = new Date().getFullYear()
let currentTab = 'midterm'

async function init() {
    document.getElementById('yearLabel').textContent = `Academic year ${year}`
    await loadResults('midterm')
}

async function switchTab(tab, btn) {
    document.querySelectorAll('.tab-btn').forEach(b => b.classList.remove('active'))
    btn.classList.add('active')
    currentTab = tab
    await loadResults(tab)
}

async function loadResults(examType) {
    document.getElementById('resultList').innerHTML = '<div class="dash-loading">Loading...</div>'

    const res = await authFetch(`${API}/results/student/${studentID}?exam_type=${examType}&year=${year}`)
    if (!res || !res.ok) return
    const list = await res.json() || []

    if (list.length === 0) {
        document.getElementById('resultList').innerHTML =
            `<div class="dash-empty">No ${examType === 'midterm' ? 'mid-term' : 'final'} results yet</div>`
        return
    }

    const total = list.reduce((sum, r) => sum + r.marks, 0)
    const maxTotal = list.reduce((sum, r) => sum + r.max_marks, 0)
    const pct = maxTotal > 0 ? Math.round((total / maxTotal) * 100) : 0

    document.getElementById('resultList').innerHTML = `
        <div class="result-summary-card">
            <div class="result-summary-label">Total Score</div>
            <div class="result-summary-score">${total} / ${maxTotal}</div>
            <div class="result-summary-pct">${pct}%</div>
            <div class="result-bar-wrap">
                <div class="result-bar" style="width:${pct}%; background:${pct >= 75 ? 'var(--accent)' : pct >= 50 ? '#c07020' : '#c0392b'}"></div>
            </div>
        </div>
        ${list.map(r => {
        const p = Math.round((r.marks / r.max_marks) * 100)
        const color = p >= 75 ? 'var(--accent)' : p >= 50 ? '#c07020' : '#c0392b'
        return `
            <div class="result-row">
                <div class="result-subject">${esc(r.subject)}</div>
                <div class="result-bar-wrap" style="flex:1; margin: 0 16px;">
                    <div class="result-bar" style="width:${p}%; background:${color}"></div>
                </div>
                <div class="result-score" style="color:${color}">${r.marks}<span>/${r.max_marks}</span></div>
            </div>`
    }).join('')}
    `
}

init()