
// this file is generated — do not edit it


/// <reference types="@sveltejs/kit" />

/**
 * Environment variables [loaded by Vite](https://vitejs.dev/guide/env-and-mode.html#env-files) from `.env` files and `process.env`. Like [`$env/dynamic/private`](https://svelte.dev/docs/kit/$env-dynamic-private), this module cannot be imported into client-side code. This module only includes variables that _do not_ begin with [`config.kit.env.publicPrefix`](https://svelte.dev/docs/kit/configuration#env) _and do_ start with [`config.kit.env.privatePrefix`](https://svelte.dev/docs/kit/configuration#env) (if configured).
 * 
 * _Unlike_ [`$env/dynamic/private`](https://svelte.dev/docs/kit/$env-dynamic-private), the values exported from this module are statically injected into your bundle at build time, enabling optimisations like dead code elimination.
 * 
 * ```ts
 * import { API_KEY } from '$env/static/private';
 * ```
 * 
 * Note that all environment variables referenced in your code should be declared (for example in an `.env` file), even if they don't have a value until the app is deployed:
 * 
 * ```
 * MY_FEATURE_FLAG=""
 * ```
 * 
 * You can override `.env` values from the command line like so:
 * 
 * ```sh
 * MY_FEATURE_FLAG="enabled" npm run dev
 * ```
 */
declare module '$env/static/private' {
	export const VITE_API_BASE_URL: string;
	export const VITE_APP_TITLE: string;
	export const VITE_ENABLE_POLLING: string;
	export const VITE_POLLING_INTERVAL: string;
	export const LESS_TERMCAP_se: string;
	export const VSCODE_CWD: string;
	export const VSCODE_ESM_ENTRYPOINT: string;
	export const USER: string;
	export const LESS_TERMCAP_ue: string;
	export const VSCODE_NLS_CONFIG: string;
	export const GOTELEMETRY_GOPLS_CLIENT_START_TIME: string;
	export const npm_config_user_agent: string;
	export const VSCODE_WSL_EXT_LOCATION: string;
	export const VSCODE_HANDLES_UNCAUGHT_ERRORS: string;
	export const DB_PORT: string;
	export const GEMINI_YOLO_MODE: string;
	export const DEBUG: string;
	export const npm_node_execpath: string;
	export const JWT_SECRET_KEY: string;
	export const SHLVL: string;
	export const LD_LIBRARY_PATH: string;
	export const npm_config_noproxy: string;
	export const HOME: string;
	export const OLDPWD: string;
	export const NVM_BIN: string;
	export const VSCODE_IPC_HOOK_CLI: string;
	export const npm_package_json: string;
	export const NVM_INC: string;
	export const HOMEBREW_PREFIX: string;
	export const LESS_TERMCAP_so: string;
	export const DB_DATABASE: string;
	export const npm_config_userconfig: string;
	export const npm_config_local_prefix: string;
	export const NMAP_PRIVILEGED: string;
	export const TABLE_TOURNAMENTS: string;
	export const DBUS_SESSION_BUS_ADDRESS: string;
	export const WSL_DISTRO_NAME: string;
	export const COLOR: string;
	export const NVM_DIR: string;
	export const QT_QPA_PLATFORMTHEME: string;
	export const WAYLAND_DISPLAY: string;
	export const VSCODE_L10N_BUNDLE_LOCATION: string;
	export const INFOPATH: string;
	export const APPLICATION_INSIGHTS_NO_STATSBEAT: string;
	export const GTK_IM_MODULE: string;
	export const LOGNAME: string;
	export const TABLE_MATCHES: string;
	export const NAME: string;
	export const WSL_INTEROP: string;
	export const VSCODE_HANDLES_SIGPIPE: string;
	export const LESS_TERMCAP_us: string;
	export const GEMINI_CLI: string;
	export const SURFACE: string;
	export const QT_AUTO_SCREEN_SCALE_FACTOR: string;
	export const PULSE_SERVER: string;
	export const _: string;
	export const npm_config_prefix: string;
	export const npm_config_npm_version: string;
	export const GOMODCACHE: string;
	export const TABLE_USERS: string;
	export const TERM: string;
	export const npm_config_cache: string;
	export const GOTELEMETRY_GOPLS_CLIENT_TOKEN: string;
	export const npm_config_node_gyp: string;
	export const PATH: string;
	export const USE_CCPA: string;
	export const HOMEBREW_CELLAR: string;
	export const NODE: string;
	export const npm_package_name: string;
	export const XDG_RUNTIME_DIR: string;
	export const DISPLAY: string;
	export const DB_ROOT_PASSWORD: string;
	export const LANG: string;
	export const XMODIFIERS: string;
	export const XAUTHORITY: string;
	export const LS_COLORS: string;
	export const npm_lifecycle_script: string;
	export const DefaultIMModule: string;
	export const SHELL: string;
	export const GOPROXY: string;
	export const GOPATH: string;
	export const npm_package_version: string;
	export const npm_lifecycle_event: string;
	export const ELECTRON_RUN_AS_NODE: string;
	export const LESS_TERMCAP_mb: string;
	export const LESS_TERMCAP_md: string;
	export const TABLE_TEAMS: string;
	export const QT_IM_MODULE: string;
	export const npm_config_globalconfig: string;
	export const npm_config_init_module: string;
	export const PWD: string;
	export const LESS_TERMCAP_me: string;
	export const LC_ALL: string;
	export const npm_execpath: string;
	export const NVM_CD_FLAGS: string;
	export const PYENV_ROOT: string;
	export const XDG_DATA_DIRS: string;
	export const npm_config_global_prefix: string;
	export const INPUt_method: string;
	export const HOMEBREW_REPOSITORY: string;
	export const npm_command: string;
	export const LETSENCRYPT_EMAIL: string;
	export const DB_HOST: string;
	export const WSL2_GUI_APPS_ENABLED: string;
	export const DB_USER: string;
	export const HOSTTYPE: string;
	export const WSLENV: string;
	export const INIT_CWD: string;
	export const EDITOR: string;
	export const NODE_ENV: string;
}

/**
 * Similar to [`$env/static/private`](https://svelte.dev/docs/kit/$env-static-private), except that it only includes environment variables that begin with [`config.kit.env.publicPrefix`](https://svelte.dev/docs/kit/configuration#env) (which defaults to `PUBLIC_`), and can therefore safely be exposed to client-side code.
 * 
 * Values are replaced statically at build time.
 * 
 * ```ts
 * import { PUBLIC_BASE_URL } from '$env/static/public';
 * ```
 */
declare module '$env/static/public' {
	
}

/**
 * This module provides access to runtime environment variables, as defined by the platform you're running on. For example if you're using [`adapter-node`](https://github.com/sveltejs/kit/tree/main/packages/adapter-node) (or running [`vite preview`](https://svelte.dev/docs/kit/cli)), this is equivalent to `process.env`. This module only includes variables that _do not_ begin with [`config.kit.env.publicPrefix`](https://svelte.dev/docs/kit/configuration#env) _and do_ start with [`config.kit.env.privatePrefix`](https://svelte.dev/docs/kit/configuration#env) (if configured).
 * 
 * This module cannot be imported into client-side code.
 * 
 * ```ts
 * import { env } from '$env/dynamic/private';
 * console.log(env.DEPLOYMENT_SPECIFIC_VARIABLE);
 * ```
 * 
 * > [!NOTE] In `dev`, `$env/dynamic` always includes environment variables from `.env`. In `prod`, this behavior will depend on your adapter.
 */
declare module '$env/dynamic/private' {
	export const env: {
		VITE_API_BASE_URL: string;
		VITE_APP_TITLE: string;
		VITE_ENABLE_POLLING: string;
		VITE_POLLING_INTERVAL: string;
		LESS_TERMCAP_se: string;
		VSCODE_CWD: string;
		VSCODE_ESM_ENTRYPOINT: string;
		USER: string;
		LESS_TERMCAP_ue: string;
		VSCODE_NLS_CONFIG: string;
		GOTELEMETRY_GOPLS_CLIENT_START_TIME: string;
		npm_config_user_agent: string;
		VSCODE_WSL_EXT_LOCATION: string;
		VSCODE_HANDLES_UNCAUGHT_ERRORS: string;
		DB_PORT: string;
		GEMINI_YOLO_MODE: string;
		DEBUG: string;
		npm_node_execpath: string;
		JWT_SECRET_KEY: string;
		SHLVL: string;
		LD_LIBRARY_PATH: string;
		npm_config_noproxy: string;
		HOME: string;
		OLDPWD: string;
		NVM_BIN: string;
		VSCODE_IPC_HOOK_CLI: string;
		npm_package_json: string;
		NVM_INC: string;
		HOMEBREW_PREFIX: string;
		LESS_TERMCAP_so: string;
		DB_DATABASE: string;
		npm_config_userconfig: string;
		npm_config_local_prefix: string;
		NMAP_PRIVILEGED: string;
		TABLE_TOURNAMENTS: string;
		DBUS_SESSION_BUS_ADDRESS: string;
		WSL_DISTRO_NAME: string;
		COLOR: string;
		NVM_DIR: string;
		QT_QPA_PLATFORMTHEME: string;
		WAYLAND_DISPLAY: string;
		VSCODE_L10N_BUNDLE_LOCATION: string;
		INFOPATH: string;
		APPLICATION_INSIGHTS_NO_STATSBEAT: string;
		GTK_IM_MODULE: string;
		LOGNAME: string;
		TABLE_MATCHES: string;
		NAME: string;
		WSL_INTEROP: string;
		VSCODE_HANDLES_SIGPIPE: string;
		LESS_TERMCAP_us: string;
		GEMINI_CLI: string;
		SURFACE: string;
		QT_AUTO_SCREEN_SCALE_FACTOR: string;
		PULSE_SERVER: string;
		_: string;
		npm_config_prefix: string;
		npm_config_npm_version: string;
		GOMODCACHE: string;
		TABLE_USERS: string;
		TERM: string;
		npm_config_cache: string;
		GOTELEMETRY_GOPLS_CLIENT_TOKEN: string;
		npm_config_node_gyp: string;
		PATH: string;
		USE_CCPA: string;
		HOMEBREW_CELLAR: string;
		NODE: string;
		npm_package_name: string;
		XDG_RUNTIME_DIR: string;
		DISPLAY: string;
		DB_ROOT_PASSWORD: string;
		LANG: string;
		XMODIFIERS: string;
		XAUTHORITY: string;
		LS_COLORS: string;
		npm_lifecycle_script: string;
		DefaultIMModule: string;
		SHELL: string;
		GOPROXY: string;
		GOPATH: string;
		npm_package_version: string;
		npm_lifecycle_event: string;
		ELECTRON_RUN_AS_NODE: string;
		LESS_TERMCAP_mb: string;
		LESS_TERMCAP_md: string;
		TABLE_TEAMS: string;
		QT_IM_MODULE: string;
		npm_config_globalconfig: string;
		npm_config_init_module: string;
		PWD: string;
		LESS_TERMCAP_me: string;
		LC_ALL: string;
		npm_execpath: string;
		NVM_CD_FLAGS: string;
		PYENV_ROOT: string;
		XDG_DATA_DIRS: string;
		npm_config_global_prefix: string;
		INPUt_method: string;
		HOMEBREW_REPOSITORY: string;
		npm_command: string;
		LETSENCRYPT_EMAIL: string;
		DB_HOST: string;
		WSL2_GUI_APPS_ENABLED: string;
		DB_USER: string;
		HOSTTYPE: string;
		WSLENV: string;
		INIT_CWD: string;
		EDITOR: string;
		NODE_ENV: string;
		[key: `PUBLIC_${string}`]: undefined;
		[key: `${string}`]: string | undefined;
	}
}

/**
 * Similar to [`$env/dynamic/private`](https://svelte.dev/docs/kit/$env-dynamic-private), but only includes variables that begin with [`config.kit.env.publicPrefix`](https://svelte.dev/docs/kit/configuration#env) (which defaults to `PUBLIC_`), and can therefore safely be exposed to client-side code.
 * 
 * Note that public dynamic environment variables must all be sent from the server to the client, causing larger network requests — when possible, use `$env/static/public` instead.
 * 
 * ```ts
 * import { env } from '$env/dynamic/public';
 * console.log(env.PUBLIC_DEPLOYMENT_SPECIFIC_VARIABLE);
 * ```
 */
declare module '$env/dynamic/public' {
	export const env: {
		[key: `PUBLIC_${string}`]: string | undefined;
	}
}
