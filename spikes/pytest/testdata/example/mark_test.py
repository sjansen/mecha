from time import sleep

import pytest


@pytest.mark.skip
def test__mark__1():
    sleep(0.25)


@pytest.mark.xfail
def test__mark__2():
    sleep(0.25)
    assert False
