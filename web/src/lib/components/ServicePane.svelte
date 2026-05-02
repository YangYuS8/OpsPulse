<script>
  export let node = null;

  const badge = (status) => {
    if (status === 'healthy') return '[ OK ]';
    if (status === 'degraded') return '[WARN]';
    return '[FAIL]';
  };

  const badgeTone = (status) => {
    if (status === 'healthy') return 'text-emerald-300';
    if (status === 'degraded') return 'text-amber-300';
    return 'text-rose-300';
  };
</script>

<section class="pane pane-scroll">
  <div class="pane-title">ServicePane</div>

  {#if node && node.services?.length}
    <div class="mt-1 space-y-px text-sm">
      {#each node.services as service}
        <div class="service-row">
          <span class={`w-[56px] ${badgeTone(service.status)}`}>{badge(service.status)}</span>
          <span class="flex-1 truncate text-terminal-fg">{service.name}</span>
          <span class="w-[120px] truncate text-terminal-dim text-right">{service.active}/{service.sub}</span>
        </div>
      {/each}
    </div>
  {:else if node}
    <div class="empty-pane py-6 text-sm text-terminal-dim">no whitelisted services reported for this node</div>
  {:else}
    <div class="empty-pane py-6 text-sm text-terminal-dim">select a node to inspect services</div>
  {/if}
</section>
