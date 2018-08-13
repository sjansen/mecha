from time import sleep

def test__g__1():
    sleep(0.25)
    raise RuntimeError


def test__g__2():
    sleep(0.25)
    raise SystemError


def test__g__3():
    sleep(0.25)
    raise ZeroDivisionError
