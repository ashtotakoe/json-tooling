const input = document.getElementById("input")
const output = document.getElementById("output")

const parseButton = document.getElementById("parse")

const parse = () => {
  fetch("http://localhost:3000/parse", {
    method: "POST",
    body: input.value
  })
    .then(
      (res) => {
        return res.text()
      }

    )
    .then(data => {
      output.textContent = data
    })
    .catch(
      (e) => {
        output.value = e.message
      }
    )

}

parseButton.addEventListener('click', parse)
