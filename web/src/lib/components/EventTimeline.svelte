<script>
  export let events = [];
  export let title = 'event log';
</script>

<section class="terminal-pane p-4 sm:p-5">
  <div class="flex items-center justify-between gap-4 border-b border-terminal-border pb-3">
    <h2 class="text-sm uppercase tracking-[0.24em] text-terminal-fg">{title}</h2>
    <span class="text-[11px] uppercase tracking-[0.2em] text-terminal-muted">recent 50</span>
  </div>
  <div class="mt-4 space-y-2 font-mono text-sm">
    {#each events as event}
      <div class="rounded border border-terminal-border bg-terminal-bg px-3 py-2">
        <div class="flex flex-wrap items-center gap-2 text-xs">
          <span class="text-terminal-muted">[{new Date(event.createdAt).toLocaleTimeString()}]</span>
          <span class={event.level === 'warn' ? 'status-chip-warn' : event.level === 'info' ? 'status-chip-ok' : 'status-chip-fail'}>
            {event.level === 'warn' ? '[WARN]' : event.level === 'info' ? '[ OK ]' : '[FAIL]'}
          </span>
          <span class="text-terminal-fg">{event.nodeId}</span>
        </div>
        <p class="mt-2 text-sm text-terminal-dim">{event.message}</p>
      </div>
    {/each}
  </div>
</section>
