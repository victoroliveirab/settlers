import typescript from '@rollup/plugin-typescript';

/**
 * @type {import('rollup').RollupOptions}
 */
const config = {
  input: './src/main.ts',
  output: {
    file: 'app.js',
    format: 'iife',
  },
  watch: {
    include: ['src/**/*.ts', 'index.html'],
  },
  plugins: [typescript()]
};
export default config;
