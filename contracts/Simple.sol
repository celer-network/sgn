pragma solidity ^0.5.1;


contract Simple {
    uint public a;
    event Test(uint i);

    function emitEvent(uint i) public {
        a = i;
        emit Test(i);
    }
}