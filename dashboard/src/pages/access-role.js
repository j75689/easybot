import React from "react";
import Layout from "../layout/Layout";
import RolebindTable from "../component/RolebindTable";

class AccessRole extends React.Component {
  state = {};
  render() {
    return (
      <Layout PageName="AccessRole">
        <RolebindTable />
      </Layout>
    );
  }
}

export default AccessRole;
