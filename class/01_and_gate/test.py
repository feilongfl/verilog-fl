
import wave
import cocotb
from cocotb.triggers import Timer
from cocotb.wavedrom import trace


async def test_map(dut, val):
    dut.a.value = val[0]
    dut.b.value = val[1]
    await Timer(2, units="ns")  # wait a bit

    # dut._log.info("out is %s", dut.out.value)
    assert dut.out.value == val[2], "my_signal_2[0] is not %s!" % val[2]


async def generate_clock(dut):
    """Generate clock pulses."""

    for cycle in range(10):
        dut.clk.value = 0
        await Timer(1, units="ns")
        dut.clk.value = 1
        await Timer(1, units="ns")


@cocotb.test()
async def test(dut):
    """Try accessing the design."""

    # run the clock "in the background"
    await cocotb.start(generate_clock(dut))
    with trace(dut.a, dut.b, dut.out, clk=dut.clk) as waves:
        await Timer(2, units="ns")  # start

        await test_map(dut, (1, 1, 1))
        await test_map(dut, (1, 0, 0))
        await test_map(dut, (0, 1, 0))
        await test_map(dut, (0, 0, 0))

        await Timer(2, units="ns")  # end

        # j = waves.dumpj()
        # print(j)
        waves.write("trace.json")
