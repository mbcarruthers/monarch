import React, {useState} from "react";
import {setupStore} from "@shipt/osmosis";

const useMapStore = () => {
    const [entityStore, updateEntities] = useState([]);

    const addEntities = (entities) => {
        updateEntities([...entities]);
    }
    const clearEntities = () => {
        updateEntities([]);
    }
    return{
        state:{
            entityStore
        },
        addEntities,
        clearEntities
    };
};

let MapStore = setupStore(useMapStore);
export default MapStore;