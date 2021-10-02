import { Chart, registerables } from 'chart.js';
import 'chartjs-adapter-moment';
Chart.register(...registerables);

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

const dataSets = [
    {
        label: 'Anaheim',
        data: [],
        fill: false,
        borderColor: 'rgb(252, 76, 2)',
    },
    {
        label: 'Arizona',
        data: [],
        fill: false,
        borderColor: 'rgb(140, 38, 51)',
    },
    {
        label: 'Boston',
        data: [],
        fill: false,
        borderColor: 'rgb(252, 181, 20)',
    },
    {
        label: 'Buffalo',
        data: [],
        fill: false,
        borderColor: 'rgb(0, 38, 84)',
    },
    {
        label: 'Calgary',
        data: [],
        fill: false,
        borderColor: 'rgb(200, 16, 46)',
    },
    {
        label: 'Carolina',
        data: [],
        fill: false,
        borderColor: 'rgb(226,24,54)',
    },
    {
        label: 'Chicago',
        data: [],
        fill: false,
        borderColor: 'rgb(207,10,44)',
    },
    {
        label: 'Colorado',
        data: [],
        fill: false,
        borderColor: 'rgb(111, 38, 61)',
    },
    {
        label: 'Columbus',
        data: [],
        fill: false,
        borderColor: 'rgb(0,38,84)',
    },
    {
        label: 'Dallas',
        data: [],
        fill: false,
        borderColor: 'rgb(0, 104, 71)',
    },
    {
        label: 'Detroit',
        data: [],
        fill: false,
        borderColor: 'rgb(206,17,38)',
    },
    {
        label: 'Edmonton',
        data: [],
        fill: false,
        borderColor: 'rgb(206,17,38)',
    },
    {
        label: 'Florida',
        data: [],
        fill: false,
        borderColor: 'rgb(4,30,66)',
    },
    {
        label: 'Los Angeles',
        data: [],
        fill: false,
        borderColor: 'rgb(17,17,17)',
    },
    {
        label: 'Minnesota',
        data: [],
        fill: false,
        borderColor: 'rgb(2, 73, 48)',
    },
    {
        label: 'Montreal',
        data: [],
        fill: false,
        borderColor: 'rgb(175, 30, 45)',
    },
    {
        label: 'Nashville',
        data: [],
        fill: false,
        borderColor: 'rgb(255,184,28)',
    },
    {
        label: 'New Jersey',
        data: [],
        fill: false,
        borderColor: 'rgb(206, 17, 38)',
    },
    {
        label: 'NY Islanders',
        data: [],
        fill: false,
        borderColor: 'rgb(0,83,155)',
    },
    {
        label: 'NY Rangers',
        data: [],
        fill: false,
        borderColor: 'rgb(0,56,168)',
    },
    {
        label: 'Ottawa',
        data: [],
        fill: false,
        borderColor: 'rgb(197 32 50)',
    },
    {
        label: 'Philadelphia',
        data: [],
        fill: false,
        borderColor: 'rgb(247, 73, 2)',
    },
    {
        label: 'Pittsburgh',
        data: [],
        fill: false,
        borderColor: 'rgb(252,181,20)',
    },
    {
        label: 'St Louis',
        data: [],
        fill: false,
        borderColor: 'rgb(0, 47, 135)',
    },
    {
        label: 'San Jose',
        data: [],
        fill: false,
        borderColor: 'rgb(0, 109, 117)',
    },
    {
        label: 'Seattle',
        data: [],
        fill: false,
        borderColor: 'rgb(153, 217, 217)',
    },
    {
        label: 'Tampa Bay',
        data: [],
        fill: false,
        borderColor: 'rgb(0, 40, 104)',
    },
    {
        label: 'Toronto',
        data: [],
        fill: false,
        borderColor: 'rgb(0, 32, 91)',
    },
    {
        label: 'Vancouver',
        data: [],
        fill: false,
        borderColor: 'rgb(0, 32, 91)',
    },
    {
        label: 'Vegas',
        data: [],
        fill: false,
        borderColor: 'rgb(185,151,91)',
    },
    {
        label: 'Washington',
        data: [],
        fill: false,
        borderColor: 'rgb(4 30 66)',
    },
    {
        label: 'Winnipeg',
        data: [],
        fill: false,
        borderColor: 'rgb(0,76,151)',
    }
]

// NHLStandingsChart - object which represents the NHL standings Chart
let NHLStandingsChart

function getSeasonData(season) {
    fetch(`http://localhost:8081/season/${season}`)
    .then(r => r.json())
    .then(data => {
        createDataSet(data, parseInt(season))
    })
}

function createDataSet(seasonData, season) {
    const currentTeamPoints = {}
    for (const dailyRecord of seasonData.dailyRecords) {
        for (const teamRecord of dailyRecord.teamRecords) {
            // skip dates where the points don't change day to day. Helps declutter the chart
            if (currentTeamPoints[teamRecord.teamName] === teamRecord.points) {
                continue
            }
            currentTeamPoints[teamRecord.teamName] = teamRecord.points
            const teamIndex = dataSets.findIndex(team => {
                return team.label === teamRecord.teamName
            })
            let date
            if (dailyRecord.month === "Oct" || dailyRecord.month === "Nov" || dailyRecord.month === "Dec") {
                date = new Date(season-1, monthIntegerLookup[dailyRecord.month], dailyRecord.day)
            } else {
                date = new Date(season, monthIntegerLookup[dailyRecord.month], dailyRecord.day)
            }
            const dataPoint = {
                'x': date,
                'y': teamRecord.points
            }
            dataSets[teamIndex].data.push(dataPoint)
        }
    }

    const ctx = document.getElementById('myChart');
    NHLStandingsChart =  new Chart(ctx, {
        type: 'line',
        data: {
            datasets: dataSets
        },
        options: {
            scales: {
                x: {
                    type: 'time',
                    time: {
                        displayFormats: {
                            day: 'DD/MM/YYYY'
                        }
                    }
                },
                y: {
                    beginAtZero: true
                }
            },
            plugins: {
                legend: {
                    display: false
                }
            }
        }
    });
}

function addTeamLogo() {
    for (const dataSet of dataSets) {
        const image = new Image()
        image.src = `images/${dataSet.label}.png`
        dataSet.pointStyle = image

    }
}

function addTeamsSelect() {
    const teamSelect = document.getElementById("teamSelection")
    for (const dataSet of dataSets) {
        const option = document.createElement("option")
        option.value = dataSet.label
        option.text = dataSet.label
        teamSelect.appendChild(option)
    }
}

function addSelectionChangeHandler() {
    const teamSelect = document.getElementById("teamSelection")
    teamSelect.addEventListener('change', () => {
        for (const dataset of NHLStandingsChart.data.datasets) {
            if (teamSelect.value === 'reset') {
                dataset.hidden = false
            } else if (dataset.label !== teamSelect.value) {
                dataset.hidden = true
            } else {
                // needed to undo previous selection
                dataset.hidden = false
            }
        }
        NHLStandingsChart.update()
    }) 
}

addTeamLogo()
addTeamsSelect()
addSelectionChangeHandler()
getSeasonData(2020)