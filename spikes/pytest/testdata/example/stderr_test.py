from sys import stderr
from time import sleep


stderr.write("loading... ")


def test__stderr__1():
    sleep(0.25)
    stderr.write("Spoon!\n")


def test__stderr__2():
    sleep(0.25)
    stderr.write("Kilroy was here.\n")
    assert False


stderr.write("done\n")
