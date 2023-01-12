import React, {useContext} from "react";
import Control from "react-leaflet-custom-control";
import "../style/MapView.css";
import MapStore from "../MapStore";

// Note: Route  /entities/search?taxon_id=XXX&date1=yy-nn-dd&date2=yy-mm-dd

const MapControls = (props) => {
    const mapStore = useContext(MapStore.Context);
    let {entityStore} = mapStore.state;
    const handleSubmitDate = () => {
        const date1 = document.getElementById("date1");
        const date2 = document.getElementById("date2")
        let d1 = date1.value;
        let d2 = date2.value; // Todo: Just pass date1.value & date2.value to the url
        console.log(d1,d2);
        fetch(`\/\/\localhost:8000/entities/search?date1="${d1}"&date2="${d2}"`)
            .then((res) => res.json())
            .then((entities) => {
                if(entities.length !== 0) {
                    mapStore.addEntities(entities)
                }
            })
            .catch((err) => console.error(err));
    }

    return (
        <Control prepend position="topright">
            <div className="form-group date-input-group">
                { entityStore.length === 0 ? undefined : <span className="badge count-box p-1 m-1">Count {entityStore.length}</span>}
                <label htmlFor="date1" className="m-1">Start date</label>
                <input type="date" id="date1" aria-valuemin={2012 - 1 - 1} aria-valuemax={2022 - 12 - 31} className="p-1"/>
                <label htmlFor="date2" className="m-1">End date</label>
                <input type="date" id="date2" aria-valuemin="2012-01-01" aria-valuemax="2022-12-31" className="p-1"/>
                <button className="btn btn-outline-info m-1 p-1" id="date-submit-button" onClick={handleSubmitDate}>Submit Dates</button>

            </div>
        </Control>
    )
};

export default MapControls;