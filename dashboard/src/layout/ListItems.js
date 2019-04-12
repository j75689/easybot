import React from "react";
import ListItem from "@material-ui/core/ListItem";
import ListItemIcon from "@material-ui/core/ListItemIcon";
import ListItemText from "@material-ui/core/ListItemText";
import ListSubheader from "@material-ui/core/ListSubheader";
import DashboardIcon from "@material-ui/icons/Dashboard";
import NoteIcon from "@material-ui/icons/Note";
import SwapHorizIcon from "@material-ui/icons/SwapHoriz";
import { Link } from "react-router-dom";
export const mainListItems = (
  <div>
    <Link to="/dashboard">
      <ListItem button>
        <ListItemIcon>
          <DashboardIcon />
        </ListItemIcon>
        <ListItemText primary="Dashboard" />
      </ListItem>
    </Link>
    <Link to="/config">
      <ListItem button>
        <ListItemIcon>
          <NoteIcon />
        </ListItemIcon>
        <ListItemText primary="Config" />
      </ListItem>
    </Link>
  </div>
);

export const systemSettingListItems = (
  <div>
    <ListSubheader inset>System Setting</ListSubheader>
    <Link to="/accessrole">
      <ListItem button>
        <ListItemIcon>
          <SwapHorizIcon />
        </ListItemIcon>
        <ListItemText primary="Access Role" />
      </ListItem>
    </Link>
  </div>
);
