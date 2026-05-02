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
  <section class="pane pane-scroll px-4 py-3 sm:px-4 sm:py-3">
    <div class="pane-title">NodeInspectorPane</div>
    <div class="mt-2 flex items-center justify-between gap-4 border-b border-terminal-border pb-2">
      <div>
        <h2 class="text-base text-terminal-fg">inspect {node.hostname}</h2>
        <p class="mt-1 text-xs text-terminal-dim">last_seen {new Date(node.lastSeen).toLocaleString()} node_id {node.nodeId}</p>
      </div>
      <div class={node.status === 'online' ? 'status-chip-ok' : node.status === 'warn' ? 'status-chip-warn' : 'status-chip-fail'}>{node.status === 'online' ? '[ OK ]' : node.status === 'warn' ? '[WARN]' : '[FAIL]'}</div>
    </div>

    <div class="inspector-grid mt-3">
      <div>
        <div class="section-label">uptime</div>
        <div class="mt-1 text-sm text-terminal-fg">{formatUptime(node.uptime)}</div>
      </div>
      <div>
        <div class="section-label">load avg</div>
        <div class="mt-1 text-sm text-terminal-fg">{node.load.one} / {node.load.five} / {node.load.fifteen}</div>
      </div>
      <div>
        <div class="section-label">docker</div>
        <div class="mt-1 text-sm text-terminal-fg">{node.docker.running} running / {node.docker.exited} exited</div>
      </div>
      <div>
        <div class="section-label">heartbeat age</div>
        <div class="mt-1 text-sm text-terminal-fg">{node.heartbeatAgeSeconds}s</div>
      </div>
    </div>

    <div class="mt-4 border-t border-terminal-border pt-3">
      <h3 class="section-label">resource usage</h3>
      <div class="mt-3 space-y-4 text-sm">
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

    <div class="mt-4 border-t border-terminal-border pt-3">
      <h3 class="section-label">summary</h3>
      <div class="mt-2 text-sm text-terminal-dim">{node.statusSummary}</div>
    </div>
  </section>
{:else}
  <section class="pane pane-scroll flex min-h-[320px] items-center justify-center px-4 py-3 text-terminal-dim">
    no nodes connected; start opspulse-agent on a host to begin telemetry
  </section>
{/if}
