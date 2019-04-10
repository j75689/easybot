import React from "react";
import PropTypes from "prop-types";
import { withStyles } from "@material-ui/core/styles";
import Button from "@material-ui/core/Button";
import Dialog from "@material-ui/core/Dialog";
import MuiDialogTitle from "@material-ui/core/DialogTitle";
import MuiDialogContent from "@material-ui/core/DialogContent";
import MuiDialogActions from "@material-ui/core/DialogActions";
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

const DialogTitle = withStyles(theme => ({
  root: {
    borderBottom: `1px solid ${theme.palette.divider}`,
    margin: 0,
    padding: theme.spacing.unit * 2
  },
  closeButton: {
    position: "absolute",
    right: theme.spacing.unit,
    top: theme.spacing.unit,
    color: theme.palette.grey[500]
  }
}))(props => {
  const { children, classes, onClose } = props;
  return (
    <MuiDialogTitle disableTypography className={classes.root}>
      <Typography variant="h6">{children}</Typography>
      {onClose ? (
        <IconButton
          aria-label="Close"
          className={classes.closeButton}
          onClick={onClose}
        >
          <CloseIcon />
        </IconButton>
      ) : null}
    </MuiDialogTitle>
  );
});

const styles = theme => ({
  container: {
    display: "flex",
    flexWrap: "wrap"
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
  dense: {
    marginTop: 19
  },
  menu: {
    width: 200
  }
});

const DialogContent = withStyles(theme => ({
  root: {
    margin: 0,
    padding: theme.spacing.unit * 2
  }
}))(MuiDialogContent);

const DialogActions = withStyles(theme => ({
  root: {
    borderTop: `1px solid ${theme.palette.divider}`,
    margin: 0,
    padding: theme.spacing.unit
  }
}))(MuiDialogActions);

const roles = [
  {
    value: "all",
    label: "All"
  },
  {
    value: "push",
    label: "Push"
  },
  {
    value: "multicast",
    label: "Multicast"
  }
];

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

class NewRoleAccountDialog extends React.Component {
  state = {
    open: false,
    selectedRole: ["all", "push", "multicast"]
  };

  handleClickOpen = () => {
    this.setState({
      open: true
    });
  };

  handleClose = () => {
    this.setState({ open: false });
  };

  render() {
    const { classes } = this.props;
    return (
      <div>
        <IconButton aria-label="New" onClick={this.handleClickOpen}>
          <PlayListAddIcon />
        </IconButton>

        <Dialog
          onClose={this.handleClose}
          aria-labelledby="customized-dialog-title"
          open={this.state.open}
        >
          <DialogTitle id="customized-dialog-title" onClose={this.handleClose}>
            Create
          </DialogTitle>
          <DialogContent>
            <form className={classes.container} autoComplete="off">
              <TextField
                required
                id="account"
                label="Account Name"
                className={classes.textField}
                margin="normal"
              />
              <TextField
                //error={false}
                id="email"
                label="Email"
                className={classes.textField}
                margin="normal"
              />
              <TextField
                required
                id="domain"
                label="Domain"
                style={{ margin: 8 }}
                fullWidth
                margin="normal"
              />
              <TextField
                required
                //error={false}
                id="provider"
                label="Provider"
                fullWidth
                style={{ margin: 8 }}
                margin="normal"
              />
              <TextField
                //error={false}
                id="active"
                label="Active"
                defaultValue="7200"
                className={classes.textField}
                placeholder="0 is no limit"
                margin="normal"
              />
              <FormControl margin="normal" className={classes.textField}>
                <InputLabel htmlFor="role">Role</InputLabel>
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
                      roles.map(obj => {
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
                  input={<Input id="role" />}
                  renderValue={selected =>
                    this.state.selectedRole.indexOf("all") > -1
                      ? "all"
                      : selected.join(",")
                  }
                  MenuProps={MenuProps}
                >
                  {roles.map(role => (
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
            </form>
          </DialogContent>
          <DialogActions>
            <Button onClick={this.handleClose} color="primary">
              Add
            </Button>
          </DialogActions>
        </Dialog>
      </div>
    );
  }
}

NewRoleAccountDialog.propTypes = {
  classes: PropTypes.object.isRequired
};

export default withStyles(styles)(NewRoleAccountDialog);
