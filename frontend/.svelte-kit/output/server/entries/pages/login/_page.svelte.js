import { e as escape_html } from "../../../chunks/escaping.js";
import "clsx";
import { v as pop, t as push } from "../../../chunks/index2.js";
import "@sveltejs/kit/internal";
import "../../../chunks/exports.js";
import "../../../chunks/utils.js";
import "../../../chunks/state.svelte.js";
const replacements = {
  translate: /* @__PURE__ */ new Map([
    [true, "yes"],
    [false, "no"]
  ])
};
function attr(name, value, is_boolean = false) {
  if (is_boolean) return "";
  const normalized = name in replacements && replacements[name].get(value) || value;
  const assignment = is_boolean ? "" : `="${escape_html(normalized, true)}"`;
  return ` ${name}${assignment}`;
}
function _page($$payload, $$props) {
  push();
  let username = "";
  let password = "";
  $$payload.out.push(`<div class="login-container svelte-4lnbfd"><div class="login-card svelte-4lnbfd"><h2 class="svelte-4lnbfd">ログイン</h2> <form><div class="input-group svelte-4lnbfd"><label for="username" class="svelte-4lnbfd">ユーザー名</label> <input type="text" id="username"${attr("value", username)} required class="svelte-4lnbfd"/></div> <div class="input-group svelte-4lnbfd"><label for="password" class="svelte-4lnbfd">パスワード</label> <input type="password" id="password"${attr("value", password)} required class="svelte-4lnbfd"/></div> `);
  {
    $$payload.out.push("<!--[!-->");
  }
  $$payload.out.push(`<!--]--> <button type="submit" class="svelte-4lnbfd">ログイン</button></form></div></div>`);
  pop();
}
export {
  _page as default
};
