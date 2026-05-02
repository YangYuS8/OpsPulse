<script>
  export let events = [];
  export let title = 'EVENTS';

  function label(level) {
    if (level === 'warn') return 'WARN';
    if (level === 'info') return 'INFO';
    return 'FAIL';
  }

  function tone(level) {
    if (level === 'warn') return 'text-amber-300';
    if (level === 'info') return 'text-terminal-fg';
    return 'text-rose-300';
  }
</script>

<section class="pane pane-scroll">
  <div class="pane-title">{title}</div>

  <div class="mt-1 space-y-px text-sm">
    {#each events as event}
      <div class="log-row">
        <span class="w-[74px] text-terminal-muted">{new Date(event.createdAt).toLocaleTimeString()}</span>
        <span class={`w-[42px] ${tone(event.level)}`}>{label(event.level)}</span>
        <span class="w-[110px] truncate text-terminal-fg">{event.nodeId}</span>
        <span class="flex-1 truncate text-terminal-dim">{event.message}</span>
      </div>
    {/each}
  </div>
</section>
