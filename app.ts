import fs from "fs";
import path from "path";
import { spawn } from "child_process";

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

function getTimeout(arg: string|undefined) {
  if (arg == undefined) {
      let now = new Date();
      return `${now.getMinutes()} ${now.getHours()} * * *`;
  }
  let regex = /((?<duration>\d{1,2})(?<time>[hwmdn]){1})/giy;
  let res;
  let matches = [];

  let time= {
    w:604800000,
    h: 3600000,
    d:86400000,
    m: 2.628e+9,
    n: 60000
  }
  
  while ((res = regex.exec(arg)) != null) {
      matches.push(res.groups);
  }
  let now = new Date();
  let timeout = matches
      .map((v) => Object.assign({}, v))
      .reduce(
  // @ts-ignore
  (acc, cur) => acc+cur.duration*time[cur.time], 0);
  
  let predicteTime = new Date(now.getTime() + timeout);
  let starOrValue = (method:string, df?:boolean)=>{
    if(df){
      // @ts-ignore
        return predicteTime[method]()
    }
    // @ts-ignore
    return predicteTime[method]()==now[method]()?"*":predicteTime[method]()
  }
  
  return `${starOrValue("getMinutes", true)} ${starOrValue("getHours", true)} ${starOrValue("getDate")} ${starOrValue("getMonth")} ${starOrValue("getDay")}`;
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
