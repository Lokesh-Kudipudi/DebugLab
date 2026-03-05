import pytest
import sys
import os

sys.path.insert(0, os.path.join(os.path.dirname(__file__), '..', 'src'))
from product_utils import sort_products, paginate, filter_by_category


SAMPLE_PRODUCTS = [
    {"name": "Laptop", "price": 999, "category": "Electronics"},
    {"name": "Shirt", "price": 29, "category": "Clothing"},
    {"name": "Phone", "price": 699, "category": "Electronics"},
    {"name": "Pants", "price": 49, "category": "Clothing"},
    {"name": "Headphones", "price": 149, "category": "Electronics"},
]


class TestSortProducts:
    def test_sort_by_price_ascending(self):
        products = list(SAMPLE_PRODUCTS)  # copy to avoid mutation issues
        result = sort_products(products, sort_by="price", ascending=True)

        assert result is not None, "sort_products should return a list, got None"
        assert len(result) == 5
        assert result[0]["price"] == 29
        assert result[-1]["price"] == 999

    def test_sort_by_price_descending(self):
        products = list(SAMPLE_PRODUCTS)
        result = sort_products(products, sort_by="price", ascending=False)

        assert result is not None, "sort_products should return a list, got None"
        assert result[0]["price"] == 999
        assert result[-1]["price"] == 29

    def test_sort_by_name(self):
        products = list(SAMPLE_PRODUCTS)
        result = sort_products(products, sort_by="name", ascending=True)

        assert result is not None
        assert result[0]["name"] == "Headphones"
        assert result[-1]["name"] == "Shirt"


class TestPaginate:
    def test_first_page(self):
        items = list(range(1, 26))  # 1 to 25
        page = paginate(items, page=1, per_page=10)

        assert len(page) == 10
        assert page[0] == 1
        assert page[-1] == 10

    def test_second_page(self):
        items = list(range(1, 26))
        page = paginate(items, page=2, per_page=10)

        assert len(page) == 10
        assert page[0] == 11
        assert page[-1] == 20

    def test_last_partial_page(self):
        items = list(range(1, 26))
        page = paginate(items, page=3, per_page=10)

        assert len(page) == 5
        assert page[0] == 21
        assert page[-1] == 25


class TestFilterByCategory:
    def test_filter_electronics(self):
        result = filter_by_category(SAMPLE_PRODUCTS, "Electronics")

        assert len(result) == 3
        assert all(p["category"] == "Electronics" for p in result)

    def test_filter_clothing(self):
        result = filter_by_category(SAMPLE_PRODUCTS, "Clothing")

        assert len(result) == 2
        assert all(p["category"] == "Clothing" for p in result)

    def test_filter_nonexistent_category(self):
        result = filter_by_category(SAMPLE_PRODUCTS, "Books")

        assert len(result) == 0
