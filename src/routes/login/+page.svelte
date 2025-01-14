<script lang="ts">
	import type { PageServerData } from './$types';
	import LoginForm from './login-form.svelte';

	let { data }: { data: PageServerData } = $props();
</script>

<div class="h-screen bg-[url('/images/login-background.jpg')] bg-cover">
	<div class="flex min-h-screen items-center justify-center">
		<LoginForm data={data.form} />
	</div>
</div>
