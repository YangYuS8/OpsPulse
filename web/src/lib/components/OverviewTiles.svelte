<script>
  export let overview;

  const items = [
    { key: 'nodesOnline', label: 'Online Nodes', tone: 'text-cyan-300' },
    { key: 'nodesWarn', label: 'Warning Nodes', tone: 'text-amber-300' },
    { key: 'nodesOffline', label: 'Offline Nodes', tone: 'text-pink-300' },
    { key: 'containersUp', label: 'Running Containers', tone: 'text-lime-300' }
  ];
</script>

<div class="grid gap-4 md:grid-cols-2 xl:grid-cols-4">
  {#each items as item}
    <section class="panel p-5">
      <div class="text-xs uppercase tracking-[0.3em] text-slate-400">{item.label}</div>
      <div class={`mt-3 text-4xl font-semibold ${item.tone}`}>{overview?.[item.key] ?? 0}</div>
    </section>
  {/each}
</div>
