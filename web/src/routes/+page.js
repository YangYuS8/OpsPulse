const apiBase = import.meta.env.PUBLIC_API_BASE || 'http://localhost:8080';

async function fetchJSON(path) {
  const response = await fetch(`${apiBase}${path}`);
  if (!response.ok) {
    throw new Error(`failed to load ${path}`);
  }
  return response.json();
}

export async function load() {
  const [overview, nodes, events] = await Promise.all([
    fetchJSON('/api/v1/overview'),
    fetchJSON('/api/v1/nodes'),
    fetchJSON('/api/v1/events')
  ]);

  return {
    apiBase,
    overview,
    nodes,
    events
  };
}
