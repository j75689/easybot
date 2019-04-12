import React from "react";
import Layout from "../layout/Layout";
import ConfigEditor from "../component/ConfigEditor";

class MessageHandlerConfig extends React.Component {
  state = {};
  render() {
    return (
      <Layout PageName="Config">
        <ConfigEditor />
      </Layout>
    );
  }
}

export default MessageHandlerConfig;
