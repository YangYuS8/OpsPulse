<script>
  import { onMount } from 'svelte';
  import CommandBar from '$components/CommandBar.svelte';
  import EventTimeline from '$components/EventTimeline.svelte';
  import NodeCard from '$components/NodeCard.svelte';
  import NodeDetail from '$components/NodeDetail.svelte';
  import OverviewTiles from '$components/OverviewTiles.svelte';
  import TerminalHeader from '$components/TerminalHeader.svelte';

  export let data;

  let overview = data.overview;
  let nodes = data.nodes ?? [];
  let events = data.events ?? [];
  let selectedId = nodes[0]?.nodeId ?? null;
  let command = '';
  let statusLine = 'type help for available commands';
  let commandFeed = [
    { id: 'boot', nodeId: 'local-shell', level: 'info', message: 'OpsPulse terminal ready', createdAt: new Date().toISOString() }
  ];

  $: selectedNode = nodes.find((node) => node.nodeId === selectedId) ?? nodes[0] ?? null;

  function applyNodeUpdate(nextNode) {
    const index = nodes.findIndex((node) => node.nodeId === nextNode.nodeId);
    if (index >= 0) {
      nodes[index] = nextNode;
      nodes = [...nodes];
      return;
    }
    nodes = [nextNode, ...nodes];
  }

  function pushCommandMessage(message, level = 'info', nodeId = 'opspulse-shell') {
    commandFeed = [
      {
        id: `${Date.now()}-${Math.random()}`,
        nodeId,
        level,
        message,
        createdAt: new Date().toISOString()
      },
      ...commandFeed
    ].slice(0, 20);
  }

  function runCommand(raw) {
    const input = raw.trim();
    if (!input) return;
    const [commandName, ...args] = input.split(/\s+/);
    switch (commandName) {
      case 'help':
        statusLine = 'help | nodes | inspect <nodeId|hostname> | events | clear';
        pushCommandMessage('available commands: help, nodes, inspect <node>, events, clear');
        break;
      case 'nodes':
        statusLine = `loaded ${nodes.length} nodes in left pane`;
        pushCommandMessage(`nodes loaded: ${nodes.map((node) => node.hostname).join(', ')}`);
        break;
      case 'inspect': {
        const target = args.join(' ');
        const node = nodes.find((item) => item.nodeId === target || item.hostname === target);
        if (node) {
          selectedId = node.nodeId;
          statusLine = `inspect switched to ${node.hostname}`;
          pushCommandMessage(`inspect -> ${node.hostname}`);
        } else {
          statusLine = `node not found: ${target || '<empty>'}`;
          pushCommandMessage(`inspect failed for ${target || '<empty>'}`, 'warn');
        }
        break;
      }
      case 'events':
        statusLine = `showing ${events.length} live events and ${commandFeed.length} shell entries`;
        pushCommandMessage('event log focused');
        break;
      case 'clear':
        commandFeed = [];
        statusLine = 'shell output cleared';
        break;
      default:
        statusLine = `unknown command: ${commandName}`;
        pushCommandMessage(`unknown command: ${commandName}`, 'warn');
    }
  }

  function selectNode(nodeId) {
    selectedId = nodeId;
  }

  function handleNodeKeydown(event, nodeId) {
    if (event.key === 'Enter' || event.key === ' ') {
      event.preventDefault();
      selectNode(nodeId);
    }
  }

  function handleCommandSubmit() {
    runCommand(command);
    command = '';
  }

  onMount(() => {
    const source = new EventSource('/api/v1/stream');
    source.onmessage = (event) => {
      const payload = JSON.parse(event.data);
      if (payload.type === 'node_update') {
        overview = payload.overview;
        applyNodeUpdate(payload.node);
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
  <meta name="description" content="Terminal-first local DevOps console for Linux nodes and services." />
</svelte:head>

<div class="min-h-screen px-3 py-4 sm:px-4 lg:px-6">
  <div class="mx-auto max-w-[1680px] space-y-4">
    <TerminalHeader {overview} />
    <OverviewTiles {overview} />

    <div class="grid gap-4 xl:grid-cols-[360px_minmax(0,1fr)]">
      <section class="terminal-pane p-4">
        <div class="mb-3 flex items-center justify-between border-b border-terminal-border pb-3">
          <h2 class="text-sm uppercase tracking-[0.24em] text-terminal-fg">node pane</h2>
          <span class="text-[11px] text-terminal-muted">{nodes.length} entries</span>
        </div>
        {#if nodes.length > 0}
          <div class="space-y-3">
          {#each nodes as node}
            <div on:click={() => selectNode(node.nodeId)} on:keydown={(event) => handleNodeKeydown(event, node.nodeId)} role="button" tabindex="0">
              <NodeCard node={node} active={node.nodeId === selectedId} />
            </div>
          {/each}
          </div>
        {:else}
          <div class="py-8 text-sm text-terminal-dim">
            no nodes available; start an agent or enable demo mode
          </div>
        {/if}
      </section>
      <NodeDetail node={selectedNode} />
    </div>

    <div class="grid gap-4 xl:grid-cols-2">
      <EventTimeline events={[...commandFeed, ...events].slice(0, 50)} title="shell + live log" />
      <EventTimeline {events} title="cluster event stream" />
    </div>

    <CommandBar bind:value={command} {statusLine} on:submit={handleCommandSubmit} />
  </div>
</div>
