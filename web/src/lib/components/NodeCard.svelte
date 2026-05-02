<script>
  export let node;
  export let active = false;

  const statusTone = {
    online: 'border-cyan-400/40 bg-cyan-400/5 text-cyan-200',
    warn: 'border-amber-400/40 bg-amber-400/10 text-amber-100',
    offline: 'border-pink-400/40 bg-pink-400/10 text-pink-100'
  };

  function pct(value) {
    return `${Math.round(value || 0)}%`;
  }
</script>

<div class={`panel w-full p-5 text-left transition hover:-translate-y-1 ${active ? 'ring-1 ring-cyan-300/60' : ''}`}>
  <div class="flex items-start justify-between gap-4">
    <div>
      <div class="text-lg font-semibold text-white">{node.hostname}</div>
      <div class="mt-1 text-xs uppercase tracking-[0.25em] text-slate-400">{node.nodeId}</div>
    </div>
    <span class={`rounded-full border px-3 py-1 text-xs uppercase tracking-[0.2em] ${statusTone[node.status]}`}>{node.status}</span>
  </div>

  <div class="mt-5 grid grid-cols-3 gap-3 text-sm text-slate-300">
    <div>
      <div class="text-slate-500">CPU</div>
      <div class="mt-1 text-cyan-200">{pct(node.cpuUsage)}</div>
    </div>
    <div>
      <div class="text-slate-500">MEM</div>
      <div class="mt-1 text-cyan-200">{pct(node.memory?.usage)}</div>
    </div>
    <div>
      <div class="text-slate-500">DISK</div>
      <div class="mt-1 text-cyan-200">{pct(node.disk?.usage)}</div>
    </div>
  </div>

  <div class="mt-5 text-sm text-slate-400">{node.statusSummary}</div>
</div>
