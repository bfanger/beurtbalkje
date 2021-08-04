#!/usr/bin/env node
const { execSync, execFile } = require("child_process");

const pos = process.argv.indexOf(__filename);
let args = process.argv;
for (let i = 0; i <= pos; i++) {
  args.shift();
}
let bin = "linux";
try {
  const platform = execSync("uname", { encoding: "utf-8" });
  if (/^Darwin/.test(platform)) {
    bin = "beurtbalkje-macos";
  } else {
    bin = "beurtbalkje-linux";
  }
} catch {
  bin = "beurtbalkje.exe";
}
const child = execFile(__dirname + "/bin/" + bin, args, (err) => {
  if (err.code) {
    process.exit(err.code);
  }
});
child.stdout.on("data", function (data) {
  process.stdout.write(data);
});
child.stderr.on("data", function (data) {
  process.stderr.write(data);
});
