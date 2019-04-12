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

class EditRoleAccountDialog extends React.Component {
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
      selectedRole: ["all"],
      account: "",
      email: "",
      domain: "",
      provider: "",
      active: "7200"
    };

    this.fetchRoles = this.fetchRoles.bind(this);
    this.fetchAccountInfo = this.fetchAccountInfo.bind(this);
    this.handleSaveAccount = this.handleSaveAccount.bind(this);
    this.handleRefreshToken = this.handleRefreshToken.bind(this);
  }

  componentDidMount() {
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

  async fetchAccountInfo(name) {
    let resp = await api.GetServiceAccount(name);
    if (resp) {
      if (resp.data) {
        let scope = resp.data.scope.split(",");
        if (scope.indexOf("all") > -1) {
          scope = [];
          this.state.roles.map(item => {
            scope.push(item.value);
          });
        }

        this.setState({
          account: resp.data.name,
          email: resp.data.email,
          domain: resp.data.domain,
          provider: resp.data.provider,
          active: resp.data.active,
          selectedRole: scope
        });
      }
    }
  }

  handleClickOpen = () => {
    this.fetchAccountInfo(this.props.account);
    this.setState({
      open: true
    });
  };

  handleClose = () => {
    let selected = [];
    this.state.roles.map(item => {
      selected.push(item.value);
    });
    this.setState({
      open: false,
      selectedRole: selected,
      account: "",
      email: "",
      domain: "",
      provider: "",
      active: "7200"
    });
  };

  async handleRefreshToken(name) {
    let resp = await api.RefreshServiceAccountToken(name);
    if (resp) {
      if (resp.data) {
        if (resp.data.success) {
          alert("Refresh.");
          this.props.refresh();
        } else {
          alert(resp.data.error);
        }
      }
    } else {
      alert("Error!");
    }
  }

  async handleSaveAccount(name, data) {
    let resp = await api.SaveServiceAccount(name, data);

    if (resp) {
      if (resp.data.success) {
        alert("Saved.");
        this.handleClose();
        this.props.refresh();
      } else {
        alert(resp.data.error);
      }
    } else {
      alert("Error!");
    }
  }

  render() {
    const { classes } = this.props;
    return (
      <div>
        <Button
          color="primary"
          className={classes.button}
          onClick={this.handleClickOpen}
        >
          Edit
        </Button>
        <Button
          color="secondary"
          className={classes.button}
          onClick={e => {
            this.handleRefreshToken(this.props.account);
          }}
        >
          Refresh
        </Button>
        <Dialog
          onClose={this.handleClose}
          aria-labelledby="responsive-dialog-title"
          open={this.state.open}
        >
          <DialogTitle id="responsive-dialog-title" onClose={this.handleClose}>
            Create
          </DialogTitle>
          <form
            className={classes.container}
            autoComplete="off"
            onSubmit={e => {
              e.preventDefault();
              let data = {
                name: this.state.account,
                email: this.state.email,
                domain: this.state.domain,
                provider: this.state.provider,
                active: this.state.active,
                scope:
                  this.state.selectedRole.indexOf("all") > -1
                    ? "all"
                    : this.state.selectedRole.join(",")
              };
              this.handleSaveAccount(this.props.account, data);
            }}
          >
            <DialogContent>
              <TextField
                required
                id="account"
                label="Account Name"
                className={classes.textField}
                margin="normal"
                value={this.state.account}
                onChange={e => {
                  this.setState({
                    account: e.target.value
                  });
                }}
              />
              <TextField
                id="email"
                label="Email"
                className={classes.textField}
                margin="normal"
                value={this.state.email}
                onChange={e => {
                  this.setState({
                    email: e.target.value
                  });
                }}
              />
              <TextField
                required
                id="domain"
                label="Domain"
                className={classes.fullWidthField}
                fullWidth
                margin="normal"
                value={this.state.domain}
                onChange={e => {
                  this.setState({
                    domain: e.target.value
                  });
                }}
              />
              <TextField
                required
                id="provider"
                label="Provider"
                fullWidth
                className={classes.fullWidthField}
                margin="normal"
                value={this.state.provider}
                onChange={e => {
                  this.setState({
                    provider: e.target.value
                  });
                }}
              />
              <TextField
                id="active"
                label="Active"
                className={classes.textField}
                placeholder="0 is no limit"
                margin="normal"
                value={this.state.active}
                onChange={e => {
                  this.setState({
                    active: e.target.value
                  });
                }}
              />
              <FormControl margin="normal" className={classes.textField}>
                <InputLabel htmlFor="scope">Scope</InputLabel>
                <Select
                  required
                  multiple
                  value={this.state.selectedRole}
                  onChange={e => {
                    if (
                      e.target.value.indexOf("all") > -1 &&
                      this.state.selectedRole.indexOf("all") == -1
                    ) {
                      var all = [];
                      this.state.roles.map(obj => {
                        all.push(obj.value);
                      });
                      this.setState({ selectedRole: all });
                    } else if (
                      e.target.value.indexOf("all") == -1 &&
                      this.state.selectedRole.indexOf("all") > -1
                    ) {
                      this.setState({ selectedRole: [] });
                    } else if (
                      e.target.value.indexOf("all") == -1 &&
                      this.state.selectedRole.indexOf("all") == -1
                    ) {
                      this.setState({ selectedRole: e.target.value });
                    }
                  }}
                  input={<Input id="scope" />}
                  renderValue={selected =>
                    this.state.selectedRole.indexOf("all") > -1
                      ? "all"
                      : selected.join(",")
                  }
                  MenuProps={MenuProps}
                >
                  {this.state.roles.map(role => (
                    <MenuItem key={role.label} value={role.value}>
                      <Checkbox
                        checked={
                          this.state.selectedRole.indexOf(role.value) > -1
                        }
                      />
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

EditRoleAccountDialog.propTypes = {
  classes: PropTypes.object.isRequired
};

export default withStyles(styles)(withMobileDialog()(EditRoleAccountDialog));
