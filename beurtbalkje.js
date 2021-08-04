#!/usr/bin/env node
const { execSync, exec } = require("child_process");

function strippedArgs() {
  let args = process.argv;
  for (let i = 0; i < 2; i++) {
    if (/\/beurtbalkje(\.js)?]$/.test(args[0])) {
      args.shift();
      break;
    }
    args.shift();
  }
  return args;
}

function detectBin() {
  let bin = "linux";
  try {
    const platform = execSync("uname", { encoding: "utf-8", stdio: "pipe" });
    if (/^Darwin/.test(platform)) {
      bin = "beurtbalkje-macos";
    } else {
      bin = "beurtbalkje-linux";
    }
  } catch {
    bin = "beurtbalkje.exe";
  }
  return bin;
}

const cli = __dirname + "/bin/" + detectBin() + " " + strippedArgs().join(" ");

const child = exec(cli, (err) => {
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
