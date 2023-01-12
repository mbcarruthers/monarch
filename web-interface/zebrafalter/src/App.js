import React from "react";
import "./App.css";
import {BrowserRouter,Routes,Route,Link} from "react-router-dom";
import Landing from "./components/Landing";
import MapView from "./components/MapView";
import MapStore from "./MapStore";
import {StoreProvider} from "@shipt/osmosis";

const App = () => {
  return(
      <BrowserRouter>
          <div className="app">
              <header className="p-3 bg-dark text-white container-fluid text-center sticky-top">
                  <div className="container">
                      <div className="d-flex flex-wrap align-items-center justify-content-center justify-content-lg-start text-center">
                          <ul className="nav col-12 col-lg-auto me-lg-auto mb-2 justify-content-center mb-md-0">
                              <li><Link to="/" className="nav-link px-2 change-color no-decor btn btn-outline-light m-2">Home</Link></li>
                              <li><Link to="/observations" className="nav-link px-2 change-color no-decor btn btn-outline-light m-2">Monarch Map</Link></li>
                          </ul>
                      </div>
                  </div>
              </header>

              <div className="local">
                  <Routes>
                      <Route path="/" element={<Landing/>}/>
                      <Route path="/observations" element={<MapView/>} />
                  </Routes>
              </div>
              <div className="bg-dark text-white container-fluid text-center fixed-bottom">
                 <h2>Footer</h2>
              </div>
          </div>
      </BrowserRouter>
  )
}

export default StoreProvider([MapStore.Provider], App);

