import React from "react";
import PropTypes from "prop-types";
import { withStyles } from "@material-ui/core/styles";
import Button from "@material-ui/core/Button";
import Dialog from "@material-ui/core/Dialog";
import DialogActions from "@material-ui/core/DialogActions";
import DialogContent from "@material-ui/core/DialogContent";
import DialogContentText from "@material-ui/core/DialogContentText";
import DialogTitle from "@material-ui/core/DialogTitle";
import IconButton from "@material-ui/core/IconButton";
import CloseIcon from "@material-ui/icons/Close";
import Typography from "@material-ui/core/Typography";
import PlayListAddIcon from "@material-ui/icons/PlaylistAdd";
import TextField from "@material-ui/core/TextField";
import MenuItem from "@material-ui/core/MenuItem";
import Select from "@material-ui/core/Select";
import FormControl from "@material-ui/core/FormControl";
import Input from "@material-ui/core/Input";
import InputLabel from "@material-ui/core/InputLabel";
import Checkbox from "@material-ui/core/Checkbox";
import ListItemText from "@material-ui/core/ListItemText";
import withMobileDialog from "@material-ui/core/withMobileDialog";
import api from "../lib/api";

const styles = theme => ({
  container: {
    display: "flex",
    flexWrap: "wrap",
    marginTop: "-40px"
  },
  formControl: {
    margin: theme.spacing.unit,
    minWidth: 120,
    maxWidth: 300
  },
  textField: {
    marginLeft: theme.spacing.unit,
    marginRight: theme.spacing.unit,
    width: 200
  },
  buttonField: {
    marginRight: theme.spacing.unit,
    width: "100%"
  },
  fullWidthField: {
    marginLeft: theme.spacing.unit,
    marginRight: theme.spacing.unit,
    width: "95%"
  },

  menu: {
    width: 200
  }
});

const ITEM_HEIGHT = 48;
const ITEM_PADDING_TOP = 8;
const MenuProps = {
  PaperProps: {
    style: {
      maxHeight: ITEM_HEIGHT * 4.5 + ITEM_PADDING_TOP,
      width: 250
    }
  }
};

class EditIptableDialog extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      open: false,
      roles: [
        {
          value: "all",
          label: "All"
        }
      ],
      id: -1,
      scope: "all",
      type: "allow",
      ip: ""
    };

    this.fetchRoles = this.fetchRoles.bind(this);
    this.handleSave = this.handleSave.bind(this);
    this.fetchIptable = this.fetchIptable.bind(this);
  }

  componentDidMount() {
    this.props.onRef(this);
    this.fetchRoles();
  }

  async fetchRoles() {
    let resp = await api.GetScopeTags();
    if (resp) {
      let roles = [...this.state.roles, ...resp.data];
      let selected = [];
      roles.map(item => {
        selected.push(item.value);
      });
      this.setState({
        roles: roles,
        selectedRole: selected
      });
    }
  }

  async fetchIptable(id) {
    let resp = await api.GetIptable(id);
    if (resp) {
      this.setState({
        scope: resp.data.scope,
        type: resp.data.type,
        ip: resp.data.ip.join(",")
      });
    }
  }

  async handleSave(id, data) {
    let resp = await api.SaveIptable(id, data);

    if (resp) {
      if (resp.data.success) {
        alert("Save.");
        this.handleClose();
        this.props.refresh();
      } else {
        alert(resp.data.error);
      }
    } else {
      alert("Error!");
    }
  }

  handleClickOpen = id => {
    this.fetchIptable(id);
    this.setState({
      open: true,
      id: id
    });
  };

  handleClose = () => {
    let selected = [];
    this.state.roles.map(item => {
      selected.push(item.value);
    });
    this.setState({
      open: false,
      id: -1,
      scope: "all",
      type: "allow",
      ip: ""
    });
  };

  render() {
    const { classes } = this.props;
    return (
      <div>
        <Dialog
          onClose={this.handleClose}
          aria-labelledby="responsive-dialog-title"
          open={this.state.open}
        >
          <DialogTitle id="responsive-dialog-title" onClose={this.handleClose}>
            Edit
          </DialogTitle>
          <form
            className={classes.container}
            autoComplete="off"
            onSubmit={e => {
              e.preventDefault();
              let data = {
                type: this.state.type,
                scope: this.state.scope,
                ip: this.state.ip.split(",")
              };
              this.handleSave(this.state.id, data);
            }}
          >
            <DialogContent>
              <TextField
                required
                fullWidth
                id="ip"
                label="IP"
                placeholder="Format x.x.x.x/mask,... ex.127.0.0.1/32,10.0.0.0/8"
                multiline
                rows="5"
                className={classes.fullWidthField}
                margin="normal"
                value={this.state.ip}
                onChange={e => {
                  this.setState({ ip: e.target.value });
                }}
              />
              <FormControl margin="normal" className={classes.textField}>
                <InputLabel htmlFor="type">Type</InputLabel>
                <Select
                  value={this.state.type}
                  onChange={e => {
                    this.setState({ type: e.target.value });
                  }}
                  input={<Input id="type" />}
                  MenuProps={MenuProps}
                >
                  <MenuItem key="allow" value="allow">
                    <ListItemText primary="Allow" />
                  </MenuItem>
                  <MenuItem key="deny" value="deny">
                    <ListItemText primary="Deny" />
                  </MenuItem>
                </Select>
              </FormControl>
              <FormControl margin="normal" className={classes.textField}>
                <InputLabel htmlFor="scope">Scope</InputLabel>
                <Select
                  value={this.state.scope}
                  onChange={e => {
                    this.setState({ scope: e.target.value });
                  }}
                  input={<Input id="scope" />}
                  MenuProps={MenuProps}
                >
                  {this.state.roles.map(role => (
                    <MenuItem key={role.value} value={role.value}>
                      <ListItemText primary={role.label} />
                    </MenuItem>
                  ))}
                </Select>
              </FormControl>

              <FormControl
                fullWidth
                margin="normal"
                style={{ marginTop: "10%" }}
                className={classes.fullWidth}
              >
                <Button color="secondary" size="large" type="submit">
                  Save
                </Button>
              </FormControl>
            </DialogContent>
          </form>
        </Dialog>
      </div>
    );
  }
}

EditIptableDialog.propTypes = {
  classes: PropTypes.object.isRequired
};

export default withStyles(styles)(withMobileDialog()(EditIptableDialog));
