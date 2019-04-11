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

var base = document.getElementsByTagName("base")[0].href;
base = base.substring(base.indexOf("//") + 2, base.length);
if (base.indexOf("/") > -1) {
  base = base.substring(base.indexOf("/"), base.length);
}
const basehref = base;
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
  <BrowserRouter basename={basehref}>
    <App />
  </BrowserRouter>,
  document.getElementById("root")
);

// If you want your app to work offline and load faster, you can change
// unregister() to register() below. Note this comes with some pitfalls.
// Learn more about service workers: http://bit.ly/CRA-PWA
serviceWorker.unregister();
