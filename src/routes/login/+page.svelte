<script lang="ts">
	import type { PageData } from './$types';
	import LoginForm from './login-form.svelte';

	let { data }: { data: PageData } = $props();
</script>

<div class="h-screen bg-[url('/images/login-background.jpg')] bg-cover">
	<div class="flex min-h-screen items-center justify-center">
		<LoginForm data={data.form} />
	</div>
</div>
