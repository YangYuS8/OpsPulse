<script>
  import { onMount } from 'svelte';
  import EventTimeline from '$components/EventTimeline.svelte';
  import NodeCard from '$components/NodeCard.svelte';
  import NodeDetail from '$components/NodeDetail.svelte';
  import OverviewTiles from '$components/OverviewTiles.svelte';

  export let data;

  let overview = data.overview;
  let nodes = data.nodes;
  let events = data.events;
  let selectedId = data.nodes[0]?.nodeId ?? null;

  $: selectedNode = nodes.find((node) => node.nodeId === selectedId) ?? nodes[0] ?? null;

  onMount(() => {
    const source = new EventSource(`${data.apiBase}/api/v1/stream`);
    source.onmessage = (event) => {
      const payload = JSON.parse(event.data);
      if (payload.type === 'node_update') {
        overview = payload.overview;
        const nextNode = payload.node;
        const index = nodes.findIndex((node) => node.nodeId === nextNode.nodeId);
        if (index >= 0) {
          nodes[index] = nextNode;
          nodes = [...nodes];
        } else {
          nodes = [nextNode, ...nodes];
        }
      }
      if (payload.type === 'event') {
        events = [payload.event, ...events].slice(0, 50);
      }
    };
    return () => source.close();
  });
</script>

<svelte:head>
  <title>OpsPulse</title>
  <meta name="description" content="Lightweight DevOps control center for Linux nodes and services." />
</svelte:head>

<div class="grid-overlay min-h-screen px-4 py-6 sm:px-6 lg:px-10">
  <div class="mx-auto max-w-[1600px] space-y-6">
    <header class="panel overflow-hidden p-6">
      <div class="flex flex-col gap-4 lg:flex-row lg:items-end lg:justify-between">
        <div>
          <div class="text-xs uppercase tracking-[0.4em] text-cyan-300">OpsPulse Command Matrix</div>
          <h1 class="mt-3 text-4xl font-semibold text-white sm:text-5xl">Cyberpunk DevOps Control Center</h1>
          <p class="mt-3 max-w-3xl text-slate-300">
            Live heartbeat intake, service watchlists, container telemetry, and offline detection for Linux nodes without opening inbound ports.
          </p>
        </div>
        <div class="rounded-2xl border border-cyan-400/20 bg-slate-900/70 px-5 py-4 text-sm text-slate-300">
          <div>Total Nodes: <span class="text-white">{overview.nodesTotal}</span></div>
          <div class="mt-1">Service Alerts: <span class="text-amber-200">{overview.servicesDown}</span></div>
        </div>
      </div>
    </header>

    <OverviewTiles {overview} />

    <div class="grid gap-6 xl:grid-cols-[420px_minmax(0,1fr)]">
      <section class="space-y-4">
        {#each nodes as node}
          <div on:click={() => (selectedId = node.nodeId)} on:keydown={() => {}} role="button" tabindex="0">
            <NodeCard node={node} active={node.nodeId === selectedId} />
          </div>
        {/each}
      </section>
      <NodeDetail node={selectedNode} />
    </div>

    <EventTimeline {events} />
  </div>
</div>
