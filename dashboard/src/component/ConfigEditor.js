import React from "react";
import ReactDOM from "react-dom";
import PropTypes from "prop-types";
import { withStyles } from "@material-ui/core/styles";
import Input from "@material-ui/core/Input";
import OutlinedInput from "@material-ui/core/OutlinedInput";
import FilledInput from "@material-ui/core/FilledInput";
import InputLabel from "@material-ui/core/InputLabel";
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

const styles = theme => ({
  root: {
    display: "flex",
    flexWrap: "wrap"
  },
  formControl: {
    margin: theme.spacing.unit,
    minWidth: 150
  },
  selectEmpty: {
    marginTop: theme.spacing.unit * 2
  }
});

class ConfigEditor extends React.Component {
  state = {
    age: "",
    name: "hai",
    labelWidth: 0
  };

  componentDidMount() {}

  handleChange = event => {
    this.setState({ [event.target.name]: event.target.value });
  };

  render() {
    const { classes } = this.props;

    return (
      <div>
        <form className={classes.root} autoComplete="off">
          <FormControl className={classes.formControl}>
            <Select
              value={this.state.age}
              onChange={this.handleChange}
              name="event"
              displayEmpty
              className={classes.selectEmpty}
            >
              <MenuItem value="" disabled>
                Choose
              </MenuItem>
            </Select>
            <FormHelperText>Event </FormHelperText>
          </FormControl>
          <FormControl className={classes.formControl}>
            <Select
              value={this.state.age}
              onChange={this.handleChange}
              name="configID"
              displayEmpty
              className={classes.selectEmpty}
            >
              <MenuItem value="" disabled>
                Choose
              </MenuItem>
              <MenuItem value={10}>Ten</MenuItem>
              <MenuItem value={20}>Twenty</MenuItem>
              <MenuItem value={30}>Thirty</MenuItem>
            </Select>
            <FormHelperText>Config </FormHelperText>
          </FormControl>
        </form>
        <div>
          <Editor
            mode="code"
            name="content"
            //schema={yourSchema}

            value={{
              id: "DefaultMessageExample",
              eventType: "message",
              defaultValues: {
                repo: "https://github.com/j75689/easybot"
              },
              stage: [
                {
                  type: "reply",
                  value: {
                    type: "text",
                    text:
                      "If you have any questions, you can open a new question on the GitHub board ${repo}/issues).\nWe will help you as soon as possible."
                  }
                }
              ]
            }}
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
