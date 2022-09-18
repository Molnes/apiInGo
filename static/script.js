
const fetchEvents = async (courseCode) => {
    const response = await fetch(`/events/${courseCode}`);
    const data = await response.json();
    return data;
}

const fetchEventsThisWeek = async (courseCode) => {

    const response = await fetch(`/eventsThisWeek/${courseCode}`);
    const data = await response.json();
    return data;
}

const fetchMultipleEventsThisWeek = async (courseCode,week) => {
   // fetch multiple events from the server with data from the courseCode array

    const response = await fetch(`/eventsByCourseCodesByWeek/${week}/${courseCode}`);
    const data = await response.json();

    console.log(data);
    return data;

}

const setTime = () => {
    let clock = document.getElementById("time");
    let time = new Date();
    let hours = time.getHours();
    let minutes = time.getMinutes();
    let seconds = time.getSeconds();
    let day = time.getDate();
    let month = time.getMonth() + 1;
    let year = time.getFullYear();

    clock.innerHTML = `${hours}:${minutes}:${seconds} ${day}/${month}/${year}`;
}

setInterval(setTime,1000);

    

const objectComparisonCallback = (a, b) => {
    if (a.from < b.from) {
        return -1;
    }
    if (a.from > b.from) {
        return 1;
    }
    return 0;
}

const showEvents = async (courseCode) => {
    const eventsThisWeek = await fetchEventsThisWeek(courseCode);
    const eventsList = document.getElementById("eventsList");
    eventsList.innerHTML = "";

    // sort the events array by date of values in the array
    eventsThisWeek.sort(objectComparisonCallback);


    eventsThisWeek.forEach(event => {
        const eventItem = document.createElement("div");
        eventItem.className = "eventItem";
        eventItem.innerHTML = `
            <div class="eventItem__title">${event.courseName.nameNob}</div>
            <div class="eventItem__date">${new Date(event.from).toLocaleString()}</div>
            <div class="eventItem__date">${new Date(event.to).toLocaleString()}</div>
            <div class="eventItem__location">${event.rooms[0].room}</div>
        `;
        eventsList.appendChild(eventItem);
    });


}
const showMultipleEvents = async (courseCode,week) => {
    const eventsThisWeek = await fetchMultipleEventsThisWeek(courseCode,week);
    const eventsList = document.getElementById("eventsList");
    eventsList.innerHTML = "";

    // display eventList as a flexbox
    eventsList.style.display = "flex";

    // sort the events array by date of values in the array
    eventsThisWeek.sort(objectComparisonCallback);

    eventsList.innerHTML = `
    <div class='eventDay' id='monday'><h1>Mandag</h1></div>
    <div class='eventDay' id='tuesday'><h1>Tirsdag</h1></div>
    <div class='eventDay' id='wednesday'><h1>Onsdag</h1></div>
    <div class='eventDay' id='thursday'><h1>Torsdag</h1></div>
    <div class='eventDay' id='friday'><h1>Fredag</h1></div>
    `;


    const eventItemTemplate = (event) => {
        return `
        <div class="eventItem">
            <div class="eventItem__title">${event.courseName.nameNob}</div>
            <div class="eventItem__date">${new Date(event.from).toLocaleString()}</div>
            <div class="eventItem__date">${new Date(event.to).toLocaleString()}</div>
            <div class="eventItem__location">${event.rooms[0].room} | ${event.rooms[0].building}</div>
        </div>
        `;
    }

    const monday = document.getElementById("monday");
    const tuesday = document.getElementById("tuesday");
    const wednesday = document.getElementById("wednesday");
    const thursday = document.getElementById("thursday");
    const friday = document.getElementById("friday");




    eventsThisWeek.forEach(event => {
        const eventItem = document.createElement("div");
        eventItem.className = "eventItem";
        eventItem.innerHTML = eventItemTemplate(event);
        switch (new Date(event.from).getDay()) {
            case 1:
                monday.appendChild(eventItem);
                break;
            case 2:
                tuesday.appendChild(eventItem);
                break;
            case 3:
                wednesday.appendChild(eventItem);
                break;
            case 4:
                thursday.appendChild(eventItem);
                break;
            case 5:
                friday.appendChild(eventItem);
                break;
            default:
                break;
        }

    });

    


}


//showEvents("IDATA2302");

showMultipleEvents(["IDATA2302", "IDATA2303", "IDATA2304", "ISTA1003"],"38");