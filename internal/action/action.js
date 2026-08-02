// SPDX-License-Identifier: Apache-2.0

import { spawnSync } from "node:child_process";
import * as console from "node:console";
import { createHash } from "node:crypto";
import { appendFile, mkdtemp, readFile } from "node:fs/promises";
import { arch, platform, tmpdir } from "node:os";
import { join } from "node:path";
import { env, exit } from "node:process";

// --- Constants ---------------------------------------------------------------
const CHECKSUM_FILE = "checksums-sha512.txt";
const GHASUM_REPOSITORY = "chains-project/ghasum";

// --- Context -----------------------------------------------------------------
const ARCH = arch().toLowerCase();
const OS = platform().toLowerCase();

const JOB = env.GITHUB_JOB;
const SHA = env.GITHUB_WORKFLOW_SHA;
const OWNER = env.GITHUB_REPOSITORY.split("/").at(0);
const PROJECT = env.GITHUB_REPOSITORY.split("/").at(1);
const WORKFLOW = env.INPUT_WORKFLOW.split(/[/@]/g).slice(2,5).join("/");
const GITHUB_TOKEN = env.INPUT_TOKEN;

let cache;
switch (OS) {
case "darwin": cache = "/Users/runner/work/_actions"; break;
case "linux":  cache = "/home/runner/work/_actions";  break;
case "win32":  cache = "D:\\a\\_actions";             break;
}

let archive;
switch (`${OS}-${ARCH}`) {
case "darwin-arm64": archive = "ghasum_darwin_arm64.tar.gz"; break;
case "darwin-x64":   archive = "ghasum_darwin_amd64.tar.gz"; break;
case "linux-arm64":  archive = "ghasum_linux_arm64.tar.gz";  break;
case "linux-x64":    archive = "ghasum_linux_amd64.tar.gz";  break;
case "win32-arm64":  archive = "ghasum_windows_arm64.zip";   break;
case "win32-x64":    archive = "ghasum_windows_amd64.zip";   break;
}

let executable;
switch (OS) {
case "darwin": executable = "ghasum";     break;
case "linux":  executable = "ghasum";     break;
case "win32":  executable = "ghasum.exe"; break;
}

// --- Inputs ------------------------------------------------------------------
const CHECKSUM = env.INPUT_CHECKSUM.replace(/^sha256:/, "");
const MODE = env.INPUT_MODE;
const VERSION = env.INPUT_VERSION;

// --- Main --------------------------------------------------------------------
try {
	if (MODE !== "install" && MODE !== "verify") {
		throw new Error(`mode must be 'install' or 'verify', got: ${MODE}`);
	}

	const cwd = await mkdtemp(join(tmpdir(), 'ghasum-'));
	exec(["gh", "release", "download", VERSION, "--repo", GHASUM_REPOSITORY, "--pattern", CHECKSUM_FILE], { cwd });
	await sum(cwd, CHECKSUM_FILE, 256, CHECKSUM, null);
	exec(["gh", "release", "download", VERSION, "--repo", GHASUM_REPOSITORY, "--pattern", archive], { cwd });
	await sum(cwd, archive, 512, null, CHECKSUM_FILE);
	exec(["tar", "-xf", archive], { cwd });

	switch (MODE) {
	case "install":
		await appendFile(env.GITHUB_PATH, cwd);
		break;
	case "verify":
		exec(
			[join(cwd, executable), "verify", "-cache", cache, "-no-evict", "-offline", `${WORKFLOW}:${JOB}`],
			{ cwd: join(cache, OWNER, PROJECT, SHA) },
		);
		break;
	}

	exit(0);
} catch (error) {
	console.error(`::error::${error.message}`);
	nuke();
	exit(1);
}

// --- Functions ---------------------------------------------------------------
function exec(command, opts) {
	console.info("$", command.join(" "));

	const cmd = command[0];
	const args = command.slice(1, command.length);
	const { status } = spawnSync(cmd, args,  {
		env: { ...env, GITHUB_TOKEN },
		stdio: "inherit" ,
		...opts
	});

	if (status !== 0) {
		throw new Error("Command failed");
	}
}

async function sum(wd, target, algo, sum, sumfile) {
	const data = await readFile(join(wd, target));
	const hasher = createHash(`sha${algo}`);
	hasher.update(data);

	const got = hasher.digest("hex");
	let want;
	if (sum) {
		want = sum;
	} else {
		const sums = await readFile(join(wd, sumfile), { encoding: "utf8" });
		const line = sums.split(/\r?\n/).find(line => line.endsWith(target));
		want = line.split(" ").at(0);
	}

	if (got !== want) {
		throw new Error(`checksum mismatch for ${target} ('${got}' != '${want}')`);
	}
}

function nuke() {
	exec(["rm", "-rf", cache]);
}
