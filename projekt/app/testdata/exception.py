"""Kezeletlen kivétel: 1-es kilépési kódot és Tracebacket kell adnia."""


def divide(a, b):
    return a / b


print("a hiba előtti kimenet még megjelenik")
divide(1, 0)
