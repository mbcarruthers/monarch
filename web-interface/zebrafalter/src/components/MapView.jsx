import { MapContainer, TileLayer, Popup,CircleMarker, LayerGroup } from 'react-leaflet'
import MapControls from "./MapControls";
import "../style/MapView.css";
import MapStore from "../MapStore";

import React,{useContext} from "react";

const ZlTid = 49766; // value 'constants' for temporary ease-ability in testing
const MonarchTid = 48662;

// createObservationMarkers - creates observation markers
// expects
//  - an array of entities matching the json struct Entity in helio project
// returns
//  - an array of leaflet observation markers made for those entities given as an argument
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
                    entityStore.length === 0 ? undefined : createObservationMarkers(entityStore)
                }
            </LayerGroup>

           <MapControls></MapControls>
        </MapContainer>
    )
}
export default MapView;

