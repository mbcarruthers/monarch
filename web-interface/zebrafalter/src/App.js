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
                  <nav className="nav-row navbar">
                      <Link to="/" className="no-decor mx-1 px-1 nav-link">
                          <h4 className="nav-text">Home</h4>
                      </Link>
                      <Link to="/observations" className="no-decor mx-1 px-1 nav-link">
                          <h4 className="nav-text">Monarch Map</h4>
                      </Link>

                  </nav>
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

