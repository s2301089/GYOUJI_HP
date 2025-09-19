import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';
import { SvelteKitPWA } from '@vite-pwa/sveltekit';

export default defineConfig({
	plugins: [
		sveltekit(),
		SvelteKitPWA({
  			registerType: 'autoUpdate',
  			includeAssets: ['icon.svg', 'icon-192.png', 'icon-512.png'],
  			manifest: {
    		name: '行事委員会結果速報アプリ',
    		short_name: 'GYOUJI_HP',
    		start_url: '/',
    		display: 'standalone',
    		background_color: '#ffffff',
    		theme_color: '#ffffff',
    		icons: [
      			{
        			src: 'icon-192.png',
        			sizes: '192x192',
        			type: 'image/png'
      			},
      			{
        			src: 'icon-512.png',
        			sizes: '512x512',
        			type: 'image/png'
      			},
      			{
        			src: 'icon-512.png',
        			sizes: '512x512',
        			type: 'image/png',
        			purpose: 'any maskable'
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