<script>
  export let events = [];
</script>

<section class="panel p-6">
  <div class="flex items-center justify-between gap-4">
    <h2 class="text-lg font-semibold text-white">Event Timeline</h2>
    <span class="text-xs uppercase tracking-[0.2em] text-slate-400">Recent 50</span>
  </div>
  <div class="mt-4 space-y-4">
    {#each events as event}
      <div class="relative border-l border-cyan-500/20 pl-4">
        <div class={`absolute -left-[5px] top-1.5 h-2.5 w-2.5 rounded-full ${event.level === 'warn' ? 'bg-amber-300' : 'bg-cyan-300'}`}></div>
        <div class="flex flex-wrap items-center gap-2 text-sm">
          <span class="font-medium text-white">{event.nodeId}</span>
          <span class={`rounded-full px-2 py-0.5 text-[10px] uppercase tracking-[0.2em] ${event.level === 'warn' ? 'bg-amber-400/20 text-amber-200' : 'bg-cyan-400/20 text-cyan-200'}`}>{event.level}</span>
          <span class="text-slate-500">{new Date(event.createdAt).toLocaleString()}</span>
        </div>
        <p class="mt-1 text-sm text-slate-300">{event.message}</p>
      </div>
    {/each}
  </div>
</section>
