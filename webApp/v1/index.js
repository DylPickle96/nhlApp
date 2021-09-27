

const monthNameLookup = {
    "Oct": "October",
    "Nov": "November",
    "Dec": "December",
    "Jan": "January",
    "Feb": "February",
    "Mar": "March",
    "Apr": "April"
}

const monthIntegerLookup = {
    "Oct": 9,
    "Nov": 10,
    "Dec": 11,
    "Jan": 0,
    "Feb": 1,
    "Mar": 2,
    "Apr": 3
}



function getSeasonData(season) {
    fetch(`http://localhost:8081/season/${season}`)
    .then(r => r.json())
    .then(data => {
        createDataSet(data, parseInt(season))
    })
}

function createDataSet(seasonData, season) {
    const dataPoints = []
    for (const dailyRecord of seasonData.dailyRecords) {
        for (const teamRecord of dailyRecord.teamRecords) {
            if (teamRecord.teamName !== "Vancouver") {
                continue
            }
            let date
            if (dailyRecord.month === "Oct" || dailyRecord.month === "Nov" || dailyRecord.month === "Dec") {
                date = new Date(season-1, monthIntegerLookup[dailyRecord.month], dailyRecord.day)
                console.log(date)
            } else {
                date = new Date(season, monthIntegerLookup[dailyRecord.month], dailyRecord.day)
            }
            const dataPoint = {
                label: 'Vancouver',
                'x': date,
                'y': teamRecord.points
            }
            dataPoints.push(dataPoint)
        }
    }

    const ctx = document.getElementById('myChart');
    console.log(dataPoints)
    new Chart(ctx, {
        type: 'line',
        data: {
            datasets: [{
                label: 'Vancouver',
                data: dataPoints,
                fill: false,
                borderColor: 'rgb(0, 32, 91)',
            }]
        },
        options: {
            scales: {
                x: {
                    type: 'time',
                    time: {
                        displayFormats: {
                            quarter: 'DD MM YYYY'
                        }
                    }
                },
                y: {
                    beginAtZero: true
                }
            }
        }
    });

}

getSeasonData(2020)