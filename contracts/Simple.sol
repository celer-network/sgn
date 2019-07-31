pragma solidity ^0.5.1;


contract Simple {
    event Test(uint i);

    function emitEvent() public {
        emit Test(1);
    }
}