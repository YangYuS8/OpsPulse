<script>
  export let node = null;

  const badgeTone = (state) => {
    if (state === 'running') return 'text-emerald-300';
    if (state === 'exited') return 'text-amber-300';
    return 'text-rose-300';
  };

  const badge = (state) => {
    if (state === 'running') return '[ OK ]';
    if (state === 'exited') return '[WARN]';
    return '[FAIL]';
  };
</script>

<section class="pane pane-scroll pane-containers">
  <div class="pane-title">CONTAINERS</div>

  {#if node && node.docker?.containers?.length}
    <div class="mt-1 space-y-px text-sm">
      {#each node.docker.containers as container}
        <div class="service-row">
          <span class={`w-[56px] ${badgeTone(container.state)}`}>{badge(container.state)}</span>
          <span class="w-[110px] truncate text-terminal-fg">{container.name}</span>
          <span class="w-[120px] truncate text-terminal-dim">{container.image}</span>
          <span class="flex-1 truncate text-terminal-dim text-right">{container.status}</span>
        </div>
      {/each}
    </div>
  {:else if node}
    <div class="empty-pane py-6 text-sm text-terminal-dim">no containers reported for this node</div>
  {:else}
    <div class="empty-pane py-6 text-sm text-terminal-dim">select a node to inspect containers</div>
  {/if}
</section>
