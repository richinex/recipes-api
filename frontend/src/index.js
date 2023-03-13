import React from "react";
import ReactDOM from "react-dom/client";
import App from "./App";
import './index.css';
import 'bootstrap/dist/css/bootstrap.min.css';
import 'bootstrap/dist/js/bootstrap.min.js';
import { Auth0Provider } from "@auth0/auth0-react";

const { REACT_APP_AUTH0_DOMAIN } = process.env;
const { REACT_APP_AUTH0_CLIENT_ID } = process.env;

const root = ReactDOM.createRoot(document.getElementById("root"));
root.render(
  <Auth0Provider
    domain={REACT_APP_AUTH0_DOMAIN}
    clientId={REACT_APP_AUTH0_CLIENT_ID}
    redirect_uri={window.location.origin}
  >
    <App />
  </Auth0Provider>
);


