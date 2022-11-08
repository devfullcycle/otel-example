const { NodeTracerProvider } = require("@opentelemetry/node");
const {
  ConsoleSpanExporter,
  SimpleSpanProcessor,
} = require("@opentelemetry/tracing");
const { ZipkinExporter } = require("@opentelemetry/exporter-zipkin");
const provider = new NodeTracerProvider();
const consoleExporter = new ConsoleSpanExporter();
const spanProcessor = new SimpleSpanProcessor(consoleExporter);
provider.addSpanProcessor(spanProcessor);
provider.register();

const zipkinExporter = new ZipkinExporter({
  url: "http://localhost:9411/api/v2/spans",
  serviceName: "course-service",
});

const zipkinProcessor = new SimpleSpanProcessor(zipkinExporter);
provider.addSpanProcessor(zipkinProcessor);

const express = require("express");
const app = express();
const port = 3000;

app.get("/", async function (req, res) {
  res.type("json");
  await new Promise((resolve) => setTimeout(resolve, 50));
  const sqlite3 = require("sqlite3").verbose();
  let db = new sqlite3.Database("db.sqlite3", sqlite3.OPEN_READWRITE, (err) => {
    if (err) {
      console.error(err.message);
    }
    console.log("Connected to the  database.");
  });
  db.all(`SELECT * FROM courses`, [], (err, rows) => {
    res.status(200).json({ rows });
  });
});

app.listen(port, () => {
  console.log(`Listening at http://localhost:${port}`);
});
