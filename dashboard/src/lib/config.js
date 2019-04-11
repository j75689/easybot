var base = document.getElementsByTagName("base")[0].href;
base = base.substring(base.indexOf("//") + 2, base.length);
if (base.indexOf("/") > -1) {
  base = base.substring(base.indexOf("/"), base.length);
}
if ("development" == process.env.NODE_ENV) {
  base = base.replace("%7B%7B.contextPath%7D%7D", "");
}
if (!base.endsWith("/")) {
  base = base + "/";
}

const Config = {
  Basehref: base,
  isDeveloper: "development" == process.env.NODE_ENV
};

export default Config;
