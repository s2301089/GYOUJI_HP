import { w as slot } from "../../chunks/index2.js";
function _layout($$payload, $$props) {
  $$payload.out.push(`<div class="app-container svelte-p8n3xy"><header class="svelte-p8n3xy"><div class="container header-content svelte-p8n3xy"><div class="header-title-group svelte-p8n3xy"><a href="/" class="svelte-p8n3xy"><img src="/icon.svg" alt="行事委員会 Icon" class="header-icon svelte-p8n3xy"/> <h1 class="svelte-p8n3xy">行事委員会</h1></a></div> <nav class="desktop-nav svelte-p8n3xy"><a href="/#home" class="svelte-p8n3xy">ホーム</a> <a href="/#about" class="svelte-p8n3xy">概要</a> <a href="/#events" class="svelte-p8n3xy">イベント</a> <a href="/#roles" class="svelte-p8n3xy">役職</a> <a href="/#join" class="svelte-p8n3xy">参加方法</a></nav> <button class="hamburger-menu svelte-p8n3xy" aria-label="メニューを開閉する"><span class="line svelte-p8n3xy"></span> <span class="line svelte-p8n3xy"></span> <span class="line svelte-p8n3xy"></span></button></div></header> `);
  {
    $$payload.out.push("<!--[!-->");
  }
  $$payload.out.push(`<!--]--> <main class="svelte-p8n3xy"><!---->`);
  slot($$payload, $$props, "default", {});
  $$payload.out.push(`<!----></main> <footer class="svelte-p8n3xy"><p>© 2025 仙台高専広瀬キャンパス 行事委員会</p></footer></div>`);
}
export {
  _layout as default
};
