import Config from "./config";
import axios from "axios";
import querystring from "querystring";

const service = axios.create({
  crossDomain: Config.isDeveloper,
  baseURL: Config.isDeveloper ? "http://localhost:8801/" : Config.Basehref,
  timeout: 5000
});

const api = {
  async Login(user, pass) {
    try {
      let form = new FormData();
      form.append("user", user);
      form.append("pass", pass);

      let res = await service.post("login", form, {
        headers: { "Content-Type": "multipart/form-data" }
      });

      return new Promise(resolve => {
        if (res.code === 0) {
          resolve(res);
        } else {
          resolve(res);
        }
      });
    } catch (err) {
      console.log(err);
    }
  },
  async BatchDeleteServiceAccount(accounts) {
    try {
      let res = await service.delete("/role/account", {
        data: accounts
      });
      return new Promise(resolve => {
        if (res.code === 0) {
          resolve(res);
        } else {
          resolve(res);
        }
      });
    } catch (err) {
      console.log(err);
    }
  },
  async GetAllServiceAccount() {
    try {
      let res = await service.get("/role/account");
      return new Promise(resolve => {
        if (res.code === 0) {
          resolve(res);
        } else {
          resolve(res);
        }
      });
    } catch (err) {
      console.log(err);
    }
  },
  async CreateServiceAccount(name, data) {
    try {
      let res = await service.post(
        "/role/account/" + name,
        querystring.stringify(data)
      );
      return new Promise(resolve => {
        if (res.code === 0) {
          resolve(res);
        } else {
          resolve(res);
        }
      });
    } catch (err) {
      console.log(err);
    }
  },
  async GetServiceAccount(name) {
    try {
      let res = await service.get("/role/account/" + name);
      return new Promise(resolve => {
        if (res.code === 0) {
          resolve(res);
        } else {
          resolve(res);
        }
      });
    } catch (err) {
      console.log(err);
    }
  },
  async SaveServiceAccount(name, data) {
    try {
      let res = await service.put(
        "/role/account/" + name,
        querystring.stringify(data)
      );
      return new Promise(resolve => {
        if (res.code === 0) {
          resolve(res);
        } else {
          resolve(res);
        }
      });
    } catch (err) {
      console.log(err);
    }
  },
  async RefreshServiceAccountToken(name) {
    try {
      let res = await service.post(`/role/account/${name}/refresh`);
      return new Promise(resolve => {
        if (res.code === 0) {
          resolve(res);
        } else {
          resolve(res);
        }
      });
    } catch (err) {
      console.log(err);
    }
  },
  async GetScopeTags() {
    try {
      let res = await service.get(`/role/scope`);
      return new Promise(resolve => {
        if (res.code === 0) {
          resolve(res);
        } else {
          resolve(res);
        }
      });
    } catch (err) {
      console.log(err);
    }
  },
  async GetAllConfigIDs() {
    try {
      let res = await service.get(`/handler/config`);
      return new Promise(resolve => {
        if (res.code === 0) {
          resolve(res);
        } else {
          resolve(res);
        }
      });
    } catch (err) {
      console.log(err);
    }
  },
  async GetHandlerConfig(configID) {
    try {
      let res = await service.get(`/handler/config/${configID}`);
      return new Promise(resolve => {
        if (res.code === 0) {
          resolve(res);
        } else {
          resolve(res);
        }
      });
    } catch (err) {
      console.log(err);
    }
  },
  async CreateHandlerConfig(configID, data) {
    try {
      let res = await service.post(`/handler/config/${configID}`, data);
      return new Promise(resolve => {
        if (res.code === 0) {
          resolve(res);
        } else {
          resolve(res);
        }
      });
    } catch (err) {
      console.log(err);
    }
  },
  async SaveHandlerConfig(configID, data) {
    try {
      let res = await service.put(`/handler/config/${configID}`, data);
      return new Promise(resolve => {
        if (res.code === 0) {
          resolve(res);
        } else {
          resolve(res);
        }
      });
    } catch (err) {
      console.log(err);
    }
  },
  async DeleteHandlerConfig(configID) {
    try {
      let res = await service.delete(`/handler/config/${configID}`);
      return new Promise(resolve => {
        if (res.code === 0) {
          resolve(res);
        } else {
          resolve(res);
        }
      });
    } catch (err) {
      console.log(err);
    }
  }
};

export default api;
