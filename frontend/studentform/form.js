document.getElementById("studentForm").addEventListener("submit", async function (e) {
    e.preventDefault()

    const formData = new FormData(this)
    const data = {}

    formData.forEach((value, key) => {
        data[key] = value
    })

    // convert class_id to number
    data.class_id = parseInt(data.class_id)

    const response = await fetch("http://localhost:8080/api/students", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(data)
    })

    const result = await response.json()

    if (response.ok) {
        alert("Student created successfully!")
        this.reset()
    } else {
        alert("Error: " + result.error)
    }
})