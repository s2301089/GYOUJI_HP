export const manifest = (() => {
function __memo(fn) {
	let value;
	return () => value ??= (value = fn());
}

return {
	appDir: "_app",
	appPath: "_app",
	assets: new Set([]),
	mimeTypes: {},
	_: {
		client: {start:"_app/immutable/entry/start.CtEohA1W.js",app:"_app/immutable/entry/app.DWnkMJjL.js",imports:["_app/immutable/entry/start.CtEohA1W.js","_app/immutable/chunks/Cr-ZIwDU.js","_app/immutable/chunks/CVVBCw2J.js","_app/immutable/chunks/SJRvxFS1.js","_app/immutable/chunks/CEKcJX13.js","_app/immutable/entry/app.DWnkMJjL.js","_app/immutable/chunks/SJRvxFS1.js","_app/immutable/chunks/CVVBCw2J.js","_app/immutable/chunks/CEKcJX13.js","_app/immutable/chunks/DsnmJJEf.js","_app/immutable/chunks/Dv_ePxNb.js"],stylesheets:[],fonts:[],uses_env_dynamic_public:false},
		nodes: [
			__memo(() => import('./nodes/0.js')),
			__memo(() => import('./nodes/1.js')),
			__memo(() => import('./nodes/2.js')),
			__memo(() => import('./nodes/3.js')),
			__memo(() => import('./nodes/4.js'))
		],
		remotes: {
			
		},
		routes: [
			{
				id: "/",
				pattern: /^\/$/,
				params: [],
				page: { layouts: [0,], errors: [1,], leaf: 2 },
				endpoint: null
			},
			{
				id: "/dashboard",
				pattern: /^\/dashboard\/?$/,
				params: [],
				page: { layouts: [0,], errors: [1,], leaf: 3 },
				endpoint: null
			},
			{
				id: "/login",
				pattern: /^\/login\/?$/,
				params: [],
				page: { layouts: [0,], errors: [1,], leaf: 4 },
				endpoint: null
			}
		],
		prerendered_routes: new Set([]),
		matchers: async () => {
			
			return {  };
		},
		server_assets: {}
	}
}
})();
