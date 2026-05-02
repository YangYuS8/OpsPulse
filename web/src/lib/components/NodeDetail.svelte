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
</script>

{#if node}
  <section class="panel p-6">
    <div class="flex items-center justify-between gap-4">
      <div>
        <h2 class="text-2xl font-semibold text-white">{node.hostname}</h2>
        <p class="mt-1 text-sm text-slate-400">Last seen {new Date(node.lastSeen).toLocaleString()}</p>
      </div>
      <div class="rounded-full border border-cyan-400/30 px-4 py-2 text-sm uppercase tracking-[0.25em] text-cyan-200">
        {node.status}
      </div>
    </div>

    <div class="mt-6 grid gap-4 md:grid-cols-2 xl:grid-cols-4">
      <div class="rounded-xl border border-slate-800 bg-slate-950/60 p-4">
        <div class="text-xs uppercase tracking-[0.25em] text-slate-500">Uptime</div>
        <div class="mt-2 text-xl text-white">{formatUptime(node.uptime)}</div>
      </div>
      <div class="rounded-xl border border-slate-800 bg-slate-950/60 p-4">
        <div class="text-xs uppercase tracking-[0.25em] text-slate-500">Load Avg</div>
        <div class="mt-2 text-xl text-white">{node.load.one} / {node.load.five} / {node.load.fifteen}</div>
      </div>
      <div class="rounded-xl border border-slate-800 bg-slate-950/60 p-4">
        <div class="text-xs uppercase tracking-[0.25em] text-slate-500">Docker</div>
        <div class="mt-2 text-xl text-white">{node.docker.running} running / {node.docker.exited} exited</div>
      </div>
      <div class="rounded-xl border border-slate-800 bg-slate-950/60 p-4">
        <div class="text-xs uppercase tracking-[0.25em] text-slate-500">Heartbeat Age</div>
        <div class="mt-2 text-xl text-white">{node.heartbeatAgeSeconds}s</div>
      </div>
    </div>

    <div class="mt-6 grid gap-6 xl:grid-cols-[1.2fr_0.8fr]">
      <div class="rounded-xl border border-slate-800 bg-slate-950/60 p-4">
        <h3 class="text-sm uppercase tracking-[0.25em] text-slate-400">Resource Usage</h3>
        <div class="mt-4 space-y-4">
          <div>
            <div class="mb-1 flex justify-between text-sm text-slate-300"><span>Memory</span><span>{Math.round(node.memory.usage)}%</span></div>
            <div class="h-2 overflow-hidden rounded-full bg-slate-900"><div class="h-full bg-cyan-400" style={`width:${node.memory.usage}%`}></div></div>
            <div class="mt-1 text-xs text-slate-500">{formatBytes(node.memory.used)} / {formatBytes(node.memory.total)}</div>
          </div>
          <div>
            <div class="mb-1 flex justify-between text-sm text-slate-300"><span>Disk</span><span>{Math.round(node.disk.usage)}%</span></div>
            <div class="h-2 overflow-hidden rounded-full bg-slate-900"><div class="h-full bg-pink-400" style={`width:${node.disk.usage}%`}></div></div>
            <div class="mt-1 text-xs text-slate-500">{formatBytes(node.disk.used)} / {formatBytes(node.disk.total)}</div>
          </div>
        </div>
      </div>

      <div class="rounded-xl border border-slate-800 bg-slate-950/60 p-4">
        <h3 class="text-sm uppercase tracking-[0.25em] text-slate-400">Services</h3>
        <div class="mt-4 space-y-3">
          {#each node.services as service}
            <div class="flex items-center justify-between rounded-lg border border-slate-800 px-3 py-2 text-sm">
              <div>
                <div class="text-white">{service.name}</div>
                <div class="text-xs text-slate-500">{service.active} / {service.sub}</div>
              </div>
              <span class={`rounded-full px-3 py-1 text-xs uppercase tracking-[0.2em] ${service.status === 'healthy' ? 'bg-lime-400/15 text-lime-200' : 'bg-pink-400/15 text-pink-200'}`}>
                {service.status}
              </span>
            </div>
          {/each}
        </div>
      </div>
    </div>
  </section>
{:else}
  <section class="panel flex min-h-[420px] items-center justify-center p-6 text-slate-400">
    Select a node to inspect live service and container state.
  </section>
{/if}
