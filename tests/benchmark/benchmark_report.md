# Goscrape Benchmark Report

This document records the performance of the Tag engine functions.

## Methodology
- **Environment**: Linux (Ubuntu), Intel(R) Core(TM) i5-12500
- **Test Data**: Static HTML file with ~1KB of semi-complex nested tags.
- **Iterations**: Handled by Go's `testing` package (automatic scale until statistically stable).

## Results (Last Recorded)

| Function | Performance (ns/op) | Memory (B/op) | Allocs (allocs/op) | Note |
| :--- | :--- | :--- | :--- | :--- |
| `Find` | 1,259 | 144 | 1 | Stops at first match. Most efficient raw search. |
| `Select` | 525 | 176 | 3 | Optimized lazy-search with CSS selectors. |
| `FindAll` | 1,979 | 1,712 | 16 | Returns a full collection using attribute search. |
| `SelectAll` | 3,377 | 1,744 | 18 | Fully parses CSS and returns a collection. |

## Strategy Summary
- **Lazy Loading**: `Select` and `FindFirst` use `limit=1` to stop tree traversal early.
- **Memory Reuse**: We minimize intermediate slice allocations by using callbacks where possible.
- **CSS Parser**: Powered by `cascadia`, but optimized with our internal `traverse` engine.
