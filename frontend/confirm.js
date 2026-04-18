// ── Custom Confirm Dialog ─────────────────────────────────────
// Returns a Promise<boolean>
// Usage: const ok = await confirmDialog({ title, message, confirmText, danger })

function confirmDialog({ title = 'Are you sure?', message = '', confirmText = 'Confirm', cancelText = 'Cancel', danger = false } = {}) {
    return new Promise(resolve => {
        // remove any existing dialog
        document.querySelectorAll('.confirm-overlay').forEach(el => el.remove())

        const overlay = document.createElement('div')
        overlay.className = 'confirm-overlay'
        overlay.innerHTML = `
            <div class="confirm-box ${danger ? 'danger' : ''}">
                <div class="confirm-icon">${danger ? '⚠' : '?'}</div>
                <div class="confirm-title">${escapeHtml(title)}</div>
                <div class="confirm-message">${escapeHtml(message)}</div>
                <div class="confirm-actions">
                    <button class="confirm-btn-cancel">${escapeHtml(cancelText)}</button>
                    <button class="confirm-btn-ok ${danger ? 'danger' : ''}">${escapeHtml(confirmText)}</button>
                </div>
            </div>
        `
        document.body.appendChild(overlay)

        const cleanup = (result) => {
            overlay.classList.add('closing')
            setTimeout(() => overlay.remove(), 150)
            resolve(result)
        }

        overlay.querySelector('.confirm-btn-cancel').onclick = () => cleanup(false)
        overlay.querySelector('.confirm-btn-ok').onclick = () => cleanup(true)
        overlay.onclick = (e) => { if (e.target === overlay) cleanup(false) }

        // Esc to cancel, Enter to confirm
        const onKey = (e) => {
            if (e.key === 'Escape') { cleanup(false); document.removeEventListener('keydown', onKey) }
            if (e.key === 'Enter')  { cleanup(true);  document.removeEventListener('keydown', onKey) }
        }
        document.addEventListener('keydown', onKey)

        // animate in
        requestAnimationFrame(() => overlay.classList.add('show'))

        // focus the confirm button so Enter works naturally
        setTimeout(() => overlay.querySelector('.confirm-btn-ok')?.focus(), 50)
    })
}
