import React, { Component } from "react";
import ReactDOM from "react-dom";
import Dashboard from "./pages/dashboard";
import SignIn from "./pages/sign-in";
import AccessRole from "./pages/access-role";
import * as serviceWorker from "./serviceWorker";
import { BrowserRouter, Switch, Route, Redirect } from "react-router-dom";

class App extends Component {
  render() {
    var loc = window.location;
    var prefix = "";
    if (loc.pathname.lastIndexOf("/") > -1) {
      prefix = loc.pathname.substring(0, loc.pathname.lastIndexOf("/"));
    }

    return (
      <>
        <Switch>
          <Route exact path={prefix + "/dashboard"} component={Dashboard} />
          <Route exact path={prefix + "/accessrole"} component={AccessRole} />
          <Route exact path={prefix + "/login"} component={SignIn} />
        </Switch>
      </>
    );
  }
}

ReactDOM.render(
  <BrowserRouter>
    <App />
  </BrowserRouter>,
  document.getElementById("root")
);

// If you want your app to work offline and load faster, you can change
// unregister() to register() below. Note this comes with some pitfalls.
// Learn more about service workers: http://bit.ly/CRA-PWA
serviceWorker.unregister();
