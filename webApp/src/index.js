import { Chart, registerables } from 'chart.js';
import 'chartjs-adapter-moment';
import 'bootstrap';
import 'bootstrap/dist/css/bootstrap.min.css';
Chart.register(...registerables);

const monthIntegerLookup = {
    "Oct": 9,
    "Nov": 10,
    "Dec": 11,
    "Jan": 0,
    "Feb": 1,
    "Mar": 2,
    "Apr": 3
}

const teamDataSets = [
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

const seasons = [2022, 2021, 2020, 2019, 2018, 2017, 2016, 2015, 2014, 2013, 2012, 2011, 2010, 2009, 2008, 2007, 2006, 2005, 2004, 2003, 2002, 2001]

// NHLStandingsChart - object which represents the NHL standings Chart
let NHLStandingsChart

// getSeasonData - Calls our API to fetch the season's data
function getSeasonData(season) {
    fetch(`http://localhost:8081/season/${season}`)
    .then(r => r.json())
    .then(data => {
        createDataSet(data, parseInt(season))
    })
}

// createDataSet - creates the dataset for each by iterating over all the daily records and matching them to the correct team
function createDataSet(seasonData, season) {
    resetTeamDataPoints()
    const currentTeamPoints = {}
    for (const dailyRecord of seasonData.dailyRecords) {
        for (const teamRecord of dailyRecord.teamRecords) {
            // skip dates where the points don't change day to day. Helps declutter the chart
            if (currentTeamPoints[teamRecord.teamName] === teamRecord.points) {
                continue
            }
            currentTeamPoints[teamRecord.teamName] = teamRecord.points
            const teamIndex = teamDataSets.findIndex(team => {
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
            teamDataSets[teamIndex].data.push(dataPoint)
        }
    }
    // renders chart
    const ctx = document.getElementById('myChart');
    NHLStandingsChart =  new Chart(ctx, {
        type: 'line',
        data: {
            datasets: teamDataSets
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

// addTeamLogo - Called to create an 'image' in the document for each team. This HTML element is then used as the label of each data point
function addTeamLogo() {
    for (const dataSet of teamDataSets) {
        const image = new Image()
        image.src = `images/${dataSet.label}.png`
        dataSet.pointStyle = image

    }
}

// addTeamsSelect - adds the teams from the declared data object to the selection element
function addTeamsSelect() {
    const teamSelect = document.getElementById("teamSelection")
    for (const dataSet of teamDataSets) {
        const option = document.createElement("option")
        option.value = dataSet.label
        option.text = dataSet.label
        teamSelect.appendChild(option)
    }
}

// addSeasonSelect - adds the seasons from the declared data array to the selection element
function addSeasonSelect() {
    const seasonSelect = document.getElementById("seasonSelection")
    for (const season of seasons) {
        const option = document.createElement("option")
        option.value = season
        option.text = `${season-1} - ${season}`
        seasonSelect.appendChild(option)
    }
}

// addTeamSelectionChangeHandler - adds the event listener to the team selection
function addTeamSelectionChangeHandler() {
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

// addSeasonSelectionChangeHandler - adds the event listener to the team selection
function addSeasonSelectionChangeHandler() {
    const seasonSelect = document.getElementById("seasonSelection")
    seasonSelect.addEventListener('change', () => {
        NHLStandingsChart.destroy()
        getSeasonData(seasonSelect.value)
    })
}

// resetTeamDataPoints - resets the team data points when a new season is selected
function resetTeamDataPoints() {
    for (const teamData of teamDataSets) {
        teamData.data = []
    }
}

addTeamLogo()
addTeamsSelect()
addSeasonSelect()
addTeamSelectionChangeHandler()
addSeasonSelectionChangeHandler()
getSeasonData(2020)