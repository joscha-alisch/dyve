import React from 'react';
import {AuthProvider} from "./context/auth";
import {NotificationProvider} from "./context/notifications";
import { ReactLocationDevtools } from 'react-location-devtools'
import { Outlet, ReactLocation, Router } from "react-location"
import { routes } from "./config/routes"

const location = new ReactLocation()

function App() {
  return <Router location={location} routes={routes} >
  <NotificationProvider>
      <AuthProvider>
          <Outlet />
          <ReactLocationDevtools initialIsOpen={false} />
      </AuthProvider>
  </NotificationProvider>
</Router>;
}

export default App;
