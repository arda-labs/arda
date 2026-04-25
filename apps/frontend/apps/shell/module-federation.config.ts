import { ModuleFederationConfig } from '@nx/module-federation';

const config: ModuleFederationConfig = {
  name: 'shell',
  remotes: [],
  shared: (libraryName, sharedConfig) => {
    if (libraryName.startsWith('rxjs') || libraryName === 'tslib') {
      return { ...sharedConfig, eager: true };
    }
    return sharedConfig;
  },
};

/**
 * Nx requires a default export of the config to allow correct resolution of the module federation graph.
 **/
export default config;
