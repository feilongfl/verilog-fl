module top_module( 
    input a, 
    input b, 
    input clk,
    output out 
    );

    assign out = a & b;

endmodule
