import React from "react";
import PropTypes from "prop-types";
import classNames from "classnames";
import { withStyles } from "@material-ui/core/styles";
import TableCell from "@material-ui/core/TableCell";
import TableSortLabel from "@material-ui/core/TableSortLabel";
import Paper from "@material-ui/core/Paper";
import Toolbar from "@material-ui/core/Toolbar";
import Tooltip from "@material-ui/core/Tooltip";
import { lighten } from "@material-ui/core/styles/colorManipulator";
import Typography from "@material-ui/core/Typography";
import IconButton from "@material-ui/core/IconButton";
import DeleteIcon from "@material-ui/icons/Delete";
import PlayListAddIcon from "@material-ui/icons/PlaylistAdd";
import Chip from "@material-ui/core/Chip";
import Avatar from "@material-ui/core/Avatar";
import FaceIcon from "@material-ui/icons/Face";
import api from "../lib/api";
import NewIptableDialog from "./NewIptableDialog";
import EditIptableDialog from "./EditIptableDialog";

import { AutoSizer, Column, SortDirection, Table } from "react-virtualized";

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

class IpTableToolbar extends React.Component {
  constructor(props) {
    super(props);
  }

  render() {
    const { numSelected, classes } = this.props;

    return (
      <Toolbar>
        <div className={classes.title}>
          <Typography variant="h6" id="tableTitle">
            Iptable
          </Typography>
        </div>
        <div className={classes.spacer} />
        <div className={classes.actions}>
          <Tooltip title="Add">
            <IconButton
              aria-label="Add"
              onClick={event => {
                this.props.open();
              }}
            >
              <PlayListAddIcon />
            </IconButton>
          </Tooltip>
        </div>
      </Toolbar>
    );
  }
}

IpTableToolbar.propTypes = {
  classes: PropTypes.object.isRequired,
  numSelected: PropTypes.number.isRequired
};

IpTableToolbar = withStyles(toolbarStyles)(IpTableToolbar);

const styles = theme => ({
  table: {
    fontFamily: theme.typography.fontFamily
  },
  flexContainer: {
    display: "flex",
    alignItems: "center",
    boxSizing: "border-box"
  },
  tableRow: {
    cursor: "pointer"
  },
  tableRowHover: {
    "&:hover": {
      backgroundColor: theme.palette.grey[200]
    }
  },
  tableCell: {
    flex: 1
  },
  noClick: {
    cursor: "initial"
  }
});

class MuiVirtualizedTable extends React.PureComponent {
  getRowClassName = ({ index }) => {
    const { classes, rowClassName, onRowClick } = this.props;

    return classNames(classes.tableRow, classes.flexContainer, rowClassName, {
      [classes.tableRowHover]: index !== -1 && onRowClick != null
    });
  };

  cellRenderer = ({ cellData, columnIndex = null }) => {
    const { columns, classes, rowHeight, onRowClick } = this.props;
    return (
      <TableCell
        component="div"
        className={classNames(classes.tableCell, classes.flexContainer, {
          [classes.noClick]: onRowClick == null
        })}
        variant="body"
        style={{ height: rowHeight }}
        align={
          (columnIndex != null && columns[columnIndex].numeric) || false
            ? "right"
            : "left"
        }
      >
        {cellData}
      </TableCell>
    );
  };

  headerRenderer = ({ label, columnIndex, dataKey, sortBy, sortDirection }) => {
    const { headerHeight, columns, classes, sort } = this.props;
    const direction = {
      [SortDirection.ASC]: "asc",
      [SortDirection.DESC]: "desc"
    };

    const inner =
      !columns[columnIndex].disableSort && sort != null ? (
        <TableSortLabel
          active={dataKey === sortBy}
          direction={direction[sortDirection]}
        >
          {label}
        </TableSortLabel>
      ) : (
        label
      );

    return (
      <>
        <TableCell
          component="div"
          className={classNames(
            classes.tableCell,
            classes.flexContainer,
            classes.noClick
          )}
          variant="head"
          style={{ height: headerHeight }}
          align={columns[columnIndex].numeric || false ? "right" : "left"}
        >
          {inner}
        </TableCell>
      </>
    );
  };

  render() {
    const { classes, columns, ...tableProps } = this.props;
    return (
      <>
        <AutoSizer>
          {({ height, width }) => (
            <>
              <Table
                className={classes.table}
                height={height - 66}
                width={width}
                {...tableProps}
                rowClassName={this.getRowClassName}
              >
                {columns.map(
                  (
                    {
                      cellContentRenderer = null,
                      className,
                      dataKey,
                      ...other
                    },
                    index
                  ) => {
                    let renderer;
                    if (cellContentRenderer != null) {
                      renderer = cellRendererProps =>
                        this.cellRenderer({
                          cellData: cellContentRenderer(cellRendererProps),
                          columnIndex: index
                        });
                    } else {
                      renderer = this.cellRenderer;
                    }

                    return (
                      <Column
                        key={dataKey}
                        headerRenderer={headerProps =>
                          this.headerRenderer({
                            ...headerProps,
                            columnIndex: index
                          })
                        }
                        className={classNames(classes.flexContainer, className)}
                        cellRenderer={renderer}
                        dataKey={dataKey}
                        {...other}
                      />
                    );
                  }
                )}
              </Table>
            </>
          )}
        </AutoSizer>
      </>
    );
  }
}

MuiVirtualizedTable.propTypes = {
  classes: PropTypes.object.isRequired,
  columns: PropTypes.arrayOf(
    PropTypes.shape({
      cellContentRenderer: PropTypes.func,
      dataKey: PropTypes.string.isRequired,
      width: PropTypes.number.isRequired
    })
  ).isRequired,
  headerHeight: PropTypes.number,
  onRowClick: PropTypes.func,
  rowClassName: PropTypes.string,
  rowHeight: PropTypes.oneOfType([PropTypes.number, PropTypes.func]),
  sort: PropTypes.func
};

MuiVirtualizedTable.defaultProps = {
  headerHeight: 56,
  rowHeight: 56
};

const WrappedVirtualizedTable = withStyles(styles)(MuiVirtualizedTable);

function createData(id, ips, type, scope, option) {
  return { id, ips, type, scope, option };
}

class IpTable extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      rows: []
    };
    this.fetchIptables = this.fetchIptables.bind(this);
    this.handleDelete = this.handleDelete.bind(this);
  }

  async fetchIptables() {
    let resp = await api.GetIptables();
    if (resp) {
      if (resp.data) {
        let rows = [];
        resp.data.map(item => {
          let ips = [];
          item.ip.map(range => {
            if (ips.length > 0) ips.push(<>&nbsp;</>);
            ips.push(<Chip label={range} />);
          });
          let option = (
            <IconButton aria-label="Delete">
              <DeleteIcon
                fontSize="small"
                onClick={e => {
                  if (window.confirm("Sure?")) {
                    this.handleDelete(item.id);
                  }
                  e.preventDefault();
                }}
              />
            </IconButton>
          );
          rows.push(createData(item.id, ips, item.type, item.scope, option));
        });
        this.setState({ rows: rows });
      }
    }
  }

  async handleDelete(id) {
    let resp = await api.DeleteIptable(id);
    if (resp) {
      if (resp.data.success) {
        alert("Deleted.");
        this.fetchIptables();
      } else {
        alert(resp.data.error);
      }
    } else {
      alert("Error!");
    }
  }

  componentDidMount() {
    this.fetchIptables();
  }

  setCreateDialog = ref => {
    this.createDialog = ref;
  };

  setEditDialog = ref => {
    this.editDialog = ref;
  };

  render() {
    return (
      <>
        <EditIptableDialog
          refresh={this.fetchIptables}
          onRef={this.setEditDialog}
        />
        <NewIptableDialog
          refresh={this.fetchIptables}
          onRef={this.setCreateDialog}
        />
        <Paper style={{ height: 400, width: "100%" }}>
          <IpTableToolbar
            open={() => {
              this.createDialog.handleClickOpen();
            }}
          />
          <WrappedVirtualizedTable
            rowCount={this.state.rows.length}
            rowGetter={({ index }) => this.state.rows[index]}
            onRowClick={e => {
              this.editDialog.handleClickOpen(e.rowData.id);
            }}
            columns={[
              {
                width: 200,
                flexGrow: 1.0,
                label: "IP",
                dataKey: "ips",
                numeric: false
              },
              {
                width: 120,
                label: "Type",
                dataKey: "type",
                numeric: false
              },
              {
                width: 120,
                label: "Scope",
                dataKey: "scope",
                numeric: false
              },
              {
                width: 120,
                label: "",
                dataKey: "option",
                numeric: false
              }
            ]}
          />
        </Paper>
      </>
    );
  }
}

export default IpTable;
