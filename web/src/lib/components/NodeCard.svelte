<script>
  export let node;
  export let active = false;

  const statusTone = {
    online: 'status-chip-ok',
    warn: 'status-chip-warn',
    offline: 'status-chip-fail'
  };

  function pct(value) {
    return `${Math.round(value || 0)}%`;
  }

  function meter(value) {
    const width = 16;
    const filled = Math.max(0, Math.min(width, Math.round(((value || 0) / 100) * width)));
    return `${'#'.repeat(filled)}${'-'.repeat(width - filled)}`;
  }
</script>

<div class={`terminal-pane w-full p-4 text-left ${active ? 'border-emerald-400/80' : ''}`}>
  <div class="flex items-start justify-between gap-4">
    <div>
      <div class="text-base text-terminal-fg">{node.hostname}</div>
      <div class="mt-1 text-[11px] uppercase tracking-[0.24em] text-terminal-muted">{node.nodeId}</div>
    </div>
    <span class={statusTone[node.status]}>{node.status === 'online' ? '[ OK ]' : node.status === 'warn' ? '[WARN]' : '[FAIL]'}</span>
  </div>

  <div class="mt-4 space-y-2 text-xs text-terminal-dim">
    <div>
      <div class="flex justify-between"><span>cpu</span><span>{pct(node.cpuUsage)}</span></div>
      <div class="meter text-terminal-green">[{meter(node.cpuUsage)}]</div>
    </div>
    <div>
      <div class="flex justify-between"><span>mem</span><span>{pct(node.memory?.usage)}</span></div>
      <div class="meter text-terminal-green">[{meter(node.memory?.usage)}]</div>
    </div>
    <div>
      <div class="flex justify-between"><span>disk</span><span>{pct(node.disk?.usage)}</span></div>
      <div class="meter text-terminal-green">[{meter(node.disk?.usage)}]</div>
    </div>
  </div>

  <div class="mt-4 text-xs text-terminal-dim">{node.statusSummary}</div>
</div>
