import React from 'react';
import { useAuth0 } from "@auth0/auth0-react";
import Profile from './Profile';
import { NavLink } from 'react-router-dom';

const Navbar = () => {
  const { isAuthenticated, loginWithRedirect, logout, user } = useAuth0();
  return (
      <nav className="navbar navbar-expand-lg navbar-light bg-light">
         <div className="navbar-nav">
          <NavLink className="navbar-brand" to="/">Recipes</NavLink>
          <button className="navbar-toggler" type="button" data-toggle="collapse" data-target="#navbarTogglerDemo02" aria-controls="navbarTogglerDemo02" aria-expanded="false" aria-label="Toggle navigation">
              <span className="navbar-toggler-icon"></span>
          </button>
          </div>
          <div className="collapse navbar-collapse" id="navbarTogglerDemo02">
              <ul className="navbar-nav ml-auto">
                  <li className="nav-item">
                      {isAuthenticated ? (<Profile />) : (
                          <a className="nav-link active" onClick={() => loginWithRedirect()}> Login</a>
                      )}
                  </li>
              </ul>
          </div>
      </nav >
  )
}

export default Navbar;
