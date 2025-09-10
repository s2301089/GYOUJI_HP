

export const index = 2;
let component_cache;
export const component = async () => component_cache ??= (await import('../entries/pages/_page.svelte.js')).default;
export const imports = ["_app/immutable/nodes/2.DWfazvA1.js","_app/immutable/chunks/DsnmJJEf.js","_app/immutable/chunks/DgMYUmlK.js","_app/immutable/chunks/SJRvxFS1.js"];
export const stylesheets = ["_app/immutable/assets/2.DSAcDFkb.css"];
export const fonts = [];
