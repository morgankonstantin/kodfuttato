"""Függvény, ciklus és többsoros kimenet."""


def fibonacci(n):
    a, b = 0, 1
    result = []
    for _ in range(n):
        result.append(a)
        a, b = b, a + b
    return result


if __name__ == "__main__":
    for i, value in enumerate(fibonacci(10)):
        print(f"{i}: {value}")
