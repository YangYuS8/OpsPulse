<script>
  export let overview;

  const lines = [
    ['NODES', 'nodesTotal'],
    ['ONLINE', 'nodesOnline'],
    ['WARN', 'nodesWarn'],
    ['OFFLINE', 'nodesOffline'],
    ['CONTAINERS', 'containersUp'],
    ['SERVICE_ALERTS', 'servicesDown']
  ];
</script>

<section class="terminal-pane p-4 sm:p-5">
  <div class="flex flex-col gap-4 xl:flex-row xl:items-end xl:justify-between">
    <div>
      <div class="terminal-title">OpsPulse / Local-First Control Plane</div>
      <div class="mt-2 text-3xl font-semibold text-terminal-fg">Web btop for Linux nodes</div>
      <div class="mt-2 max-w-3xl text-sm text-terminal-dim">
        Read-only observability console with heartbeat telemetry, service watchlists, offline detection and live event stream.
      </div>
    </div>
    <div class="grid gap-2 sm:grid-cols-2 xl:grid-cols-3">
      {#each lines as [label, key]}
        <div class="rounded border border-terminal-border bg-terminal-panel px-3 py-2">
          <div class="text-[10px] uppercase tracking-[0.24em] text-terminal-muted">{label}</div>
          <div class="mt-1 text-lg text-terminal-fg">{overview?.[key] ?? 0}</div>
        </div>
      {/each}
    </div>
  </div>
</section>
