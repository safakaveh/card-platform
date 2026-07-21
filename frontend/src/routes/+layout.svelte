<script lang="ts">
	import './layout.css';
	import favicon from '$lib/assets/favicon.png';
	import NavigationMenu from '@/components/app/navigation-menu.svelte';

	let { children } = $props();
</script>



<svelte:head>
	<link rel="icon" href={favicon} />
	<title>برنامه دستگاه پرس</title>
</svelte:head>
<main class="relative min-h-screen w-full overflow-hidden bg-slate-100">
  <div class="absolute inset-0 bg-[linear-gradient(to_right,rgba(148,163,184,0.18)_1px,transparent_1px),linear-gradient(to_bottom,rgba(148,163,184,0.18)_1px,transparent_1px)] bg-[size:56px_56px]" />

  <div class="absolute -top-24 left-0 h-[320px] w-[320px] rounded-full bg-cyan-300/35 blur-3xl" />
  <div class="absolute right-0 top-0 h-[280px] w-[280px] rounded-full bg-sky-300/30 blur-3xl" />
  <div class="absolute bottom-0 left-1/3 h-[260px] w-[260px] rounded-full bg-blue-300/25 blur-3xl" />

  <div class="relative z-20 w-full px-4 pt-4">
    <div class="rounded-2xl border border-white/70 bg-white/75 px-4 py-3 shadow-[0_8px_30px_rgba(15,23,42,0.08)] backdrop-blur-xl">
      <NavigationMenu />
    </div>
  </div>

  <div class="relative z-10 flex min-h-[calc(100vh-80px)] w-full items-center justify-center p-6">
    {@render children()}
  </div>
</main>
