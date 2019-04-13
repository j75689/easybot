import React from "react";
import ReactDOM from "react-dom";
import PropTypes from "prop-types";
import { withStyles } from "@material-ui/core/styles";
import Input from "@material-ui/core/Input";
import OutlinedInput from "@material-ui/core/OutlinedInput";
import FilledInput from "@material-ui/core/FilledInput";
import InputLabel from "@material-ui/core/InputLabel";
import Button from "@material-ui/core/Button";
import MenuItem from "@material-ui/core/MenuItem";
import FormHelperText from "@material-ui/core/FormHelperText";
import FormControl from "@material-ui/core/FormControl";
import Select from "@material-ui/core/Select";
import "brace";
import "brace/mode/json";
import "brace/theme/github";
import { JsonEditor as Editor } from "jsoneditor-react";
import "jsoneditor-react/es/editor.min.css";
import "../css/jsoneditor.css";
import api from "../lib/api";
import NewConfigDialog from "./NewConfigDialog";

const styles = theme => ({
  root: {
    display: "flex",
    flexWrap: "wrap"
  },
  formControl: {
    margin: theme.spacing.unit,
    minWidth: 150
  },
  formControlSpan: {
    margin: theme.spacing.unit,
    minWidth: 50
  },
  selectEmpty: {
    marginTop: theme.spacing.unit * 2
  },
  button: {
    margin: theme.spacing.unit
  },
  input: {
    display: "none"
  }
});

class ConfigEditor extends React.Component {
  constructor(props) {
    super();
    this.state = {
      selectEvent: "",
      selectConfigID: "",
      configs: {
        //test: ["test1", "test2"]
      }
    };

    this.fetchConfig = this.fetchConfig.bind(this);
    this.fetchConfigID = this.fetchConfigID.bind(this);
    this.saveConfig = this.saveConfig.bind(this);
    this.deleteConfig = this.deleteConfig.bind(this);
  }

  componentDidMount() {
    this.fetchConfigID();
  }

  async fetchConfigID() {
    let resp = await api.GetAllConfigIDs();
    if (resp) {
      if (resp.data) {
        this.setState({
          configs: resp.data
        });
      }
    }
  }

  async fetchConfig(id) {
    let resp = await api.GetHandlerConfig(id);
    if (resp) {
      if (resp.data) {
        this.jsonEditor.set(resp.data);
      } else {
        this.jsonEditor.set();
      }
    }
  }

  async saveConfig(id, data) {
    let resp = await api.SaveHandlerConfig(id, data);
    if (resp) {
      if (resp.data.success) {
        alert("Saved.");
      } else {
        alert(resp.data.error);
      }
    } else {
      alert("Error!");
    }
  }

  async deleteConfig(id) {
    let resp = await api.DeleteHandlerConfig(id);
    if (resp) {
      if (resp.data.success) {
        await this.fetchConfigID();
        this.setState({
          selectEvent: "",
          selectConfigID: ""
        });
        this.jsonEditor.set({});
        alert("Delete.");
      } else {
        alert(resp.data.error);
      }
    } else {
      alert("Error!");
    }
  }

  handleChange = event => {
    this.setState({ [event.target.name]: event.target.value });
  };

  setJsonEditorRef = instance => {
    if (instance) {
      const { jsonEditor } = instance;
      this.jsonEditor = jsonEditor;
    } else {
      this.jsonEditor = null;
    }
  };

  setNewConfigDialog = ref => {
    this.newConfigDialog = ref;
  };

  render() {
    const { classes } = this.props;

    return (
      <div>
        <form className={classes.root} autoComplete="off">
          <FormControl className={classes.formControl}>
            <Select
              value={this.state.selectEvent}
              onChange={e => {
                this.setState({
                  selectEvent: e.target.value,
                  selectConfigID: ""
                });
              }}
              name="event"
              displayEmpty
              className={classes.selectEmpty}
            >
              <MenuItem value="" disabled>
                Choose
              </MenuItem>
              {Object.keys(this.state.configs).map(key => {
                return <MenuItem value={key}>{key}</MenuItem>;
              })}
            </Select>
            <FormHelperText>Event </FormHelperText>
          </FormControl>
          <FormControl className={classes.formControl}>
            <Select
              value={this.state.selectConfigID}
              onChange={e => {
                this.setState({
                  selectConfigID: e.target.value
                });
                this.fetchConfig(e.target.value);
              }}
              name="configID"
              displayEmpty
              className={classes.selectEmpty}
            >
              <MenuItem value="" disabled>
                Choose
              </MenuItem>
              {this.state.configs[this.state.selectEvent] &&
                this.state.configs[this.state.selectEvent].map(item => {
                  return <MenuItem value={item}>{item}</MenuItem>;
                })}
            </Select>
            <FormHelperText>Config </FormHelperText>
          </FormControl>
          <FormControl className={classes.formControl} />
          <Button
            color="primary"
            className={classes.button}
            onClick={e => {
              if (this.state.selectConfigID != "") {
                this.saveConfig(
                  this.state.selectConfigID,
                  this.jsonEditor.get()
                );
              } else {
                alert("please choose config");
              }
            }}
          >
            Save
          </Button>
          <Button
            color="inherit"
            className={classes.button}
            onClick={e => {
              if (this.state.selectConfigID != "") {
                this.deleteConfig(this.state.selectConfigID);
              } else {
                alert("please choose config");
              }
            }}
          >
            Delete
          </Button>
          <FormControl className={classes.formControlSpan} />
          <Button
            color="secondary"
            className={classes.button}
            onClick={e => {
              this.newConfigDialog.handleClickOpen();
            }}
          >
            Create
          </Button>
          <NewConfigDialog
            onRef={this.setNewConfigDialog}
            refresh={this.fetchConfigID}
          />
        </form>
        <div>
          <Editor
            ref={this.setJsonEditorRef}
            mode="code"
            name="content"
            //schema={yourSchema}
            value={{}}
          />
        </div>
      </div>
    );
  }
}

ConfigEditor.propTypes = {
  classes: PropTypes.object.isRequired
};

export default withStyles(styles)(ConfigEditor);
