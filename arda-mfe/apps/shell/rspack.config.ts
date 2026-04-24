import { createConfig } from '@nx/angular-rspack';
import baseWebpackConfig from './webpack.config';
import webpackMerge from 'webpack-merge';

export default async () => {
  const baseConfig = await createConfig(
    {
      options: {
        root: __dirname,

        outputPath: {
          base: '../../dist/apps/shell',
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
          port: 3000,
          publicHost: 'http://localhost:3000',
        },
      },
    },
    {
      production: {
        options: {
          budgets: [],
          outputHashing: 'all',
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
  const merged = webpackMerge(baseConfig[0], baseWebpackConfig);

  merged.devServer = {
    ...merged.devServer,
    proxy: [
      {
        context: ['/common'],
        target: 'http://localhost:3001',
        pathRewrite: { '^/common': '' },
      },
      {
        context: ['/api'],
        target: 'http://localhost:8000',
        pathRewrite: { '^/api': '' },
        secure: false,
      },
    ],
  };

  return merged;
};
