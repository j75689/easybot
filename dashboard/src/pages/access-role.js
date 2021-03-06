import React from "react";
import Layout from "../layout/Layout";
import RolebindTable from "../component/RolebindTable";
import IpTable from "../component/Iptable";
class AccessRole extends React.Component {
  state = {};
  render() {
    return (
      <Layout PageName="AccessRole">
        <RolebindTable {...this.props} />
        <div style={{ margin: "50px" }} />
        <IpTable {...this.props} />
      </Layout>
    );
  }
}

export default AccessRole;
