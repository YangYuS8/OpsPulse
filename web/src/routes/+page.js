export const ssr = false;

export async function load({ fetch }) {
  async function fetchJSON(path) {
    const response = await fetch(path);
    if (!response.ok) {
      throw new Error(`failed to load ${path}`);
    }
    return response.json();
  }

  const [overview, nodes, events] = await Promise.all([
    fetchJSON('/api/v1/overview'),
    fetchJSON('/api/v1/nodes'),
    fetchJSON('/api/v1/events')
  ]);

  return {
    overview,
    nodes,
    events
  };
}
