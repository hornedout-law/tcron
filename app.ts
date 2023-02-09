import fs from "fs";
import path from "path";
import { spawn } from "child_process";

type Options = {
  timeout: string;
};

let args = process.argv
  .map((v, i) => {
    return i > 1 ? v : undefined;
  })
  .filter((v) => v != undefined) as string[];

function parseArguments(opts: Array<string>, args: Array<string>) {
  let options: { [key: string]: string | true } = {};
  let leftArgs: string[] = args;
  for (let option of opts) {
    let formats = option.split(",").map((v) => {
      let splited = v.split(/<.+>/i);
      if (splited.length > 1) {
        return splited[0].trim();
      } else {
        return v;
      }
    });
    let optionIndex =
      args.indexOf(formats[0]) > -1
        ? args.indexOf(formats[0])
        : args.indexOf(formats[1]) > -1
        ? args.indexOf(formats[1])
        : -1;
    if (optionIndex != -1) {
      options[formats[0].split("--")[1]] =
        optionIndex + 1 < args.length ? args[optionIndex + 1] : true;
      leftArgs = leftArgs
        .map((v, i) => {
          if (i == optionIndex || i == optionIndex + 1) return undefined;
          else return v;
        })
        .filter((v) => v != undefined) as string[];
    }
  }
  return { options, leftArgs };
}

function getTimeout(arg: string | undefined) {
    if(arg==undefined){
      let now = new Date()
        return `${now.getMinutes()} ${now.getHours()} * * *`
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
  let now = Date.now()
  let timeout =  matches
    .map((v) => Object.assign({}, v))
    .reduce(
      // @ts-ignore
      (acc, cur) => ({ ...acc, [time.get(cur.time) as string]: parseInt(cur.duration) }),
      {}
    );
  let timeinMili = 0
  for (let k in timeout) timeinMili+=(timeToSecs.get(k)as number)*parseInt(timeout[k])
  let predicteTime = new Date(now+timeinMili)
  
  return `${predicteTime.getMinutes()} ${predicteTime.getHours()} ${predicteTime.getDate()} ${predicteTime.getMonth()} ${predicteTime.getDay()}`
}

export function createDirWithDelete() {
  let parsed = parseArguments(["--timeout <duration>,-t"], args);
  let dirPath = path.resolve(process.cwd(), parsed.leftArgs[0]);
  
  let timeString = getTimeout(parsed.options?.timeout as string)
  
  fs.mkdirSync(dirPath);
  let cronScript = spawn("bash", [path.resolve(__dirname, "../cron_script.sh"), timeString, dirPath])
  cronScript.stdout.on("data", (chunck)=>console.log(chunck.toString()))
  cronScript.stderr.on("data", (chunck)=>console.log(chunck.toString()))
}
