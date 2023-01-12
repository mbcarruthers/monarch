import React from "react";
import "../style/Landing.css";

// 49766 48662

const Landing = () => {
    return(
        <div className="container-fluid text-center outline">
            <h1>Landing</h1>
            <div className="container-fluid outline">
                <img src={"\/\/\localhost:8025/btrfly.jpg"} alt="not working" className="btrfly-image m-2"/>
            </div>
        </div>
    )
}

// <img src={"\/\/\localhost:8025/btrfly.jpg"} alt="not working"/>

export default Landing;