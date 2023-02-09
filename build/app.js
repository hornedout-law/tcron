"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.createDirWithDelete = void 0;
const fs_1 = __importDefault(require("fs"));
const path_1 = __importDefault(require("path"));
const child_process_1 = require("child_process");
let args = process.argv
    .map((v, i) => {
    return i > 1 ? v : undefined;
})
    .filter((v) => v != undefined);
function parseArguments(opts, args) {
    let options = {};
    let leftArgs = args;
    for (let option of opts) {
        let formats = option.split(",").map((v) => {
            let splited = v.split(/<.+>/i);
            if (splited.length > 1) {
                return splited[0].trim();
            }
            else {
                return v;
            }
        });
        let optionIndex = args.indexOf(formats[0]) > -1
            ? args.indexOf(formats[0])
            : args.indexOf(formats[1]) > -1
                ? args.indexOf(formats[1])
                : -1;
        if (optionIndex != -1) {
            options[formats[0].split("--")[1]] =
                optionIndex + 1 < args.length ? args[optionIndex + 1] : true;
            leftArgs = leftArgs
                .map((v, i) => {
                if (i == optionIndex || i == optionIndex + 1)
                    return undefined;
                else
                    return v;
            })
                .filter((v) => v != undefined);
        }
    }
    return { options, leftArgs };
}
function getTimeout(arg) {
    if (arg == undefined) {
        let now = new Date();
        return `${now.getMinutes()} ${now.getHours()} * * *`;
    }
    let regex = /((?<duration>\d){1,2}(?<time>[hwmd]){1})/giy;
    let res;
    let matches = [];
    let time = new Map([
        ["w", "week"],
        ["h", "hour"],
        ["d", "day"],
        ["m", "month"],
    ]);
    let timeToSecs = new Map([
        ["week", 604800000],
        ["hour", 3600000],
        ["day", 86400000],
        ["month", 2.628e+9],
    ]);
    while ((res = regex.exec(arg)) != null) {
        matches.push(res.groups);
    }
    let now = Date.now();
    let timeout = matches
        .map((v) => Object.assign({}, v))
        .reduce(
    // @ts-ignore
    (acc, cur) => (Object.assign(Object.assign({}, acc), { [time.get(cur.time)]: parseInt(cur.duration) })), {});
    let timeinMili = 0;
    for (let k in timeout)
        timeinMili += timeToSecs.get(k) * parseInt(timeout[k]);
    let predicteTime = new Date(now + timeinMili);
    return `${predicteTime.getMinutes()} ${predicteTime.getHours()} ${predicteTime.getDate()} ${predicteTime.getMonth()} ${predicteTime.getDay()}`;
}
function createDirWithDelete() {
    var _a;
    let parsed = parseArguments(["--timeout <duration>,-t"], args);
    let dirPath = path_1.default.resolve(process.cwd(), parsed.leftArgs[0]);
    let timeString = getTimeout((_a = parsed.options) === null || _a === void 0 ? void 0 : _a.timeout);
    fs_1.default.mkdirSync(dirPath);
    let cronScript = (0, child_process_1.spawn)("bash", [path_1.default.resolve(__dirname, "../cron_script.sh"), timeString, dirPath]);
    cronScript.stdout.on("data", (chunck) => console.log(chunck.toString()));
    cronScript.stderr.on("data", (chunck) => console.log(chunck.toString()));
}
exports.createDirWithDelete = createDirWithDelete;
