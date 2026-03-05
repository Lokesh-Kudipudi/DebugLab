import pytest
import json
import sys
import os

sys.path.insert(0, os.path.join(os.path.dirname(__file__), '..', 'src'))
from api_parser import parse_response, flatten_dict


class TestParseResponse:
    def test_parse_valid_response(self):
        raw = json.dumps({
            "status": "ok",
            "data": [{"id": 1, "name": "Alice"}, {"id": 2, "name": "Bob"}]
        })
        result = parse_response(raw)

        assert result["status"] == "ok"
        assert result["count"] == 2
        assert result["items"][0]["name"] == "Alice"

    def test_second_parse_does_not_contain_first_results(self):
        """Calling parse_response twice should NOT accumulate items."""
        raw1 = json.dumps({"status": "ok", "data": [{"id": 1}]})
        raw2 = json.dumps({"status": "ok", "data": [{"id": 2}]})

        result1 = parse_response(raw1)
        result2 = parse_response(raw2)

        # Second call should only have 1 item, not 2
        assert result2["count"] == 1, (
            f"Expected 1 item but got {result2['count']} — "
            "items from previous call are leaking in"
        )
        assert result2["items"][0]["id"] == 2

    def test_invalid_json_raises_error(self):
        """Malformed JSON should raise an exception, not silently fail."""
        with pytest.raises(Exception):
            parse_response("this is not json {{{")


class TestFlattenDict:
    def test_flat_dict_unchanged(self):
        result = flatten_dict({"a": 1, "b": 2})
        assert result == {"a": 1, "b": 2}

    def test_nested_dict_flattened(self):
        result = flatten_dict({
            "user": {
                "name": "Alice",
                "age": 30,
            },
            "active": True,
        })
        assert result == {
            "user.name": "Alice",
            "user.age": 30,
            "active": True,
        }

    def test_deeply_nested_dict(self):
        result = flatten_dict({
            "a": {
                "b": {
                    "c": 42
                }
            }
        })
        assert result == {"a.b.c": 42}
