Audit the existing ARGUS RCP/get_context implementation for exposing the
persisted belief edge graph.

Do not change code yet.

We have verified:
- get_context returns ALL beliefs in the current scenario
- get_context returns ALL evidence
- get_context returns ALL intents
- belief_edge records are persisted
- but Snapshot currently has no Edges field and GetSnapshot does not query
  belief_edge

Determine the minimum implementation needed to expose the existing edge
graph through argus.get_context.

Answer:

1. Exact current belief_edge schema.
2. Exact existing edge kinds.
3. Exact write path used by submit_packet.
4. Whether edge data can be added to Snapshot without changing the MCP tool
   surface.
5. Exact response shape you recommend.
6. Whether this is a small projection-only change or requires architectural
   changes.
7. Tests needed to prove:
   - derives edges are returned
   - contradicts edges are returned
   - edges reference the correct immutable belief IDs
   - no authority capability is added

Do not implement yet.