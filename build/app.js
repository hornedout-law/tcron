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
    let regex = /((?<duration>\d{1,2})(?<time>[hwmdn]){1})/giy;
    let res;
    let matches = [];
    let time = {
        w: 604800000,
        h: 3600000,
        d: 86400000,
        m: 2.628e+9,
        n: 60000
    };
    while ((res = regex.exec(arg)) != null) {
        matches.push(res.groups);
    }
    let now = new Date();
    let timeout = matches
        .map((v) => Object.assign({}, v))
        .reduce(
    // @ts-ignore
    (acc, cur) => acc + cur.duration * time[cur.time], 0);
    let predicteTime = new Date(now.getTime() + timeout);
    let starOrValue = (method, df) => {
        if (df) {
            // @ts-ignore
            return predicteTime[method]();
        }
        // @ts-ignore
        return predicteTime[method]() == now[method]() ? "*" : predicteTime[method]();
    };
    return `${starOrValue("getMinutes", true)} ${starOrValue("getHours", true)} ${starOrValue("getDate")} ${starOrValue("getMonth")} ${starOrValue("getDay")}`;
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
