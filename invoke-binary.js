const os = require("os");
const childProcess = require("child_process");

function chooseBinary() {
    const platform = os.platform().startsWith("win") ? "windows" : os.platform();
    let arch = os.arch();

    if (arch === "x64") {
        arch = "amd64";
    } else if (arch === "x32") {
        arch = "386";
    }

    return `${__dirname}/bin/main-${platform}-${arch}`;
}

function run() {
    const binary = chooseBinary();

    childProcess.spawnSync(binary, {stdio: "inherit"});
}

run();