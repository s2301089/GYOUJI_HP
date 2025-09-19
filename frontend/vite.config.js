import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';
import { SvelteKitPWA } from '@vite-pwa/sveltekit';

export default defineConfig({
	plugins: [
		sveltekit(),
		SvelteKitPWA({
			registerType: 'autoUpdate',
			includeAssets: ['icon.svg'],
			manifest: {
				name: '行事委員会結果速報アプリ',
				short_name: 'GYOUJI_HP',
				theme_color: '#ffffff',
				icons: [
					{
						src: 'icon.svg',
						sizes: 'any',
						type: 'image/svg+xml'
					}
				]
			}
		})
	],
	server: {
		proxy: {
			'/swagger': {
				target: 'http://localhost:8080',
				changeOrigin: true,
			},
			'/api': {
				target: 'http://localhost:8080',
				changeOrigin: true,
				secure: false,
			}
		}
	}
});