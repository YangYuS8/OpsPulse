<script>
  import { createEventDispatcher } from 'svelte';

  export let value = '';
  export let statusLine = '';

  const dispatch = createEventDispatcher();

  function submit() {
    dispatch('submit', value);
  }
</script>

<section class="terminal-pane p-3">
  <div class="mb-2 text-xs uppercase tracking-[0.24em] text-terminal-muted">Command Bar</div>
  <div class="flex items-center gap-2 rounded border border-terminal-border bg-terminal-bg px-3 py-2 text-sm">
    <span class="text-terminal-green">opspulse&gt;</span>
    <input
      bind:value
      class="min-w-0 flex-1 border-0 bg-transparent p-0 font-mono text-terminal-fg focus:ring-0"
      placeholder="help | nodes | inspect <node> | events | clear"
      on:keydown={(event) => event.key === 'Enter' && submit()}
    />
  </div>
  <div class="mt-2 min-h-5 text-xs text-terminal-dim">{statusLine}</div>
</section>
