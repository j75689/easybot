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
  }
};

export default api;
