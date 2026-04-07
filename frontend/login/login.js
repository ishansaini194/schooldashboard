const roles = document.querySelectorAll(".role")
let selectedRole = "admin"

// Role selection
roles.forEach(btn => {
    btn.addEventListener("click", () => {
        roles.forEach(b => b.classList.remove("active"))
        btn.classList.add("active")
        selectedRole = btn.dataset.role
    })
})

// Form submit
document.getElementById("loginForm").addEventListener("submit", async function (e) {
    e.preventDefault()

    const formData = new FormData(this)
    const data = {
        username: formData.get("username"),
        password: formData.get("password"),
        role: selectedRole
    }

    console.log("Login Data:", data)

    // Example API call (you'll connect later)
    /*
    const res = await fetch("http://localhost:8080/api/login", {
        method: "POST",
        headers: {"Content-Type": "application/json"},
        body: JSON.stringify(data)
    })
    */
})