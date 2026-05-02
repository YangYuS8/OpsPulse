<script>
  import { createEventDispatcher } from 'svelte';

  export let value = '';
  export let statusLine = '';

  const dispatch = createEventDispatcher();

  function submit() {
    dispatch('submit', value);
  }
</script>

<section class="pane px-3 py-2">
  <div class="pane-title mb-2 border-b-0 pb-0">SHELL</div>
  <div class="flex items-center gap-2 text-sm">
    <span class="text-terminal-green">opspulse&gt;</span>
    <input
      bind:value
      class="min-w-0 flex-1 border-0 bg-transparent p-0 font-mono text-terminal-fg focus:ring-0"
      placeholder="help | nodes | inspect <node> | events | clear"
      on:keydown={(event) => event.key === 'Enter' && submit()}
    />
    <span class="text-xs text-terminal-dim">{statusLine}</span>
  </div>
</section>
