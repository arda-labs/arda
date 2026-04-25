import { createConfig } from '@nx/angular-rspack';
import baseWebpackConfig from './webpack.config';
import webpackMerge from 'webpack-merge';

export default async () => {
  const baseConfig = await createConfig(
    {
      options: {
        root: __dirname,

        outputPath: {
          base: '../../dist/apps/common',
        },
        index: './src/index.html',
        browser: './src/main.ts',
        tsConfig: './tsconfig.app.json',
        assets: [
          {
            glob: '**/*',
            input: './public',
          },
          {
            glob: '**/*',
            input: '../../node_modules/primeicons',
            output: 'primeicons',
          },
        ],
        styles: ['./src/styles.css'],
        devServer: {
          port: 3001,
          publicHost: 'http://localhost:3001',
        },
      },
    },
    {
      production: {
        options: {
          budgets: [],
          outputHashing: 'all',
          deployUrl: '/common/',
          devServer: {},
        },
      },

      development: {
        options: {
          optimization: false,
          vendorChunk: true,
          extractLicenses: false,
          sourceMap: true,
          namedChunks: true,
          devServer: {},
        },
      },
    },
  );
  return webpackMerge(baseConfig[0], baseWebpackConfig);
};
