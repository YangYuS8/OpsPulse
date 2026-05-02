<script>
  export let nodes = [];
  export let selectedId = null;
  export let onSelect = () => {};
  export let onKeySelect = () => {};

  function pct(value) {
    return `${Math.round(value || 0)}%`;
  }

  function age(node) {
    return `${node.heartbeatAgeSeconds ?? 0}s`;
  }

  function badge(status) {
    if (status === 'online') return '[ OK ]';
    if (status === 'warn') return '[WARN]';
    return '[FAIL]';
  }

  function badgeTone(status) {
    if (status === 'online') return 'text-emerald-300';
    if (status === 'warn') return 'text-amber-300';
    return 'text-rose-300';
  }
</script>

<section class="pane pane-scroll">
  <div class="pane-title">NodeListPane</div>
  <div class="table-head text-terminal-muted">state host cpu mem disk age</div>

  {#if nodes.length > 0}
    <div class="mt-1 space-y-px">
      {#each nodes as node}
        <div
          class={`node-row ${node.nodeId === selectedId ? 'node-row-active' : ''}`}
          on:click={() => onSelect(node.nodeId)}
          on:keydown={(event) => onKeySelect(event, node.nodeId)}
          role="button"
          tabindex="0"
        >
          <span class={`w-[56px] ${badgeTone(node.status)}`}>{badge(node.status)}</span>
          <span class="w-[112px] truncate text-terminal-fg">{node.hostname}</span>
          <span class="w-[64px] text-right">cpu {pct(node.cpuUsage)}</span>
          <span class="w-[64px] text-right">mem {pct(node.memory?.usage)}</span>
          <span class="w-[68px] text-right">disk {pct(node.disk?.usage)}</span>
          <span class="w-[52px] text-right">{age(node)}</span>
        </div>
      {/each}
    </div>
  {:else}
    <div class="empty-pane py-6 text-sm text-terminal-dim">
      no nodes connected; start opspulse-agent on a host to begin telemetry
    </div>
  {/if}
</section>
