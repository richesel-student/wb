const BASE = "http://localhost:8080"

// CREATE
async function create() {
  const data = {
    asset_name: document.getElementById("asset").value,
    type: document.getElementById("type").value,
    amount_usd: parseFloat(document.getElementById("amount").value)
  }
   console.log("chart data:", data)

  await fetch(BASE + "/items", {
    method: "POST",
    headers: {"Content-Type":"application/json"},
    body: JSON.stringify(data)
  })

  loadItems()
}

// GET ITEMS
async function loadItems() {
  const from = document.getElementById("from").value
  const to = document.getElementById("to").value

  let url = BASE + "/items"
  if (from && to) url += `?from=${from}&to=${to}`

  const res = await fetch(url)
  const items = await res.json()

  const table = document.getElementById("table")
  table.innerHTML = ""

  items.forEach(i => {
    table.innerHTML += `
      <tr>
        <td>${i.id}</td>
        <td>${i.asset_name}</td>
        <td>${i.type}</td>
        <td>${i.amount_usd}</td>
        <td>${i.price.toFixed(2)}</td>
        <td>${i.quantity.toFixed(6)}</td>
        <td>${i.fee}</td>
        <td>${new Date(i.timestamp).toLocaleString()}</td>
      </tr>
    `
  })
}

// ANALYTICS
async function loadAnalytics() {
  const from = document.getElementById("from").value
  const to = document.getElementById("to").value

  let url = BASE + "/analytics"
  if (from && to) url += `?from=${from}&to=${to}`

  const res = await fetch(url)
  const data = await res.json()

  document.getElementById("analytics").innerHTML = `
    Sum: ${data.sum}<br>
    Avg: ${data.avg}<br>
    Count: ${data.count}<br>
    Median: ${data.median}<br>
    P90: ${data.p90}
  `
}

// CSV
function exportCSV() {
  window.location = BASE + "/export"
}

// ✅ ГРАФИК
let chart

async function loadChart() {
  const res = await fetch("http://localhost:8080/chart")
  const data = await res.json()

  console.log(data)
/*
  const labels = data.map(x => x.date)
  const values = data.map(x => x.total)

    const labels = ["A", "B", "C"]
    const values = [1000, 2000, 1500]
*/
    const labels = data.map(x => 
  new Date(x.date).toLocaleDateString()
)

    const values = data.map(x => x.total)
  

  const ctx = document.getElementById("chart").getContext("2d")

  if (chart) chart.destroy()

    chart = new Chart(ctx, {
    type: "line",
    data: {
        labels: labels,
        datasets: [{
        label: "USD per day",
        data: values,
        borderWidth: 2,
        pointRadius: 6,              
        pointBackgroundColor: "red"  
        }]
    }
    })
}


// INIT
loadItems()