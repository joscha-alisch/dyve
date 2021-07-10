import React from 'react';
import ReactDOM from 'react-dom';
import './index.css';
import App from './App';
import {
    BrowserRouter as Router, Route,
} from "react-router-dom";
import { QueryParamProvider } from 'use-query-params';

ReactDOM.render(
  <React.StrictMode>
      <Router>
          <QueryParamProvider ReactRouterRoute={Route}>
              <App />
          </QueryParamProvider>
      </Router>
  </React.StrictMode>,
  document.getElementById('root')
);

