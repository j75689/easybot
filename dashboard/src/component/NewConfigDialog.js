import React from "react";
import PropTypes from "prop-types";
import Button from "@material-ui/core/Button";
import Dialog from "@material-ui/core/Dialog";
import DialogActions from "@material-ui/core/DialogActions";
import DialogContent from "@material-ui/core/DialogContent";
import DialogContentText from "@material-ui/core/DialogContentText";
import DialogTitle from "@material-ui/core/DialogTitle";
import withMobileDialog from "@material-ui/core/withMobileDialog";
import { withStyles } from "@material-ui/core/styles";
import Input from "@material-ui/core/Input";
import InputLabel from "@material-ui/core/InputLabel";
import FormControl from "@material-ui/core/FormControl";

import "brace";
import "brace/mode/json";
import "brace/theme/github";
import { JsonEditor as Editor } from "jsoneditor-react";
import "jsoneditor-react/es/editor.min.css";
import "../css/jsoneditor.css";
import api from "../lib/api";

const styles = theme => ({
  button: {
    margin: theme.spacing.unit
  },
  container: {
    display: "flex",
    flexWrap: "wrap"
  },
  formControl: {
    margin: theme.spacing.unit
  }
});

class NewConfigDialog extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      open: false,
      noFillIDErr: false,
      configID: ""
    };
    this.handleCreateConfig = this.handleCreateConfig.bind(this);
  }

  componentDidMount() {
    this.props.onRef(this);
  }

  handleClickOpen = () => {
    this.setState({ open: true });
  };

  handleClose = () => {
    this.setState({ open: false, noFillIDErr: false, configID: "" });
  };

  async handleCreateConfig(id, data) {
    let resp = await api.CreateHandlerConfig(id, data);
    if (resp) {
      if (resp.data.success) {
        alert("Created.");
        this.handleClose();
        this.props.refresh();
      } else {
        alert(resp.data.error);
      }
    } else {
      alert("Error!");
    }
  }

  setJsonEditorRef = instance => {
    if (instance) {
      const { jsonEditor } = instance;
      this.jsonEditor = jsonEditor;
    } else {
      this.jsonEditor = null;
    }
  };

  render() {
    const { fullScreen } = this.props;
    const { classes } = this.props;

    return (
      <div>
        <Dialog
          fullScreen={fullScreen}
          open={this.state.open}
          onClose={this.handleClose}
          aria-labelledby="responsive-dialog-title"
        >
          <DialogTitle id="responsive-dialog-title">
            Create Handler Config
          </DialogTitle>
          <DialogContent>
            <FormControl
              error={this.state.noFillIDErr}
              required
              className={classes.formControl}
            >
              <InputLabel htmlFor="component-simple">ID</InputLabel>
              <Input
                id="component-simple"
                value={this.state.configID}
                onChange={e => {
                  this.setState({ configID: e.target.value });
                }}
              />
            </FormControl>
            <Editor
              ref={this.setJsonEditorRef}
              mode="code"
              name="content"
              value={{}}
            />
          </DialogContent>
          <DialogActions>
            <Button
              onClick={this.handleClose}
              color="secondary"
              autoFocus
              onClick={e => {
                if (this.state.configID != "") {
                  this.handleCreateConfig(
                    this.state.configID,
                    this.jsonEditor.get()
                  );
                } else {
                  this.setState({
                    noFillIDErr: true
                  });
                }
              }}
            >
              Create
            </Button>
            <Button onClick={this.handleClose} color="primary">
              Cancel
            </Button>
          </DialogActions>
        </Dialog>
      </div>
    );
  }
}

NewConfigDialog.propTypes = {
  fullScreen: PropTypes.bool.isRequired
};

export default withStyles(styles)(
  withMobileDialog({ breakpoint: "lg" })(NewConfigDialog)
);
