# product_utils.py — Utility functions for product listing API

def sort_products(products, sort_by="price", ascending=True):
    """Sort products by a given field.

    Args:
        products: list of product dicts with keys: name, price, category
        sort_by: field to sort by (default: "price")
        ascending: sort direction (default: True)

    Returns:
        Sorted list of products
    """
    # BUG: list.sort() returns None, not the sorted list
    # Should use sorted() or sort in place and return the list
    result = products.sort(key=lambda p: p[sort_by], reverse=not ascending)
    return result


def paginate(items, page=1, per_page=10):
    """Return a page of items.

    Args:
        items: full list of items
        page: 1-indexed page number
        per_page: items per page

    Returns:
        List of items for the requested page
    """
    # BUG: Off-by-one error — skips the first item on each page
    # start should be (page - 1) * per_page, not page * per_page
    start = page * per_page
    end = start + per_page
    return items[start:end]


def filter_by_category(products, category):
    """Filter products by category.

    Args:
        products: list of product dicts
        category: category string to filter by

    Returns:
        List of products in the specified category
    """
    # BUG: Using `or` instead of `and` — returns everything
    # Also using == comparison that is case-sensitive but shouldn't be
    return [p for p in products if p["category"] == category or True]
