<script>
  export let node;

  const formatBytes = (value) => {
    if (!value) return '0 B';
    const units = ['B', 'KB', 'MB', 'GB', 'TB'];
    let size = value;
    let index = 0;
    while (size >= 1024 && index < units.length - 1) {
      size /= 1024;
      index += 1;
    }
    return `${size.toFixed(size >= 10 ? 0 : 1)} ${units[index]}`;
  };

  const formatUptime = (seconds) => {
    const days = Math.floor(seconds / 86400);
    const hours = Math.floor((seconds % 86400) / 3600);
    const minutes = Math.floor((seconds % 3600) / 60);
    return `${days}d ${hours}h ${minutes}m`;
  };

  const serviceBadge = (status) => {
    if (status === 'healthy') return '[ OK ]';
    if (status === 'degraded') return '[WARN]';
    return '[FAIL]';
  };

  const meter = (value) => {
    const width = 24;
    const filled = Math.max(0, Math.min(width, Math.round(((value || 0) / 100) * width)));
    return `${'='.repeat(filled)}${'.'.repeat(width - filled)}`;
  };
</script>

{#if node}
  <section class="terminal-pane p-4 sm:p-5">
    <div class="flex items-center justify-between gap-4 border-b border-terminal-border pb-3">
      <div>
        <h2 class="text-xl text-terminal-fg">inspect {node.hostname}</h2>
        <p class="mt-1 text-xs text-terminal-dim">last_seen {new Date(node.lastSeen).toLocaleString()}</p>
      </div>
      <div class={node.status === 'online' ? 'status-chip-ok' : node.status === 'warn' ? 'status-chip-warn' : 'status-chip-fail'}>{node.status === 'online' ? '[ OK ]' : node.status === 'warn' ? '[WARN]' : '[FAIL]'}</div>
    </div>

    <div class="mt-4 grid gap-3 md:grid-cols-2 xl:grid-cols-4">
      <div class="rounded border border-terminal-border bg-terminal-bg p-3">
        <div class="text-[11px] uppercase tracking-[0.24em] text-terminal-muted">uptime</div>
        <div class="mt-2 text-sm text-terminal-fg">{formatUptime(node.uptime)}</div>
      </div>
      <div class="rounded border border-terminal-border bg-terminal-bg p-3">
        <div class="text-[11px] uppercase tracking-[0.24em] text-terminal-muted">load avg</div>
        <div class="mt-2 text-sm text-terminal-fg">{node.load.one} / {node.load.five} / {node.load.fifteen}</div>
      </div>
      <div class="rounded border border-terminal-border bg-terminal-bg p-3">
        <div class="text-[11px] uppercase tracking-[0.24em] text-terminal-muted">docker</div>
        <div class="mt-2 text-sm text-terminal-fg">{node.docker.running} running / {node.docker.exited} exited</div>
      </div>
      <div class="rounded border border-terminal-border bg-terminal-bg p-3">
        <div class="text-[11px] uppercase tracking-[0.24em] text-terminal-muted">heartbeat age</div>
        <div class="mt-2 text-sm text-terminal-fg">{node.heartbeatAgeSeconds}s</div>
      </div>
    </div>

    <div class="mt-4 grid gap-4 xl:grid-cols-[1.15fr_0.85fr]">
      <div class="rounded border border-terminal-border bg-terminal-bg p-4">
        <h3 class="text-[11px] uppercase tracking-[0.24em] text-terminal-muted">resource usage</h3>
        <div class="mt-4 space-y-4 text-sm">
          <div>
            <div class="mb-1 flex justify-between text-terminal-dim"><span>cpu</span><span>{Math.round(node.cpuUsage)}%</span></div>
            <div class="meter text-terminal-green">[{meter(node.cpuUsage)}]</div>
          </div>
          <div>
            <div class="mb-1 flex justify-between text-terminal-dim"><span>memory</span><span>{Math.round(node.memory.usage)}%</span></div>
            <div class="meter text-terminal-green">[{meter(node.memory.usage)}]</div>
            <div class="mt-1 text-xs text-terminal-muted">{formatBytes(node.memory.used)} / {formatBytes(node.memory.total)}</div>
          </div>
          <div>
            <div class="mb-1 flex justify-between text-terminal-dim"><span>disk</span><span>{Math.round(node.disk.usage)}%</span></div>
            <div class="meter text-terminal-green">[{meter(node.disk.usage)}]</div>
            <div class="mt-1 text-xs text-terminal-muted">{formatBytes(node.disk.used)} / {formatBytes(node.disk.total)}</div>
          </div>
        </div>
      </div>

      <div class="rounded border border-terminal-border bg-terminal-bg p-4">
        <h3 class="text-[11px] uppercase tracking-[0.24em] text-terminal-muted">services</h3>
        <div class="mt-4 space-y-3 text-sm">
          {#each node.services as service}
            <div class="flex items-center justify-between rounded border border-terminal-border px-3 py-2">
              <div>
                <div class="text-terminal-fg">{service.name}</div>
                <div class="text-xs text-terminal-muted">{service.active} / {service.sub}</div>
              </div>
              <span class={service.status === 'healthy' ? 'status-chip-ok' : service.status === 'degraded' ? 'status-chip-warn' : 'status-chip-fail'}>{serviceBadge(service.status)}</span>
            </div>
          {/each}
        </div>
      </div>
    </div>
  </section>
{:else}
  <section class="terminal-pane flex min-h-[420px] items-center justify-center p-6 text-terminal-dim">
    inspect pane idle; select a node or run inspect &lt;node&gt;
  </section>
{/if}
