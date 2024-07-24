import * as fs from 'fs'
import { exec as cbExec } from 'child_process'
import * as path from 'path'
import { promisify } from 'util'

const app = process && process.type === 'renderer' ? require('@electron/remote').app : require('electron').app
const AiAdmin = app.isPackaged ? path.join(process.resourcesPath, 'AiAdmin') : path.resolve(process.cwd(), '..', 'AiAdmin')
const exec = promisify(cbExec)
const symlinkPath = '/usr/local/bin/AiAdmin'

export function installed() {
  return fs.existsSync(symlinkPath) && fs.readlinkSync(symlinkPath) === AiAdmin
}

export async function install() {
  const command = `do shell script "mkdir -p ${path.dirname(
    symlinkPath
  )} && ln -F -s \\"${AiAdmin}\\" \\"${symlinkPath}\\"" with administrator privileges`

  await exec(`osascript -e '${command}'`)
}
