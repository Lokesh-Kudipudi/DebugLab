# api_parser.py — Utility for parsing API responses

import json


def parse_response(raw_json, collected_items=[]):
    """Parse a JSON API response and collect items.

    Args:
        raw_json: a JSON string like '{"status": "ok", "data": [...]}'
        collected_items: list to accumulate items into

    Returns:
        dict with keys: status, items, count

    BUG 1: Mutable default argument — collected_items persists between calls,
    so calling parse_response twice accumulates items from both calls.
    """
    try:
        response = json.loads(raw_json)
    except:
        # BUG 2: Bare except catches EVERYTHING including KeyboardInterrupt
        # Should only catch json.JSONDecodeError
        # Also silently returns instead of raising — caller never knows it failed
        return {"status": "error", "items": [], "count": 0}

    status = response.get("status", "unknown")

    # BUG 3: Should raise an error for non-ok status, but doesn't
    data = response.get("data", [])

    for item in data:
        collected_items.append(item)

    return {
        "status": status,
        "items": collected_items,
        "count": len(collected_items),
    }


def flatten_dict(d, parent_key="", sep="."):
    """Flatten a nested dictionary.

    Example:
        {"a": {"b": 1, "c": 2}} -> {"a.b": 1, "a.c": 2}

    Args:
        d: dictionary to flatten
        parent_key: prefix for keys (used in recursion)
        sep: separator between parent and child keys

    Returns:
        Flattened dictionary
    """
    items = {}
    for k, v in d.items():
        new_key = f"{parent_key}{sep}{k}" if parent_key else k
        if isinstance(v, dict):
            # BUG 4: Doesn't merge recursive result — just calls flatten_dict
            # but doesn't use the return value
            flatten_dict(v, new_key, sep)
        else:
            items[new_key] = v
    return items
