import { MapContainer, TileLayer, Popup,CircleMarker, LayerGroup } from 'react-leaflet'
import MapControls from "./MapControls";
import "../style/MapView.css";
import MapStore from "../MapStore";
import {DateTime} from "luxon";
import React,{useContext} from "react";

// monthColors - set of 12 colors to represent the months of the year
const monthColor = [
    "#1C2331", // navy blue - January
    "#5C7FFF", // baby blue - February
    "#0F9B0F", // forest green - March
    "#48C9B0", // mint green - April
    "#9B0F0F", // maroon - May
    "#FC4F4F", // salmon - June
    "#9B910F", // olive green - July
    "#F7DC6F", // pale goldenrod - August
    "#F7A631", // tangerine - September
    "#7DCEF2", // baby blue - October
    "#1E90FF", // dodger blue - November
    "#87CEFA", // light blue - December
]

// compare - used to sort observations chronologically
const compare = (a, b) => {
    if(a.observed_on < b.observed_on) {
        return -1;
    }
    if (a.observed_on > b.observed_on) {
        return 1
    }
    return 0;
}
// setMarkerStyle - provides the correct marker color for the associated month observation
// was observed in
const setMarkerStyle = (dateObserved) => {
    // const monthNumber = new Date(dateObserved).getUTCMonth();
    const monthNumber = DateTime.fromFormat(dateObserved,"yyyy-MM-dd").month - 1;
    return {
        color: monthColor[monthNumber],
        fillColor: monthColor[monthNumber],
        fillOpacity: 0.5,
        radius: 11,
    }
}
// parseYear - because Date().getFullYear() return 2011 for '2012-01-01'
const parseYear = (dateObserved) => {
    return parseInt(dateObserved);
}
// pareMonth - Exists because getMonth() was returning 11 for the month of '2012-01-01'
const parseMonth = (dateObserved) => {
    return new Date(dateObserved).getUTCMonth() + 1;
}
// sortByMonthAndYear - sort an array of observations by month and year, returns a 2 dimensional array
// of sorted observations.
// Note: Not sure how necessary this is, but even if not it may proved by be useful
function sortByMonthAndYear(entities) {
    const elements = entities.sort(compare);
    const result = [];
    let currentMonth = null;
    let currentYear = null;
    let currentArray = null;

    for (const element of elements) {

        const month = parseMonth(element.observed_on);
        const year = parseYear(element.observed_on);

        if (currentMonth === null || month !== currentMonth || year !== currentYear) {

            currentMonth = month;
            currentYear = year;
            currentArray = [];
            result.push(currentArray);
        }
        currentArray.push(element);
    }
    console.log(result);
    return result;
}

// createObservationMarkers - creates observation markers
// Note: Only makes red markers
// expects- an array of entities matching the json struct Entity in helio project
// returns- an array of leaflet observation markers made for those entities given as an argument
const createObservationMarkers = (entities) => {
    return entities.map( (entity,idx) => {
        let _link = `https://www.inaturalist.org/observations/${entity.id}`;
        return <CircleMarker center={[entity.latitude,entity.longitude]}
                             color="red"
                             fillColor="#f03"
                             fillOpacity={0.5}
                             radius={5}
                             key={idx}>
            <Popup className="popup">
                <small className="">Locale:{entity.place_guess}</small><br/>
                <small>{entity.species_guess}</small><br/>
                <small>Date:{entity.observed_on}</small><br/>
                <small><a target="_blank" rel="noopener noreferrer" href={_link}>iNaturalist</a></small>
            </Popup>
        </CircleMarker>
    })
}

const createMarkers = (entities) => {
    const elements = sortByMonthAndYear(entities);

    return elements.map((element, index) => {
        return(
            <div>
                {
                    element.map((observation,idx) => {
                        const _link = `https://www.inaturalist.org/observations/${observation.id}`;
                        return(
                         <CircleMarker center={[observation.latitude, observation.longitude]}
                                       {...setMarkerStyle(observation.observed_on)}
                                             key={idx}>
                            <Popup className="popup">
                                <small className="">Locale:{observation.place_guess}</small><br/>
                                <small>{observation.species_guess}</small><br/>
                                <small>Date:{observation.observed_on}</small><br/>
                                <small><a target="_blank" rel="noopener noreferrer" href={_link}>iNaturalist</a></small>
                            </Popup>
                        </CircleMarker>
                        )
                    })
                }
            </div>
        )
    })
}


// MapView serves as the view for the map component.
const MapView = () => {
    const mapStore = useContext(MapStore.Context);
    let {entityStore} = mapStore.state;
    return(
        <MapContainer center={[27.732161, -84.00095]} zoom={6} scrollWheelZoom={true} id="mapview">
            <TileLayer
                attribution='&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
                url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
            />
            <LayerGroup>
                {
                    entityStore.length === 0 ? undefined : createMarkers(entityStore)
                }
            </LayerGroup>

           <MapControls></MapControls>
        </MapContainer>
    )
}
export default MapView;

// createObservationMarkers(entityStore)