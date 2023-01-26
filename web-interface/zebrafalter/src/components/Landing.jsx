import React from "react";
import "../style/Landing.css";

// Monarch and ZebraLongwing taxon_id's 49766 48662

const Landing = () => {
    return(
        <div className="container-fluid text-center outline">
            <h1 className="comfort">Landing</h1>
            <div className="container-fluid outline">
                <img src={"\/\/\localhost:8025/btrfly.jpg"} alt="not working" className="btrfly-image m-2"/>
            </div>
        </div>
    )
}

// <img src={"\/\/\localhost:8025/btrfly.jpg"} alt="not working"/>

export default Landing;