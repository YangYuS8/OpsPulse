<script>
  export let overview;

  const items = [
    { key: 'nodesOnline', label: '[ OK ] online', tone: 'text-emerald-300' },
    { key: 'nodesWarn', label: '[WARN] watch', tone: 'text-amber-300' },
    { key: 'nodesOffline', label: '[FAIL] down', tone: 'text-rose-300' },
    { key: 'containersUp', label: '[ OK ] ctrs', tone: 'text-terminal-green' }
  ];
</script>

<div class="grid gap-3 md:grid-cols-2 xl:grid-cols-4">
  {#each items as item}
    <section class="terminal-pane p-4">
      <div class="text-[11px] uppercase tracking-[0.3em] text-terminal-muted">{item.label}</div>
      <div class={`mt-2 text-3xl ${item.tone}`}>{overview?.[item.key] ?? 0}</div>
    </section>
  {/each}
</div>
