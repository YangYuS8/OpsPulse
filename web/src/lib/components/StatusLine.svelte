<script>
  export let overview;
  export let transport = 'disconnected';
</script>

<div class="status-line px-3 py-2 text-xs text-terminal-fg">
  OPSPULSE mode=live transport={transport} nodes={overview?.nodesTotal ?? 0} online={overview?.nodesOnline ?? 0} warn={overview?.nodesWarn ?? 0} offline={overview?.nodesOffline ?? 0} containers={overview?.containersUp ?? 0}
</div>
