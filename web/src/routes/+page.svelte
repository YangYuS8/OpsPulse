<script>
  import { onMount } from 'svelte';
  import CommandBar from '$components/CommandBar.svelte';
  import EventLogPane from '$components/EventLogPane.svelte';
  import NodeListPane from '$components/NodeListPane.svelte';
  import NodeDetail from '$components/NodeDetail.svelte';
  import ServicePane from '$components/ServicePane.svelte';
  import StatusLine from '$components/StatusLine.svelte';

  export let data;

  let overview = data.overview;
  let nodes = data.nodes ?? [];
  let events = data.events ?? [];
  let selectedId = nodes[0]?.nodeId ?? null;
  let command = '';
  let statusLine = 'type help for available commands';
  let transport = 'disconnected';
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

  function mergeEventLog() {
    return [...commandFeed, ...events].slice(0, 80);
  }

  onMount(() => {
    const source = new EventSource('/api/v1/stream');
    source.onopen = () => {
      transport = 'sse';
    };
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
    source.onerror = () => {
      transport = 'disconnected';
    };
    return () => {
      transport = 'disconnected';
      source.close();
    };
  });
</script>

<svelte:head>
  <title>OpsPulse</title>
  <meta name="description" content="Terminal-first local DevOps console for Linux nodes and services." />
</svelte:head>

<div class="min-h-screen px-2 py-2 sm:px-3 lg:px-4">
  <div class="terminal-frame mx-auto max-w-[1800px]">
    <StatusLine {overview} {transport} />

    <div class="frame-main">
      <NodeListPane
        {nodes}
        {selectedId}
        onSelect={selectNode}
        onKeySelect={handleNodeKeydown}
      />
      <NodeDetail node={selectedNode} />
      <ServicePane node={selectedNode} />
    </div>

    <EventLogPane events={mergeEventLog()} title="EventLogPane" />

    <CommandBar bind:value={command} {statusLine} on:submit={handleCommandSubmit} />
  </div>
</div>
