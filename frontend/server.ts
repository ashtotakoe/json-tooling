import bun from "bun";
import fs from "fs/promises";
import { ServerResponse } from "http";

console.log("listening on http://localhost:8080/");
bun.serve({
  port: 8080,
  routes: {
    "/*": async (req) => {
      const allowed = ["style.css", "main.js", "index.html"];
      const file =
        req.url.slice("http://localhost:8080/".length).split("/")[0] ||
        "index.html";

      console.log(file);
      if (!allowed.includes(file)) {
        return new Response("", { status: 404 });
      }

      const resp = new Response(await fs.readFile(`./${file}`), {
        status: 200,
      });

      if (file === "index.html")
        resp.headers.append("content-type", "text/html");
      if (file === "style.css") resp.headers.append("content-type", "text/css");

      return resp;
    },
  },
});
