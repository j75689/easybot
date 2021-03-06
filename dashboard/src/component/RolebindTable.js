import React from "react";
import classNames from "classnames";
import PropTypes from "prop-types";
import { withStyles } from "@material-ui/core/styles";
import Table from "@material-ui/core/Table";
import TableBody from "@material-ui/core/TableBody";
import TableCell from "@material-ui/core/TableCell";
import TableHead from "@material-ui/core/TableHead";
import TablePagination from "@material-ui/core/TablePagination";
import TableRow from "@material-ui/core/TableRow";
import TableSortLabel from "@material-ui/core/TableSortLabel";
import Toolbar from "@material-ui/core/Toolbar";
import Typography from "@material-ui/core/Typography";
import Paper from "@material-ui/core/Paper";
import Checkbox from "@material-ui/core/Checkbox";
import IconButton from "@material-ui/core/IconButton";
import Tooltip from "@material-ui/core/Tooltip";
import DeleteIcon from "@material-ui/icons/Delete";
import PlayListAddIcon from "@material-ui/icons/PlaylistAdd";
import Button from "@material-ui/core/Button";
import { lighten } from "@material-ui/core/styles/colorManipulator";
import NewRoleAccountDialog from "./NewRoleAccountDialog";
import EditRoleAccountDialog from "./EditRoleAccountDialog";
import api from "../lib/api";

const copyToClipboard = data => {
  const textField = document.createElement("textarea");
  textField.innerText = data;
  document.body.appendChild(textField);
  textField.select();
  document.execCommand("copy");
  textField.remove();
};

function createData(account, scope, generate, expire, token) {
  return {
    account: account,
    scope: scope,
    generate: generate,
    expire: expire,
    token: token
  };
}

function desc(a, b, orderBy) {
  if (b[orderBy] < a[orderBy]) {
    return -1;
  }
  if (b[orderBy] > a[orderBy]) {
    return 1;
  }
  return 0;
}

function stableSort(array, cmp) {
  const stabilizedThis = array.map((el, index) => [el, index]);
  stabilizedThis.sort((a, b) => {
    const order = cmp(a[0], b[0]);
    if (order !== 0) return order;
    return a[1] - b[1];
  });
  return stabilizedThis.map(el => el[0]);
}

function getSorting(order, orderBy) {
  return order === "desc"
    ? (a, b) => desc(a, b, orderBy)
    : (a, b) => -desc(a, b, orderBy);
}

const rows = [
  {
    id: "account",
    numeric: false,
    disablePadding: true,
    label: "Account Name"
  },
  { id: "scope", numeric: false, disablePadding: false, label: "Scope" },
  { id: "generate", numeric: false, disablePadding: false, label: "Generate" },
  { id: "expire", numeric: false, disablePadding: false, label: "Expired" },
  { id: "token", numeric: false, disablePadding: false, label: "Token" },
  { id: "editButton", numeric: false, disablePadding: false, label: "" }
];

class RolebindTableHead extends React.Component {
  createSortHandler = property => event => {
    this.props.onRequestSort(event, property);
  };

  render() {
    const {
      onSelectAllClick,
      order,
      orderBy,
      numSelected,
      rowCount
    } = this.props;

    return (
      <TableHead>
        <TableRow>
          <TableCell padding="checkbox">
            <Checkbox
              indeterminate={numSelected > 0 && numSelected < rowCount}
              checked={numSelected === rowCount}
              onChange={onSelectAllClick}
            />
          </TableCell>
          {rows.map(
            row => (
              <TableCell
                key={row.id}
                align={row.numeric ? "right" : "left"}
                padding={row.disablePadding ? "none" : "default"}
                sortDirection={orderBy === row.id ? order : false}
              >
                <Tooltip
                  title="Sort"
                  placement={row.numeric ? "bottom-end" : "bottom-start"}
                  enterDelay={300}
                >
                  <TableSortLabel
                    active={orderBy === row.id}
                    direction={order}
                    onClick={this.createSortHandler(row.id)}
                  >
                    {row.label}
                  </TableSortLabel>
                </Tooltip>
              </TableCell>
            ),
            this
          )}
        </TableRow>
      </TableHead>
    );
  }
}

RolebindTableHead.propTypes = {
  numSelected: PropTypes.number.isRequired,
  onRequestSort: PropTypes.func.isRequired,
  onSelectAllClick: PropTypes.func.isRequired,
  order: PropTypes.string.isRequired,
  orderBy: PropTypes.string.isRequired,
  rowCount: PropTypes.number.isRequired
};

const toolbarStyles = theme => ({
  root: {
    paddingRight: theme.spacing.unit
  },
  highlight:
    theme.palette.type === "light"
      ? {
          color: theme.palette.secondary.main,
          backgroundColor: lighten(theme.palette.secondary.light, 0.85)
        }
      : {
          color: theme.palette.text.primary,
          backgroundColor: theme.palette.secondary.dark
        },
  spacer: {
    flex: "1 1 100%"
  },
  actions: {
    color: theme.palette.text.secondary
  },
  title: {
    flex: "0 0 auto"
  }
});

class RolebindTableToolbar extends React.Component {
  constructor(props) {
    super(props);
  }

  render() {
    const { numSelected, classes } = this.props;

    return (
      <Toolbar
        className={classNames(classes.root, {
          [classes.highlight]: numSelected > 0
        })}
      >
        <div className={classes.title}>
          {numSelected > 0 ? (
            <Typography color="inherit" variant="subtitle1">
              {numSelected} selected
            </Typography>
          ) : (
            <Typography variant="h6" id="tableTitle">
              Service API Token
            </Typography>
          )}
        </div>
        <div className={classes.spacer} />
        <div className={classes.actions}>
          {numSelected > 0 ? (
            <Tooltip title="Delete">
              <IconButton
                aria-label="Delete"
                onClick={event => {
                  let ok = window.confirm(
                    `make sure delete [${this.props.selected}] ?`
                  );
                  if (ok) {
                    this.props.delete();
                  }
                }}
              >
                <DeleteIcon />
              </IconButton>
            </Tooltip>
          ) : (
            <NewRoleAccountDialog {...this.props} />
          )}
        </div>
      </Toolbar>
    );
  }
}

RolebindTableToolbar.propTypes = {
  classes: PropTypes.object.isRequired,
  numSelected: PropTypes.number.isRequired
};

RolebindTableToolbar = withStyles(toolbarStyles)(RolebindTableToolbar);

const styles = theme => ({
  root: {
    width: "100%",
    marginTop: theme.spacing.unit * 3
  },
  table: {
    minWidth: 1020
  },
  tableWrapper: {
    overflowX: "auto"
  },
  button: {
    margin: theme.spacing.unit
  }
});

class RolebindTable extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      order: "asc",
      orderBy: "calories",
      selected: [],
      data: [
        // createData(
        //   "mock",
        //   "all",
        //   "2019-04-12",
        //   "no expiry",
        //   "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE1NTQ5NzkxMzIsIm5iZiI6MTU1NDk3OTEzMiwiTmFtZSI6Imh5bGliIiwiRU1haWwiOiJrMDBAaHl3ZWIuY29tLnR3IiwiRG9tYWluIjoiaHl3ZWIuY29tLnR3IiwiUHJvdmlkZXIiOiJoeXdlYiIsIlNjb3BlIjoiYWxsIiwiQWN0aXZlIjowfQ.kxCQ9bhZuZHjtgbrCjh946Ynr7eRjQIKsYgGREKArPzgdt0DWrAJDMVVEHxk46H_6t7R73QrBrWlDrUUdxEa3Q"
        // )
      ],
      page: 0,
      rowsPerPage: 5
    };

    this.RefreshAccount = this.RefreshAccount.bind(this);
    this.DeleteAccount = this.DeleteAccount.bind(this);
  }

  componentDidMount() {
    this.RefreshAccount();
  }

  async RefreshAccount() {
    let resp = await api.GetAllServiceAccount();
    let accounts = [];
    if (resp) {
      resp.data.map(item => {
        let generate = new Date(item.generate * 1000).toISOString();
        let expired =
          item.expired > 0
            ? new Date(item.expired * 1000).toISOString()
            : "no expiry";
        accounts.push(
          createData(item.name, item.scope, generate, expired, item.token)
        );
      });

      this.setState({
        data: accounts
      });
    }
  }

  async DeleteAccount() {
    let resp = await api.BatchDeleteServiceAccount(this.state.selected);
    if (resp) {
      if (resp.data.success) {
        alert("Complete.");
        this.setState({ selected: [] });
      } else {
        alert(`Error! ${JSON.stringify(resp.data.error)}`);
      }
    }
    this.RefreshAccount();
  }

  handleRequestSort = (event, property) => {
    const orderBy = property;
    let order = "desc";

    if (this.state.orderBy === property && this.state.order === "desc") {
      order = "asc";
    }

    this.setState({ order, orderBy });
  };

  handleSelectAllClick = event => {
    if (event.target.checked) {
      this.setState(state => ({ selected: state.data.map(n => n.account) }));
      return;
    }
    this.setState({ selected: [] });
  };

  handleClick = (event, id) => {
    const { selected } = this.state;
    const selectedIndex = selected.indexOf(id);
    let newSelected = [];

    if (selectedIndex === -1) {
      newSelected = newSelected.concat(selected, id);
    } else if (selectedIndex === 0) {
      newSelected = newSelected.concat(selected.slice(1));
    } else if (selectedIndex === selected.length - 1) {
      newSelected = newSelected.concat(selected.slice(0, -1));
    } else if (selectedIndex > 0) {
      newSelected = newSelected.concat(
        selected.slice(0, selectedIndex),
        selected.slice(selectedIndex + 1)
      );
    }

    this.setState({ selected: newSelected });
  };

  handleChangePage = (event, page) => {
    this.setState({ page });
  };

  handleChangeRowsPerPage = event => {
    this.setState({ rowsPerPage: event.target.value });
  };

  isSelected = id => this.state.selected.indexOf(id) !== -1;

  render() {
    const { classes } = this.props;
    const { data, order, orderBy, selected, rowsPerPage, page } = this.state;
    const emptyRows =
      rowsPerPage - Math.min(rowsPerPage, data.length - page * rowsPerPage);

    return (
      <Paper className={classes.root}>
        <RolebindTableToolbar
          numSelected={selected.length}
          selected={selected}
          refresh={this.RefreshAccount}
          delete={this.DeleteAccount}
        />
        <div className={classes.tableWrapper}>
          <Table className={classes.table} aria-labelledby="tableTitle">
            <RolebindTableHead
              numSelected={selected.length}
              order={order}
              orderBy={orderBy}
              onSelectAllClick={this.handleSelectAllClick}
              onRequestSort={this.handleRequestSort}
              rowCount={data.length}
            />
            <TableBody>
              {stableSort(data, getSorting(order, orderBy))
                .slice(page * rowsPerPage, page * rowsPerPage + rowsPerPage)
                .map(n => {
                  const isSelected = this.isSelected(n.account);
                  return (
                    <TableRow
                      hover
                      role="checkbox"
                      aria-checked={isSelected}
                      tabIndex={-1}
                      key={n.account}
                      selected={isSelected}
                    >
                      <TableCell padding="checkbox">
                        <Checkbox
                          checked={isSelected}
                          onClick={event => this.handleClick(event, n.account)}
                        />
                      </TableCell>
                      <TableCell component="th" scope="row" padding="none">
                        {n.account}
                      </TableCell>
                      <TableCell style={{ maxWidth: "30px" }} align="left">
                        {n.scope}
                      </TableCell>
                      <TableCell style={{ maxWidth: "30px" }} align="left">
                        {n.generate}
                      </TableCell>
                      <TableCell style={{ maxWidth: "30px" }} align="left">
                        {n.expire}
                      </TableCell>
                      <TableCell style={{ maxWidth: "50px" }} align="left">
                        <div
                          onClick={e => {
                            copyToClipboard(n.token);
                            alert("copied.");
                          }}
                        >
                          {n.token.length > 30
                            ? `${n.token.substring(0, 30)}...`
                            : n.token}
                        </div>
                      </TableCell>
                      <TableCell align="right">
                        <EditRoleAccountDialog
                          account={n.account}
                          refresh={this.RefreshAccount}
                        />
                      </TableCell>
                    </TableRow>
                  );
                })}
              {emptyRows > 0 && (
                <TableRow style={{ height: 49 * emptyRows }}>
                  <TableCell colSpan={6} />
                </TableRow>
              )}
            </TableBody>
          </Table>
        </div>
        <TablePagination
          rowsPerPageOptions={[5, 10, 25]}
          component="div"
          count={data.length}
          rowsPerPage={rowsPerPage}
          page={page}
          backIconButtonProps={{
            "aria-label": "Previous Page"
          }}
          nextIconButtonProps={{
            "aria-label": "Next Page"
          }}
          onChangePage={this.handleChangePage}
          onChangeRowsPerPage={this.handleChangeRowsPerPage}
        />
      </Paper>
    );
  }
}

RolebindTable.propTypes = {
  classes: PropTypes.object.isRequired
};

export default withStyles(styles)(RolebindTable);
