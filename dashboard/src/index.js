import React, { Component } from "react";
import ReactDOM from "react-dom";
import Dashboard from "./pages/dashboard";
import SignIn from "./pages/sign-in";
import AccessRole from "./pages/access-role";
import * as serviceWorker from "./serviceWorker";
import {
  BrowserRouter,
  Switch,
  Route,
  Redirect,
  useRouterHistory
} from "react-router-dom";
import Config from "./lib/config";
class App extends Component {
  render() {
    return (
      <>
        <Switch>
          <Route exact path="/dashboard" component={Dashboard} />
          <Route exact path="/accessrole" component={AccessRole} />
          <Route exact path="/login" component={SignIn} />
        </Switch>
      </>
    );
  }
}

ReactDOM.render(
  <BrowserRouter basename={Config.Basehref}>
    <App />
  </BrowserRouter>,
  document.getElementById("root")
);

// If you want your app to work offline and load faster, you can change
// unregister() to register() below. Note this comes with some pitfalls.
// Learn more about service workers: http://bit.ly/CRA-PWA
serviceWorker.unregister();
