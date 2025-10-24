import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';
import { enhancedImages } from '@sveltejs/enhanced-img';
import { exec } from 'node:child_process';

/** @type {import("vite").Plugin}*/
const compileWASM = {
	name: 'compileWASM',
	handleHotUpdate({ file }) {
		if (file.endsWith('static/main.wasm')) return;
		console.log(`File changed ${file}`);
		exec(
			'cd ../wisim-backend/wasm && GOOS=js GOARCH=wasm go build -o main.wasm main.go && cd ../../ && cp wisim-backend/wasm/main.wasm wisim-nouveau/static/main.wasm && echo "main.wasm compiled"',
			(err, stdout, stderr) => {
				if (err) {
					console.error(`Error executing script: ${err}`);
					return;
				}
				console.info(`Script output: ${stdout}`);
				if (stderr) {
					console.error(`Script error output:{stderr}`);
				}
			}
		);
	}
};

export default defineConfig({
	plugins: [enhancedImages(), sveltekit(), compileWASM]
});
